package txserver

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	evidencetypes "cosmossdk.io/x/evidence/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/RWAs-labs/ethermint/crypto/hd"
	etherminttypes "github.com/RWAs-labs/ethermint/types"
	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/samber/lo"

	"github.com/RWAs-labs/muse/app"
	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	emissionstypes "github.com/RWAs-labs/muse/x/emissions/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	lightclienttypes "github.com/RWAs-labs/muse/x/lightclient/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// SystemContractAddresses contains the addresses of the system contracts deployed
type SystemContractAddresses struct {
	UniswapV2FactoryAddr, UniswapV2RouterAddr, MEVMConnectorAddr, WMUSEAddr, ERC20mrc20Addr string
}

// MRC20Deployment configures deployment of MRC20 contracts
type MRC20Deployment struct {
	ERC20Addr common.Address
	SPLAddr   *solana.PublicKey // if nil - no SPL MRC20 is deployed
}

// MRC20Addresses contains the addresses of deployed MRC20 contracts
type MRC20Addresses struct {
	ERC20MRC20Addr common.Address
	SPLMRC20Addr   common.Address
}

// EmissionsPoolAddress is the address of the emissions pool
// This address is constant for all networks because it is derived from emissions name
const EmissionsPoolAddress = "muse1w43fn2ze2wyhu5hfmegr6vp52c3dgn0srdgymy"

// MuseTxServer is a MuseChain tx server for E2E test
type MuseTxServer struct {
	ctx               context.Context
	clientCtx         client.Context
	txFactory         tx.Factory
	name              []string
	mnemonic          []string
	address           []string
	blockTimeout      time.Duration
	authorityClient   authoritytypes.QueryClient
	validatorsKeyring keyring.Keyring
}

// NewMuseTxServer returns a new TxServer with provided account
func NewMuseTxServer(rpcAddr string, names []string, privateKeys []string, chainID string) (*MuseTxServer, error) {
	ctx := context.Background()

	if len(names) == 0 {
		return nil, errors.New("no account provided")
	}

	if len(names) != len(privateKeys) {
		return nil, errors.New("invalid names and privateKeys")
	}

	// initialize rpc and check status
	rpc, err := rpchttp.New(rpcAddr, "/websocket")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rpc: %s", err.Error())
	}
	if _, err = rpc.Status(ctx); err != nil {
		return nil, fmt.Errorf("failed to query rpc: %s", err.Error())
	}

	// initialize codec
	cdc, reg := newCodec()

	// initialize keyring
	kr := keyring.NewInMemory(cdc, hd.EthSecp256k1Option())

	addresses := make([]string, 0, len(names))

	// create accounts
	for i := range names {
		err = kr.ImportPrivKeyHex(names[i], privateKeys[i], string(hd.EthSecp256k1Type))
		if err != nil {
			return nil, fmt.Errorf("failed to create account: %w", err)
		}
		r, err := kr.Key(names[i])
		if err != nil {
			return nil, fmt.Errorf("failed to get account key: %w", err)
		}
		accAddr, err := r.GetAddress()
		if err != nil {
			return nil, fmt.Errorf("failed to get account address: %w", err)
		}

		addresses = append(addresses, accAddr.String())
	}

	// initialize validators keyring. The validator keys have been copied to this location in the docker image.
	// Refer contrib/localnet/orchestrator/Dockerfile.fastbuild for more details
	validatorsKeyring, err := keyring.New(
		"validatorKeys",
		keyring.BackendTest,
		"/root/.musecored/",
		nil,
		cdc,
		hd.EthSecp256k1Option(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator keyring: %w", err)
	}

	clientCtx := newContext(rpc, cdc, reg, kr, chainID)
	txf := newFactory(clientCtx)

	return &MuseTxServer{
		ctx:               ctx,
		clientCtx:         clientCtx,
		txFactory:         txf,
		name:              names,
		address:           addresses,
		blockTimeout:      2 * time.Minute,
		validatorsKeyring: validatorsKeyring,
	}, nil
}

// UpdateKeyring updates the keyring being used by the MuseTxServer
// This returns a copy of the MuseTxServer with the updated keyring
func (zts *MuseTxServer) UpdateKeyring(kr keyring.Keyring) MuseTxServer {
	ztsUpdated := *zts
	ztsUpdated.clientCtx = ztsUpdated.clientCtx.WithKeyring(kr)
	ztsUpdated.txFactory = newFactory(ztsUpdated.clientCtx)
	return ztsUpdated
}

// GetAccountName returns the account name from the given index
// returns empty string if index is out of bound, caller should handle error
func (zts *MuseTxServer) GetAccountName(index int) string {
	if index >= len(zts.name) {
		return ""
	}
	return zts.name[index]
}

// GetCodec returns the codec for the MuseTxServer
func (zts *MuseTxServer) GetCodec() codec.Codec {
	return zts.clientCtx.Codec
}

// GetAccountAddress returns the account address from the given index
// returns empty string if index is out of bound, error should be handled by caller
func (zts *MuseTxServer) GetAccountAddress(index int) string {
	if index >= len(zts.address) {
		return ""
	}
	return zts.address[index]
}

// GetAccountAddressFromName returns the account address from the given name
func (zts *MuseTxServer) GetAccountAddressFromName(name string) (string, error) {
	acc, err := zts.clientCtx.Keyring.Key(name)
	if err != nil {
		return "", err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

// MustGetAccountAddressFromName returns the account address from the given name.It panics on error
// and should be used in tests only
func (zts *MuseTxServer) MustGetAccountAddressFromName(name string) string {
	acc, err := zts.clientCtx.Keyring.Key(name)
	if err != nil {
		panic(err)
	}
	addr, err := acc.GetAddress()
	if err != nil {
		panic(err)
	}
	return addr.String()
}

// GetValidatorsKeyring returns the validators keyring
func (zts *MuseTxServer) GetValidatorsKeyring() keyring.Keyring {
	return zts.validatorsKeyring
}

// GetAllAccountAddress returns all account addresses
func (zts *MuseTxServer) GetAllAccountAddress() []string {
	return zts.address
}

// GetAccountMnemonic returns the account name from the given index
// returns empty string if index is out of bound, error should be handled by caller
func (zts *MuseTxServer) GetAccountMnemonic(index int) string {
	if index >= len(zts.mnemonic) {
		return ""
	}
	return zts.mnemonic[index]
}

// BroadcastTx broadcasts a tx to MuseChain with the provided msg from the account
// and waiting for blockTime for tx to be included in the block
func (zts *MuseTxServer) BroadcastTx(account string, msgs ...sdktypes.Msg) (*sdktypes.TxResponse, error) {
	// Find number and sequence and set it
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return nil, err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return nil, err
	}
	accountNumber, accountSeq, err := zts.clientCtx.AccountRetriever.GetAccountNumberSequence(zts.clientCtx, addr)
	if err != nil {
		return nil, err
	}
	zts.txFactory = zts.txFactory.WithAccountNumber(accountNumber).WithSequence(accountSeq)

	txBuilder, err := zts.txFactory.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, err
	}
	// increase gas and fees if multiple messages are provided
	txBuilder.SetGasLimit(zts.txFactory.Gas() * uint64(len(msgs)))
	txBuilder.SetFeeAmount(zts.txFactory.Fees().MulInt(sdkmath.NewInt(int64(len(msgs)))))

	// Sign tx
	err = tx.Sign(context.TODO(), zts.txFactory, account, txBuilder, true)
	if err != nil {
		return nil, err
	}
	txBytes, err := zts.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}
	return broadcastWithBlockTimeout(zts, txBytes)
}

func broadcastWithBlockTimeout(zts *MuseTxServer, txBytes []byte) (*sdktypes.TxResponse, error) {
	res, err := zts.clientCtx.BroadcastTx(txBytes)
	if err != nil {
		if res == nil {
			return nil, err
		}
		return &sdktypes.TxResponse{
			Code:      res.Code,
			Codespace: res.Codespace,
			TxHash:    res.TxHash,
		}, err
	}
	if res.Code != 0 {
		return res, fmt.Errorf("broadcast failed: %s", res.RawLog)
	}

	exitAfter := time.After(zts.blockTimeout)
	hash, err := hex.DecodeString(res.TxHash)
	if err != nil {
		return nil, err
	}
	for {
		select {
		case <-exitAfter:
			return nil, fmt.Errorf("timed out after waiting for tx to get included in the block: %d", zts.blockTimeout)
		case <-time.After(time.Millisecond * 100):
			resTx, err := zts.clientCtx.Client.Tx(zts.ctx, hash, false)
			if err == nil {
				return mkTxResult(zts.ctx, zts.clientCtx, resTx)
			}
		}
	}
}

func mkTxResult(
	ctx context.Context,
	clientCtx client.Context,
	resTx *coretypes.ResultTx,
) (*sdktypes.TxResponse, error) {
	if resTx.TxResult.Code != 0 {
		return nil, fmt.Errorf("tx failed: %s", resTx.TxResult.Log)
	}
	txb, err := clientCtx.TxConfig.TxDecoder()(resTx.Tx)
	if err != nil {
		return nil, err
	}
	p, ok := txb.(intoAny)
	if !ok {
		return nil, fmt.Errorf("expecting a type implementing intoAny, got: %T", txb)
	}
	resBlock, err := clientCtx.Client.Block(ctx, &resTx.Height)
	if err != nil {
		return nil, err
	}
	return sdktypes.NewResponseResultTx(resTx, p.AsAny(), resBlock.Block.Time.Format(time.RFC3339)), nil
}

type intoAny interface {
	AsAny() *codectypes.Any
}

// UpdateChainParams updates the chain params for the given chain ID
// It takes into account that the required policy get updated in v28 operational -> admin
func (zts *MuseTxServer) UpdateChainParams(chainParams *observertypes.ChainParams) error {
	_, err := zts.BroadcastTx(utils.AdminPolicyName, observertypes.NewMsgUpdateChainParams(
		zts.MustGetAccountAddressFromName(utils.AdminPolicyName),
		chainParams,
	))
	return err
}

// EnableHeaderVerification enables the header verification for the given chain IDs
func (zts *MuseTxServer) EnableHeaderVerification(account string, chainIDList []int64) error {
	// retrieve account
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return err
	}

	_, err = zts.BroadcastTx(account, lightclienttypes.NewMsgEnableHeaderVerification(
		addr.String(),
		chainIDList,
	))
	return err
}

// UpdateGatewayAddress updates the gateway address
func (zts *MuseTxServer) UpdateGatewayAddress(account, gatewayAddr string) error {
	// retrieve account
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return err
	}

	_, err = zts.BroadcastTx(account, fungibletypes.NewMsgUpdateGatewayContract(
		addr.String(),
		gatewayAddr,
	))

	return err
}

// DeploySystemContracts deploys the system contracts
// returns the addresses of uniswap factory, router
func (zts *MuseTxServer) DeploySystemContracts(
	accountOperational, accountAdmin string,
) (SystemContractAddresses, error) {
	// retrieve account
	accOperational, err := zts.clientCtx.Keyring.Key(accountOperational)
	if err != nil {
		return SystemContractAddresses{}, err
	}
	addrOperational, err := accOperational.GetAddress()
	if err != nil {
		return SystemContractAddresses{}, err
	}
	accAdmin, err := zts.clientCtx.Keyring.Key(accountAdmin)
	if err != nil {
		return SystemContractAddresses{}, err
	}
	addrAdmin, err := accAdmin.GetAddress()
	if err != nil {
		return SystemContractAddresses{}, err
	}

	// deploy new system contracts
	res, err := zts.BroadcastTx(accountOperational, fungibletypes.NewMsgDeploySystemContracts(addrOperational.String()))
	if err != nil {
		return SystemContractAddresses{}, fmt.Errorf("failed to deploy system contracts: %s", err.Error())
	}

	deployedEvent, ok := EventOfType[*fungibletypes.EventSystemContractsDeployed](res.Events)
	if !ok {
		return SystemContractAddresses{}, fmt.Errorf("no EventSystemContractsDeployed in %s", res.TxHash)
	}

	// get system contract
	_, err = zts.BroadcastTx(
		accountAdmin,
		fungibletypes.NewMsgUpdateSystemContract(addrAdmin.String(), deployedEvent.SystemContract),
	)
	if err != nil {
		return SystemContractAddresses{}, fmt.Errorf("failed to set system contract: %s", err.Error())
	}

	return SystemContractAddresses{
		UniswapV2FactoryAddr: deployedEvent.UniswapV2Factory,
		UniswapV2RouterAddr:  deployedEvent.UniswapV2Router,
		MEVMConnectorAddr:    deployedEvent.ConnectorMevm,
		WMUSEAddr:            deployedEvent.Wmuse,
	}, nil
}

// DeployMRC20s deploys the MRC20 contracts
// returns the addresses of erc20 and spl mrc20
func (zts *MuseTxServer) DeployMRC20s(
	mrc20Deployment MRC20Deployment,
	skipChain func(chainID int64) bool,
) (*MRC20Addresses, error) {
	// retrieve account
	accOperational, err := zts.clientCtx.Keyring.Key(utils.OperationalPolicyName)
	if err != nil {
		return nil, err
	}
	addrOperational, err := accOperational.GetAddress()
	if err != nil {
		return nil, err
	}
	accAdmin, err := zts.clientCtx.Keyring.Key(utils.AdminPolicyName)
	if err != nil {
		return nil, err
	}
	addrAdmin, err := accAdmin.GetAddress()
	if err != nil {
		return nil, err
	}

	// authorization for deploying new MRC20 has changed from accountOperational to accountAdmin in v19
	// we use this query to check the current authorization for the message
	// if pre v19 the query is not implement and authorization is operational
	deployerAccount := utils.AdminPolicyName
	deployerAddr := addrAdmin.String()
	authorization, preV19, err := zts.fetchMessagePermissions(&fungibletypes.MsgDeployFungibleCoinMRC20{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch message permissions: %s", err.Error())
	}
	if preV19 || authorization == authoritytypes.PolicyType_groupOperational {
		deployerAccount = utils.OperationalPolicyName
		deployerAddr = addrOperational.String()
	}

	deployMsgs := []*fungibletypes.MsgDeployFungibleCoinMRC20{
		fungibletypes.NewMsgDeployFungibleCoinMRC20(
			deployerAddr,
			"",
			chains.GoerliLocalnet.ChainId,
			18,
			"ETH",
			"gETH",
			coin.CoinType_Gas,
			100000,
			nil,
		),
		fungibletypes.NewMsgDeployFungibleCoinMRC20(
			deployerAddr,
			"",
			chains.BitcoinRegtest.ChainId,
			8,
			"BTC",
			"tBTC",
			coin.CoinType_Gas,
			100000,
			nil,
		),
		fungibletypes.NewMsgDeployFungibleCoinMRC20(
			deployerAddr,
			"",
			chains.SolanaLocalnet.ChainId,
			9,
			"Solana",
			"SOL",
			coin.CoinType_Gas,
			100000,
			nil,
		),
		fungibletypes.NewMsgDeployFungibleCoinMRC20(
			deployerAddr,
			"",
			chains.TONLocalnet.ChainId,
			9,
			"TON",
			"TON",
			coin.CoinType_Gas,
			100_000,
			nil,
		),
		fungibletypes.NewMsgDeployFungibleCoinMRC20(
			deployerAddr,
			mrc20Deployment.ERC20Addr.Hex(),
			chains.GoerliLocalnet.ChainId,
			6,
			"USDT",
			"USDT",
			coin.CoinType_ERC20,
			100000,
			nil,
		),
	}

	if mrc20Deployment.SPLAddr != nil {
		deployMsgs = append(deployMsgs, fungibletypes.NewMsgDeployFungibleCoinMRC20(
			deployerAddr,
			mrc20Deployment.SPLAddr.String(),
			chains.SolanaLocalnet.ChainId,
			9,
			"USDT",
			"USDT",
			coin.CoinType_ERC20,
			100000,
			nil,
		))
	}

	// apply skipChain filter and convert to sdk.Msg
	deployMsgsIface := lo.FilterMap(
		deployMsgs,
		func(msg *fungibletypes.MsgDeployFungibleCoinMRC20, _ int) (sdktypes.Msg, bool) {
			if skipChain(msg.ForeignChainId) {
				return nil, false
			}
			return msg, true
		},
	)

	res, err := zts.BroadcastTx(deployerAccount, deployMsgsIface...)
	if err != nil {
		return nil, fmt.Errorf("deploy mrc20s: %w", err)
	}

	deployedEvents, ok := EventsOfType[*fungibletypes.EventMRC20Deployed](res.Events)
	if !ok {
		return nil, fmt.Errorf("no EventMRC20Deployed in %s", res.TxHash)
	}

	mrc20Addrs := lo.Map(deployedEvents, func(ev *fungibletypes.EventMRC20Deployed, _ int) string {
		return ev.Contract
	})

	err = zts.InitializeLiquidityCaps(mrc20Addrs...)
	if err != nil {
		return nil, fmt.Errorf("initialize liquidity cap: %w", err)
	}

	// find erc20 mrc20
	erc20mrc20, ok := lo.Find(deployedEvents, func(ev *fungibletypes.EventMRC20Deployed) bool {
		return ev.ChainId == chains.GoerliLocalnet.ChainId && ev.CoinType == coin.CoinType_ERC20
	})
	if !ok {
		return nil, fmt.Errorf("unable to find erc20 mrc20")
	}

	// find spl mrc20
	splmrc20Addr := common.Address{}
	if mrc20Deployment.SPLAddr != nil {
		splmrc20, ok := lo.Find(deployedEvents, func(ev *fungibletypes.EventMRC20Deployed) bool {
			return ev.ChainId == chains.SolanaLocalnet.ChainId && ev.CoinType == coin.CoinType_ERC20
		})
		if !ok {
			return nil, fmt.Errorf("unable to find spl mrc20")
		}

		splmrc20Addr = common.HexToAddress(splmrc20.Contract)
	}

	return &MRC20Addresses{
		ERC20MRC20Addr: common.HexToAddress(erc20mrc20.Contract),
		SPLMRC20Addr:   splmrc20Addr,
	}, nil
}

// FundEmissionsPool funds the emissions pool with the given amount
func (zts *MuseTxServer) FundEmissionsPool(account string, amount *big.Int) error {
	// retrieve account
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return err
	}

	// retrieve account address
	emissionPoolAccAddr, err := sdktypes.AccAddressFromBech32(EmissionsPoolAddress)
	if err != nil {
		return err
	}

	// convert amount
	amountInt := sdkmath.NewIntFromBigInt(amount)

	// fund emissions pool
	_, err = zts.BroadcastTx(account, banktypes.NewMsgSend(
		addr,
		emissionPoolAccAddr,
		sdktypes.NewCoins(sdktypes.NewCoin(config.BaseDenom, amountInt)),
	))
	return err
}

func (zts *MuseTxServer) WithdrawAllEmissions(withdrawAmount sdkmath.Int, account, observer string) error {
	// retrieve account
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return fmt.Errorf("failed to get withdrawer account: %w", err)
	}
	withdrawerAddress, err := acc.GetAddress()
	if err != nil {
		return fmt.Errorf("failed to get withdrawer account address: %w", err)
	}

	msg := emissionstypes.MsgWithdrawEmission{
		Creator: observer,
		Amount:  withdrawAmount,
	}

	authzMessage := authz.NewMsgExec(withdrawerAddress, []sdktypes.Msg{&msg})

	_, err = zts.BroadcastTx(account, &authzMessage)
	return err
}

// UpdateKeygen sets a new keygen height . The new height is the current height + 30
func (zts *MuseTxServer) UpdateKeygen(height int64) error {
	keygenHeight := height + 30
	_, err := zts.BroadcastTx(zts.GetAccountName(0), observertypes.NewMsgUpdateKeygen(
		zts.GetAccountAddress(0),
		keygenHeight,
	))
	return err
}

// SetAuthorityClient sets the authority client
func (zts *MuseTxServer) SetAuthorityClient(authorityClient authoritytypes.QueryClient) {
	zts.authorityClient = authorityClient
}

// InitializeLiquidityCaps initializes the liquidity cap for the given coin with a large value
func (zts *MuseTxServer) InitializeLiquidityCaps(mrc20s ...string) error {
	liquidityCap := sdkmath.NewUint(1e18).MulUint64(1e12)

	msgs := make([]sdktypes.Msg, len(mrc20s))
	for i, mrc20 := range mrc20s {
		msgs[i] = fungibletypes.NewMsgUpdateMRC20LiquidityCap(
			zts.MustGetAccountAddressFromName(utils.OperationalPolicyName),
			mrc20,
			liquidityCap,
		)
	}
	_, err := zts.BroadcastTx(utils.OperationalPolicyName, msgs...)
	return err
}

// fetchMessagePermissions fetches the message permissions for a given message
// return a bool preV19 to indicate the node is preV19 and the query doesn't exist
func (zts *MuseTxServer) fetchMessagePermissions(msg sdktypes.Msg) (authoritytypes.PolicyType, bool, error) {
	msgURL := sdktypes.MsgTypeURL(msg)

	res, err := zts.authorityClient.Authorization(zts.ctx, &authoritytypes.QueryAuthorizationRequest{
		MsgUrl: msgURL,
	})

	// check if error is unknown method
	if err != nil {
		if strings.Contains(err.Error(), "unknown method") {
			return authoritytypes.PolicyType_groupOperational, true, nil
		}
		return authoritytypes.PolicyType_groupOperational, false, err
	}

	return res.Authorization.AuthorizedPolicy, false, nil
}

// newCodec returns the codec for msg server
func newCodec() (*codec.ProtoCodec, codectypes.InterfaceRegistry) {
	encodingConfig := app.MakeEncodingConfig()
	interfaceRegistry := encodingConfig.InterfaceRegistry
	cdc := codec.NewProtoCodec(interfaceRegistry)

	sdktypes.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	authz.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)
	stakingtypes.RegisterInterfaces(interfaceRegistry)
	slashingtypes.RegisterInterfaces(interfaceRegistry)
	upgradetypes.RegisterInterfaces(interfaceRegistry)
	distrtypes.RegisterInterfaces(interfaceRegistry)
	evidencetypes.RegisterInterfaces(interfaceRegistry)
	crisistypes.RegisterInterfaces(interfaceRegistry)
	evmtypes.RegisterInterfaces(interfaceRegistry)
	etherminttypes.RegisterInterfaces(interfaceRegistry)
	crosschaintypes.RegisterInterfaces(interfaceRegistry)
	emissionstypes.RegisterInterfaces(interfaceRegistry)
	fungibletypes.RegisterInterfaces(interfaceRegistry)
	observertypes.RegisterInterfaces(interfaceRegistry)
	lightclienttypes.RegisterInterfaces(interfaceRegistry)
	authoritytypes.RegisterInterfaces(interfaceRegistry)

	return cdc, interfaceRegistry
}

// newContext returns the client context for msg server
func newContext(
	rpc *rpchttp.HTTP,
	cdc *codec.ProtoCodec,
	reg codectypes.InterfaceRegistry,
	kr keyring.Keyring,
	chainID string,
) client.Context {
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	return client.Context{}.
		WithChainID(chainID).
		WithInterfaceRegistry(reg).
		WithCodec(cdc).
		WithTxConfig(txConfig).
		WithLegacyAmino(codec.NewLegacyAmino()).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithBroadcastMode(flags.BroadcastSync).
		WithClient(rpc).
		WithSkipConfirmation(true).
		WithFromName("creator").
		WithFromAddress(sdktypes.AccAddress{}).
		WithKeyring(kr).
		WithAccountRetriever(authtypes.AccountRetriever{})
}

// newFactory returns the tx factory for msg server
func newFactory(clientCtx client.Context) tx.Factory {
	return tx.Factory{}.
		WithChainID(clientCtx.ChainID).
		WithKeybase(clientCtx.Keyring).
		WithGas(10000000).
		WithGasAdjustment(1).
		WithSignMode(signing.SignMode_SIGN_MODE_UNSPECIFIED).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithTxConfig(clientCtx.TxConfig).
		WithFees("100000000000000000amuse")
}

// EventsOfType gets events of a specified type
func EventsOfType[T proto.Message](events []abci.Event) ([]T, bool) {
	var filteredEvents []T
	for _, ev := range events {
		pEvent, err := sdktypes.ParseTypedEvent(ev)
		if err != nil {
			continue
		}
		if typedEvent, ok := pEvent.(T); ok {
			filteredEvents = append(filteredEvents, typedEvent)
		}
	}
	return filteredEvents, len(filteredEvents) > 0
}

// EventOfType gets one event of a specific type
func EventOfType[T proto.Message](events []abci.Event) (T, bool) {
	var event T
	for _, ev := range events {
		pEvent, err := sdktypes.ParseTypedEvent(ev)
		if err != nil {
			continue
		}
		if typedEvent, ok := pEvent.(T); ok {
			return typedEvent, true
		}
	}
	return event, false
}
