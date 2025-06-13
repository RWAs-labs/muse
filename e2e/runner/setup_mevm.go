package runner

import (
	"math/big"
	"time"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	connectormevm "github.com/RWAs-labs/protocol-contracts/pkg/museconnectormevm.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/systemcontract.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/wmuse.sol"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/erc1967proxy"
	"github.com/RWAs-labs/muse/e2e/contracts/testdappv2"
	"github.com/RWAs-labs/muse/e2e/txserver"
	e2eutils "github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/contracts/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/RWAs-labs/muse/pkg/contracts/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// EmissionsPoolFunding represents the amount of MUSE to fund the emissions pool with
// This is the same value as used originally on mainnet (20M MUSE)
var EmissionsPoolFunding = big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(2e7))

// SetTSSAddresses set TSS addresses from information queried from MuseChain
func (r *E2ERunner) SetTSSAddresses() error {
	btcChainID, err := chains.GetBTCChainIDFromChainParams(r.BitcoinParams)
	if err != nil {
		return err
	}

	res := &observertypes.QueryGetTssAddressResponse{}
	for i := 0; ; i++ {
		res, err = r.ObserverClient.GetTssAddress(r.Ctx, &observertypes.QueryGetTssAddressRequest{
			BitcoinChainId: btcChainID,
		})
		if err != nil {
			if i%10 == 0 {
				r.Logger.Info("ObserverClient.TSS error %s", err.Error())
				r.Logger.Info("TSS not ready yet, waiting for TSS to be appear in musecore network...")
			}
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	tssAddress := ethcommon.HexToAddress(res.Eth)

	btcTSSAddress, err := btcutil.DecodeAddress(res.Btc, r.BitcoinParams)
	require.NoError(r, err)

	r.TSSAddress = tssAddress
	r.BTCTSSAddress = btcTSSAddress
	r.SuiTSSAddress = res.Sui

	return nil
}

// SetupMEVMMRC20s setup MRC20 for the MEVM
func (r *E2ERunner) SetupMEVMMRC20s(mrc20Deployment txserver.MRC20Deployment) {
	r.Logger.Print("⚙️ deploying MRC20s on MEVM")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("System contract deployments took %s\n", time.Since(startTime))
	}()

	// deploy system contracts and MRC20 contracts on MuseChain
	deployedMRC20Addresses, err := r.MuseTxServer.DeployMRC20s(
		mrc20Deployment,
		r.skipChainOperations,
	)
	require.NoError(r, err)

	// Set ERC20MRC20Addr
	r.ERC20MRC20Addr = deployedMRC20Addresses.ERC20MRC20Addr
	r.ERC20MRC20, err = mrc20.NewMRC20(r.ERC20MRC20Addr, r.MEVMClient)
	require.NoError(r, err)

	// Set SPLMRC20Addr if set
	if deployedMRC20Addresses.SPLMRC20Addr != (ethcommon.Address{}) {
		r.SPLMRC20Addr = deployedMRC20Addresses.SPLMRC20Addr
		r.SPLMRC20, err = mrc20.NewMRC20(r.SPLMRC20Addr, r.MEVMClient)
		require.NoError(r, err)
	}

	// set MRC20 contracts
	r.SetupETHMRC20()
	r.SetupBTCMRC20()
	r.SetupSOLMRC20()
	r.SetupTONMRC20()
}

// SetupETHMRC20 sets up the ETH MRC20 in the runner from the values queried from the chain
func (r *E2ERunner) SetupETHMRC20() {
	ethMRC20Addr, err := r.SystemContract.GasCoinMRC20ByChainId(
		&bind.CallOpts{},
		big.NewInt(chains.GoerliLocalnet.ChainId),
	)
	require.NoError(r, err)
	require.NotEqual(r, ethcommon.Address{}, ethMRC20Addr, "eth mrc20 not found")

	r.ETHMRC20Addr = ethMRC20Addr
	ethMRC20, err := mrc20.NewMRC20(ethMRC20Addr, r.MEVMClient)
	require.NoError(r, err)

	r.ETHMRC20 = ethMRC20
}

// SetupBTCMRC20 sets up the BTC MRC20 in the runner from the values queried from the chain
func (r *E2ERunner) SetupBTCMRC20() {
	BTCMRC20Addr, err := r.SystemContract.GasCoinMRC20ByChainId(
		&bind.CallOpts{},
		big.NewInt(chains.BitcoinRegtest.ChainId),
	)
	require.NoError(r, err)
	r.BTCMRC20Addr = BTCMRC20Addr
	r.Logger.Info("BTCMRC20Addr: %s", BTCMRC20Addr.Hex())
	BTCMRC20, err := mrc20.NewMRC20(BTCMRC20Addr, r.MEVMClient)
	require.NoError(r, err)
	r.BTCMRC20 = BTCMRC20
}

// SetupSOLMRC20 sets up the SOL MRC20 in the runner from the values queried from the chain
func (r *E2ERunner) SetupSOLMRC20() {
	// set SOLMRC20 address by chain ID
	SOLMRC20Addr, err := r.SystemContract.GasCoinMRC20ByChainId(
		&bind.CallOpts{},
		big.NewInt(chains.SolanaLocalnet.ChainId),
	)
	require.NoError(r, err)

	// set SOLMRC20 address
	r.SOLMRC20Addr = SOLMRC20Addr
	r.Logger.Info("SOLMRC20Addr: %s", SOLMRC20Addr.Hex())

	// set SOLMRC20 contract
	SOLMRC20, err := mrc20.NewMRC20(SOLMRC20Addr, r.MEVMClient)
	require.NoError(r, err)
	r.SOLMRC20 = SOLMRC20
}

// SetupTONMRC20 sets up the TON MRC20 in the runner from the values queried from the chain
func (r *E2ERunner) SetupTONMRC20() {
	chainID := chains.TONLocalnet.ChainId

	// noop
	if r.skipChainOperations(chainID) {
		return
	}

	TONMRC20Addr, err := r.SystemContract.GasCoinMRC20ByChainId(&bind.CallOpts{}, big.NewInt(chainID))
	require.NoError(r, err)

	r.TONMRC20Addr = TONMRC20Addr
	r.Logger.Info("TON MRC20 address: %s", TONMRC20Addr.Hex())

	TONMRC20, err := mrc20.NewMRC20(TONMRC20Addr, r.MEVMClient)
	require.NoError(r, err)

	r.TONMRC20 = TONMRC20
}

// SetupSUIMRC20 sets up the SUI MRC20 in the runner from the values queried from the chain
func (r *E2ERunner) SetupSUIMRC20() {
	chainID := chains.SuiLocalnet.ChainId

	// noop
	if r.skipChainOperations(chainID) {
		return
	}

	SUIMRC20Addr, err := r.SystemContract.GasCoinMRC20ByChainId(&bind.CallOpts{}, big.NewInt(chainID))
	require.NoError(r, err)

	r.SUIMRC20Addr = SUIMRC20Addr
	r.Logger.Info("SUI MRC20 address: %s", SUIMRC20Addr.Hex())

	SUIMRC20, err := mrc20.NewMRC20(SUIMRC20Addr, r.MEVMClient)
	require.NoError(r, err)

	r.SUIMRC20 = SUIMRC20
}

// EnableHeaderVerification enables the header verification for the given chain IDs
func (r *E2ERunner) EnableHeaderVerification(chainIDList []int64) error {
	r.Logger.Print("⚙️ enabling verification flags for block headers")

	return r.MuseTxServer.EnableHeaderVerification(e2eutils.AdminPolicyName, chainIDList)
}

// SetupMEVMProtocolContracts setup protocol contracts for the MEVM
func (r *E2ERunner) SetupMEVMProtocolContracts() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := e2eutils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage+" tx hash: "+tx.Hash().Hex())
	}

	r.Logger.Print("⚙️ setting up MEVM protocol contracts")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("MEVM protocol contracts took %s\n", time.Since(startTime))
	}()

	// deploy system contracts and MRC20 contracts on MuseChain
	addresses, err := r.MuseTxServer.DeploySystemContracts(
		e2eutils.OperationalPolicyName,
		e2eutils.AdminPolicyName,
	)
	require.NoError(r, err)

	// UniswapV2FactoryAddr
	r.UniswapV2FactoryAddr = ethcommon.HexToAddress(addresses.UniswapV2FactoryAddr)
	r.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(r.UniswapV2FactoryAddr, r.MEVMClient)
	require.NoError(r, err)

	// UniswapV2RouterAddr
	r.UniswapV2RouterAddr = ethcommon.HexToAddress(addresses.UniswapV2RouterAddr)
	r.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(r.UniswapV2RouterAddr, r.MEVMClient)
	require.NoError(r, err)

	// MevmConnectorAddr
	r.ConnectorMEVMAddr = ethcommon.HexToAddress(addresses.MEVMConnectorAddr)
	r.ConnectorMEVM, err = connectormevm.NewMuseConnectorMEVM(r.ConnectorMEVMAddr, r.MEVMClient)
	require.NoError(r, err)

	// WMuseAddr
	r.WMuseAddr = ethcommon.HexToAddress(addresses.WMUSEAddr)
	r.WMuse, err = wmuse.NewWETH9(r.WMuseAddr, r.MEVMClient)
	require.NoError(r, err)

	// query system contract address from the chain
	systemContractRes, err := r.FungibleClient.SystemContract(
		r.Ctx,
		&fungibletypes.QueryGetSystemContractRequest{},
	)
	require.NoError(r, err)

	systemContractAddr := ethcommon.HexToAddress(systemContractRes.SystemContract.SystemContract)
	systemContract, err := systemcontract.NewSystemContract(
		systemContractAddr,
		r.MEVMClient,
	)
	require.NoError(r, err)

	r.SystemContract = systemContract
	r.SystemContractAddr = systemContractAddr

	r.Logger.Info("Deploying Gateway MEVM")
	gatewayMEVMAddr, txGateway, _, err := gatewaymevm.DeployGatewayMEVM(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err)

	ensureTxReceipt(txGateway, "Gateway deployment failed")

	gatewayMEVMABI, err := gatewaymevm.GatewayMEVMMetaData.GetAbi()
	require.NoError(r, err)

	// Encode the initializer data
	initializerData, err := gatewayMEVMABI.Pack("initialize", r.WMuseAddr, r.Account.EVMAddress())
	require.NoError(r, err)

	// Deploy the proxy contract
	r.Logger.Info(
		"Deploying proxy with %s and %s, address: %s",
		r.WMuseAddr.Hex(),
		r.Account.EVMAddress().Hex(),
		gatewayMEVMAddr.Hex(),
	)
	proxyAddress, txProxy, _, err := erc1967proxy.DeployERC1967Proxy(
		r.MEVMAuth,
		r.MEVMClient,
		gatewayMEVMAddr,
		initializerData,
	)
	require.NoError(r, err)

	r.GatewayMEVMAddr = proxyAddress
	r.GatewayMEVM, err = gatewaymevm.NewGatewayMEVM(proxyAddress, r.MEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Gateway MEVM contract address: %s, tx hash: %s", gatewayMEVMAddr.Hex(), txGateway.Hash().Hex())

	// Set the gateway address in the protocol
	err = r.MuseTxServer.UpdateGatewayAddress(e2eutils.AdminPolicyName, r.GatewayMEVMAddr.Hex())
	require.NoError(r, err)

	// deploy test dapp v2
	testDAppV2Addr, txTestDAppV2, _, err := testdappv2.DeployTestDAppV2(
		r.MEVMAuth,
		r.MEVMClient,
		true,
		r.GatewayEVMAddr,
	)
	require.NoError(r, err)

	r.TestDAppV2MEVMAddr = testDAppV2Addr
	r.TestDAppV2MEVM, err = testdappv2.NewTestDAppV2(testDAppV2Addr, r.MEVMClient)
	require.NoError(r, err)

	ensureTxReceipt(txProxy, "Gateway proxy deployment failed")
	ensureTxReceipt(txTestDAppV2, "TestDAppV2 deployment failed")

	// check isMuseChain is true
	isMuseChain, err := r.TestDAppV2MEVM.IsMuseChain(&bind.CallOpts{})
	require.NoError(r, err)
	require.True(r, isMuseChain)
}

// UpdateProtocolContractsInChainParams update the erc20 custody contract and gateway address in the chain params
// TODO: should be used for all protocol contracts including the MUSE connector
// https://github.com/RWAs-labs/muse/issues/3257
func (r *E2ERunner) UpdateProtocolContractsInChainParams() {
	res, err := r.ObserverClient.GetChainParams(r.Ctx, &observertypes.QueryGetChainParamsRequest{})
	require.NoError(r, err)

	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	// find old chain params
	var (
		chainParams *observertypes.ChainParams
		found       bool
	)
	for _, cp := range res.ChainParams.ChainParams {
		if cp.ChainId == evmChainID.Int64() {
			chainParams = cp
			found = true
			break
		}
	}
	require.True(r, found, "Chain params not found for chain id %d", evmChainID)

	// update with the new ERC20 custody contract address
	chainParams.Erc20CustodyContractAddress = r.ERC20CustodyAddr.Hex()

	// update with the new gateway address
	chainParams.GatewayAddress = r.GatewayEVMAddr.Hex()

	// update the chain params
	err = r.MuseTxServer.UpdateChainParams(chainParams)
	require.NoError(r, err)
}
