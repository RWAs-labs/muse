package e2etests

import (
	"context"
	"fmt"
	"math/big"
	"os/exec"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/config"
	"github.com/RWAs-labs/muse/e2e/e2etests/legacy"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/txserver"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// EVM2RPCURL is the RPC URL for the additional EVM localnet
// Only this test currently uses a additional EVM localnet, and this test is only run locally
// Therefore, we hardcode RPC urls and addresses for simplicity
const EVM2RPCURL = "http://eth2:8545"

// EVMSepoliaChainID is the chain ID for the additional EVM localnet
// We set Sepolia testnet although the value is not important, only used to differentiate
var EVMSepoliaChainID = chains.Sepolia.ChainId

func TestMigrateChainSupport(r *runner.E2ERunner, _ []string) {
	// deposit most of the MUSE supply on MuseChain
	museAmount := big.NewInt(1e18)
	museAmount = museAmount.Mul(museAmount, big.NewInt(20_000_000_000)) // 20B Muse
	r.LegacyDepositMuseWithAmount(r.EVMAddress(), museAmount)

	// do an ethers withdraw on the previous chain (0.01eth) for some interaction
	legacy.TestEtherWithdraw(r, []string{"10000000000000000"})

	// create runner for the new EVM and set it up
	newRunner, err := configureEVM2(r)
	require.NoError(r, err)

	newRunner.LegacySetupEVM(false, false)

	// mint some ERC20
	newRunner.MintERC20OnEVM(10000)

	// we deploy connectorETH in this test to simulate a new "canonical" chain emitting MUSE
	// to represent the MUSE already existing on MuseChain we manually send the minted MUSE to the connector
	newRunner.LegacySendMuseOnEvm(newRunner.ConnectorEthAddr, 20_000_000_000)

	// update the chain params to set up the chain
	chainParams := getNewEVMChainParams(newRunner)

	err = r.MuseTxServer.UpdateChainParams(chainParams)
	require.NoError(r, err)

	// setup the gas token
	require.NoError(r, err)
	_, err = newRunner.MuseTxServer.BroadcastTx(
		utils.AdminPolicyName,
		fungibletypes.NewMsgDeployFungibleCoinMRC20(
			r.MuseTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
			"",
			chainParams.ChainId,
			18,
			"Sepolia ETH",
			"sETH",
			coin.CoinType_Gas,
			100000,
			nil,
		),
	)
	require.NoError(r, err)

	// set the gas token in the runner
	ethMRC20Addr, err := newRunner.SystemContract.GasCoinMRC20ByChainId(
		&bind.CallOpts{},
		big.NewInt(chainParams.ChainId),
	)
	require.NoError(r, err)
	require.NotEqual(r, ethcommon.Address{}, ethMRC20Addr)

	newRunner.ETHMRC20Addr = ethMRC20Addr
	ethMRC20, err := mrc20.NewMRC20(ethMRC20Addr, newRunner.MEVMClient)
	require.NoError(r, err)
	newRunner.ETHMRC20 = ethMRC20

	// set the chain nonces for the new chain
	_, err = r.MuseTxServer.BroadcastTx(utils.OperationalPolicyName, observertypes.NewMsgResetChainNonces(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		chainParams.ChainId,
		0,
		0,
	))
	require.NoError(r, err)

	// deactivate the previous chain
	chainParams = observertypes.GetDefaultGoerliLocalnetChainParams()
	chainParams.IsSupported = false
	err = r.MuseTxServer.UpdateChainParams(chainParams)
	require.NoError(r, err)

	// restart MuseClient to pick up the new chain
	r.Logger.Print("🔄 restarting MuseClient to pick up the new chain")
	require.NoError(r, restartMuseClient())

	// wait 10 set for the chain to start
	time.Sleep(10 * time.Second)

	// emitting a withdraw with the previous chain should fail
	txWithdraw, err := r.ETHMRC20.Withdraw(r.MEVMAuth, r.EVMAddress().Bytes(), big.NewInt(10000000000000000))
	if err == nil {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, txWithdraw, r.Logger, r.ReceiptTimeout)
		utils.RequiredTxFailed(r, receipt)
	}

	// test cross-chain functionalities on the new network
	// we use a Go routine to manually mine blocks because Anvil network only mine blocks on tx by default
	// we need automatic block mining to get the necessary confirmations for the cross-chain functionalities
	stopMining, err := newRunner.AnvilMineBlocks(EVM2RPCURL, 3)
	require.NoError(r, err)

	// deposit Ethers and ERC20 on MuseChain
	etherAmount := big.NewInt(1e18)
	etherAmount = etherAmount.Mul(etherAmount, big.NewInt(10))
	txEtherDeposit := newRunner.LegacyDepositEtherWithAmount(etherAmount)
	newRunner.WaitForMinedCCTX(txEtherDeposit)

	// perform withdrawals on the new chain
	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(10))
	newRunner.LegacyDepositAndApproveWMuse(amount)
	tx := newRunner.LegacyWithdrawMuse(amount, true)
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "muse withdraw")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	legacy.TestEtherWithdraw(newRunner, []string{"50000000000000000"})

	// finally try to deposit Muse back
	legacy.TestMuseDeposit(newRunner, []string{"100000000000000000"})

	// ERC20 test

	// whitelist erc20 mrc20
	newRunner.Logger.Info("whitelisting ERC20 on new network")
	res, err := newRunner.MuseTxServer.BroadcastTx(utils.AdminPolicyName, crosschaintypes.NewMsgWhitelistERC20(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		newRunner.ERC20Addr.Hex(),
		chains.Sepolia.ChainId,
		"USDT",
		"USDT",
		18,
		100000,
		sdkmath.NewUintFromString("100000000000000000000000000"),
	))
	require.NoError(r, err)

	event, ok := txserver.EventOfType[*crosschaintypes.EventERC20Whitelist](res.Events)
	require.True(r, ok, "no EventERC20Whitelist in %s", res.TxHash)
	erc20mrc20Addr := event.Mrc20Address
	whitelistCCTXIndex := event.WhitelistCctxIndex

	// wait for the whitelist cctx to be mined
	newRunner.WaitForMinedCCTXFromIndex(whitelistCCTXIndex)

	// set erc20 mrc20 contract address
	require.True(r, ethcommon.IsHexAddress(erc20mrc20Addr), "invalid contract address: %s", erc20mrc20Addr)

	erc20MRC20, err := mrc20.NewMRC20(ethcommon.HexToAddress(erc20mrc20Addr), newRunner.MEVMClient)
	require.NoError(r, err)

	newRunner.ERC20MRC20 = erc20MRC20

	// deposit ERC20 on MuseChain
	txERC20Deposit := newRunner.DepositERC20Deployer()
	newRunner.WaitForMinedCCTX(txERC20Deposit)

	// stop mining
	stopMining()
}

// configureEVM2 takes a runner and configures it to use the additional EVM localnet
func configureEVM2(r *runner.E2ERunner) (*runner.E2ERunner, error) {
	// initialize a new runner with previous runner values
	newRunner := runner.NewE2ERunner(
		r.Ctx,
		"admin-evm2",
		r.CtxCancel,
		r.Account,
		r.Clients,
		runner.NewLogger(true, color.FgHiYellow, "admin-evm2"),
		runner.WithMuseTxServer(r.MuseTxServer),
	)

	// All existing fields of the runner are the same except for the RPC URL and client for EVM
	ewvmClient, evmAuth, err := getEVMClient(newRunner.Ctx, EVM2RPCURL, r.Account)
	if err != nil {
		return nil, err
	}
	newRunner.EVMClient = ewvmClient
	newRunner.EVMAuth = evmAuth

	// Copy the MuseChain contract addresses from the original runner
	if err := newRunner.CopyAddressesFrom(r); err != nil {
		return nil, err
	}

	// reset evm contracts to ensure they are re-initialized
	newRunner.MuseEthAddr = ethcommon.Address{}
	newRunner.MuseEth = nil
	newRunner.ConnectorEthAddr = ethcommon.Address{}
	newRunner.ConnectorEth = nil
	newRunner.ERC20CustodyAddr = ethcommon.Address{}
	newRunner.ERC20Custody = nil
	newRunner.ERC20Addr = ethcommon.Address{}
	newRunner.ERC20 = nil

	return newRunner, nil
}

// getEVMClient get evm client from rpc and private key
func getEVMClient(
	ctx context.Context,
	rpc string,
	account config.Account,
) (*ethclient.Client, *bind.TransactOpts, error) {
	evmClient, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, nil, err
	}

	chainid, err := evmClient.ChainID(ctx)
	if err != nil {
		return nil, nil, err
	}
	deployerPrivkey, err := account.PrivateKey()
	if err != nil {
		return nil, nil, err
	}
	evmAuth, err := bind.NewKeyedTransactorWithChainID(deployerPrivkey, chainid)
	if err != nil {
		return nil, nil, err
	}

	return evmClient, evmAuth, nil
}

// getNewEVMChainParams returns the chain params for the new EVM chain
func getNewEVMChainParams(r *runner.E2ERunner) *observertypes.ChainParams {
	// goerli local as base
	chainParams := observertypes.GetDefaultGoerliLocalnetChainParams()

	// set the chain id to the new chain id
	chainParams.ChainId = EVMSepoliaChainID

	// set contracts
	chainParams.ConnectorContractAddress = r.ConnectorEthAddr.Hex()
	chainParams.Erc20CustodyContractAddress = r.ERC20CustodyAddr.Hex()
	chainParams.MuseTokenContractAddress = r.MuseEthAddr.Hex()

	// set supported
	chainParams.IsSupported = true

	return chainParams
}

// restartMuseClient restarts the Muse client
func restartMuseClient() error {
	sshCommandFilePath := "/work/restart-museclientd.sh"
	cmd := exec.Command("/bin/sh", sshCommandFilePath)

	// Execute the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error restarting MuseClient: %s - %s", err.Error(), output)
	}
	return nil
}
