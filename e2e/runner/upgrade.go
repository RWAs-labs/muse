package runner

import (
	"github.com/RWAs-labs/protocol-contracts/pkg/erc20custody.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/config"
	"github.com/RWAs-labs/muse/e2e/utils"
)

// UpgradeGatewayOptions is the options for the gateway upgrade tests
type UpgradeGatewayOptions struct {
	TestSolana bool
	TestSui    bool
}

// UpgradeGatewaysAndERC20Custody upgrades gateways and ERC20Custody contracts
// It deploys new contract implementation with the current imported artifacts and upgrades the contract
func (r *E2ERunner) UpgradeGatewaysAndERC20Custody() {
	r.UpgradeGatewayMEVM()
	r.UpgradeGatewayEVM()
	r.UpgradeERC20Custody()
}

// RunGatewayUpgradeTestsExternalChains runs the gateway upgrade tests for external chains
func (r *E2ERunner) RunGatewayUpgradeTestsExternalChains(conf config.Config, opts UpgradeGatewayOptions) {
	if opts.TestSolana {
		r.SolanaVerifyGatewayContractsUpgrade(conf.AdditionalAccounts.UserSolana.SolanaPrivateKey.String())
	}

	if opts.TestSui {
		r.SuiVerifyGatewayPackageUpgrade()
	}
}

// UpgradeGatewayMEVM upgrades the GatewayMEVM contract
func (r *E2ERunner) UpgradeGatewayMEVM() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage+" tx hash: "+tx.Hash().Hex())
	}

	r.Logger.Info("Upgrading Gateway MEVM contract")
	// Deploy the new gateway contract implementation
	newImplementationAddress, txDeploy, _, err := gatewaymevm.DeployGatewayMEVM(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err)
	ensureTxReceipt(txDeploy, "New GatewayMEVM implementation deployment failed")

	// Upgrade
	txUpgrade, err := r.GatewayMEVM.UpgradeToAndCall(r.MEVMAuth, newImplementationAddress, []byte{})
	require.NoError(r, err)
	ensureTxReceipt(txUpgrade, "GatewayMEVM upgrade failed")
}

// UpgradeGatewayEVM upgrades the GatewayEVM contract
func (r *E2ERunner) UpgradeGatewayEVM() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage+" tx hash: "+tx.Hash().Hex())
	}

	r.Logger.Info("Upgrading Gateway EVM contract")
	// Deploy the new gateway contract implementation
	newImplementationAddress, txDeploy, _, err := gatewayevm.DeployGatewayEVM(r.EVMAuth, r.EVMClient)
	require.NoError(r, err)
	ensureTxReceipt(txDeploy, "New GatewayEVM implementation deployment failed")

	// Upgrade
	txUpgrade, err := r.GatewayEVM.UpgradeToAndCall(r.EVMAuth, newImplementationAddress, []byte{})
	require.NoError(r, err)
	ensureTxReceipt(txUpgrade, "GatewayEVM upgrade failed")
}

// UpgradeERC20Custody upgrades the ERC20Custody contract
func (r *E2ERunner) UpgradeERC20Custody() {
	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.requireTxSuccessful(receipt, failMessage+" tx hash: "+tx.Hash().Hex())
	}

	r.Logger.Info("Upgrading ERC20Custody contract")
	// Deploy the new erc20Custody contract implementation
	newImplementationAddress, txDeploy, _, err := erc20custody.DeployERC20Custody(r.EVMAuth, r.EVMClient)
	require.NoError(r, err)
	ensureTxReceipt(txDeploy, "New ERC20Custody implementation deployment failed")

	// Upgrade
	txUpgrade, err := r.ERC20Custody.UpgradeToAndCall(r.EVMAuth, newImplementationAddress, []byte{})
	require.NoError(r, err)
	ensureTxReceipt(txUpgrade, "ERC20Custody upgrade failed")
}
