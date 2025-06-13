package runner

import (
	"math/big"
	"time"

	museeth "github.com/RWAs-labs/protocol-contracts/pkg/muse.eth.sol"
	museconnectoreth "github.com/RWAs-labs/protocol-contracts/pkg/museconnector.eth.sol"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/config"
	"github.com/RWAs-labs/muse/e2e/contracts/testdapp"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/constant"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

const (
	ContractsConfigFile = "contracts.toml"
)

// LegacySetEVMContractsFromConfig set legacy EVM contracts for e2e test from the config
func (r *E2ERunner) LegacySetEVMContractsFromConfig() {
	conf, err := config.ReadConfig(ContractsConfigFile, true)
	require.NoError(r, err)

	// Set MuseEthAddr
	r.MuseEthAddr = ethcommon.HexToAddress(conf.Contracts.EVM.MuseEthAddr.String())
	r.MuseEth, err = museeth.NewMuseEth(r.MuseEthAddr, r.EVMClient)
	require.NoError(r, err)

	// Set ConnectorEthAddr
	r.ConnectorEthAddr = ethcommon.HexToAddress(conf.Contracts.EVM.ConnectorEthAddr.String())
	r.ConnectorEth, err = museconnectoreth.NewMuseConnectorEth(r.ConnectorEthAddr, r.EVMClient)
	require.NoError(r, err)
}

// LegacySetupEVM setup legacy contracts on EVM for e2e test
func (r *E2ERunner) LegacySetupEVM(contractsDeployed bool, legacyTestRunning bool) {
	r.Logger.Print("⚙️ setting up EVM network legacy contracts")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("EVM setup took %s\n", time.Since(startTime))
	}()

	// We use this config to be consistent with the old implementation
	if contractsDeployed {
		r.LegacySetEVMContractsFromConfig()
		return
	}
	conf := config.DefaultConfig()

	r.Logger.InfoLoud("Deploy MuseETH ConnectorETH ERC20Custody ERC20\n")

	// donate to the TSS address to avoid account errors because deploying gas token MRC20 will automatically mint
	// gas token on MuseChain to initialize the pool
	txDonation, err := r.LegacySendEther(r.TSSAddress, big.NewInt(101000000000000000), []byte(constant.DonationMessage))
	require.NoError(r, err)

	r.Logger.Info("Deploying MuseEth contract")
	museEthAddr, txMuseEth, MuseEth, err := museeth.DeployMuseEth(
		r.EVMAuth,
		r.EVMClient,
		r.EVMAddress(),
		big.NewInt(21_000_000_000),
	)
	require.NoError(r, err)

	r.MuseEth = MuseEth
	r.MuseEthAddr = museEthAddr
	conf.Contracts.EVM.MuseEthAddr = config.DoubleQuotedString(museEthAddr.String())
	r.Logger.Info("MuseEth contract address: %s, tx hash: %s", museEthAddr.Hex(), txMuseEth.Hash())

	r.Logger.Info("Deploying MuseConnectorEth contract")
	connectorEthAddr, txConnector, ConnectorEth, err := museconnectoreth.DeployMuseConnectorEth(
		r.EVMAuth,
		r.EVMClient,
		museEthAddr,
		r.TSSAddress,
		r.EVMAddress(),
		r.EVMAddress(),
	)
	require.NoError(r, err)

	r.ConnectorEth = ConnectorEth
	r.ConnectorEthAddr = connectorEthAddr
	conf.Contracts.EVM.ConnectorEthAddr = config.DoubleQuotedString(connectorEthAddr.String())

	r.Logger.Info(
		"MuseConnectorEth contract address: %s, tx hash: %s",
		connectorEthAddr.Hex(),
		txConnector.Hash().Hex(),
	)

	// deploy TestDApp contract
	appAddr, txApp, _, err := testdapp.DeployTestDApp(
		r.EVMAuth,
		r.EVMClient,
		r.ConnectorEthAddr,
		r.MuseEthAddr,
	)
	require.NoError(r, err)

	r.EvmTestDAppAddr = appAddr
	r.Logger.Info("TestDApp contract address: %s, tx hash: %s", appAddr.Hex(), txApp.Hash().Hex())

	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage)
	}

	// check contract deployment receipt
	ensureTxReceipt(txDonation, "EVM donation tx failed")
	ensureTxReceipt(txMuseEth, "MuseEth deployment failed")
	ensureTxReceipt(txConnector, "MuseConnectorEth deployment failed")
	ensureTxReceipt(txApp, "TestDApp deployment failed")

	// save config containing contract addresses
	// TODO: put this logic outside of this function in a general config
	// We use this config to be consistent with the old implementation
	// https://github.com/RWAs-labs/muse-private/issues/41
	require.NoError(r, config.WriteConfig(ContractsConfigFile, conf))

	// chain params will need to be updated if they do not match the default params
	// this be required if the deployer account changes
	currentChainParamsRes, err := r.ObserverClient.GetChainParamsForChain(
		r.Ctx,
		&observertypes.QueryGetChainParamsForChainRequest{
			ChainId: chains.GoerliLocalnet.ChainId,
		},
	)
	require.NoError(r, err, "failed to get chain params for chain %d", chains.GoerliLocalnet.ChainId)

	chainParams := currentChainParamsRes.ChainParams
	chainParams.ConnectorContractAddress = r.ConnectorEthAddr.Hex()
	chainParams.MuseTokenContractAddress = r.MuseEthAddr.Hex()
	if legacyTestRunning {
		chainParams.DisableTssBlockScan = false
	}

	err = r.MuseTxServer.UpdateChainParams(chainParams)
	require.NoError(r, err, "failed to update chain params")
	r.Logger.Print("🔄 updated chain params")
}
