package config

import (
	"context"
	"fmt"

	"github.com/RWAs-labs/muse/e2e/config"
	"github.com/RWAs-labs/muse/e2e/runner"
)

// RunnerFromConfig create test runner from config
func RunnerFromConfig(
	ctx context.Context,
	name string,
	ctxCancel context.CancelCauseFunc,
	conf config.Config,
	account config.Account,
	logger *runner.Logger,
	opts ...runner.E2ERunnerOption,
) (*runner.E2ERunner, error) {
	// initialize all clients for E2E tests
	e2eClients, err := getClientsFromConfig(ctx, conf, account)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients from config: %w", err)
	}

	// initialize E2E test runner
	newRunner := runner.NewE2ERunner(
		ctx,
		name,
		ctxCancel,
		account,
		e2eClients,
		logger,
		opts...,
	)

	// set contracts
	err = setContractsFromConfig(newRunner, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to set contracts from config: %w", err)
	}

	// set bitcoin params
	chainParams, err := conf.RPCs.Bitcoin.Params.GetParams()
	if err != nil {
		return nil, fmt.Errorf("failed to get bitcoin params: %w", err)
	}
	newRunner.BitcoinParams = &chainParams

	return newRunner, err
}

// ExportContractsFromRunner export contracts from the runner to config using a source config
func ExportContractsFromRunner(r *runner.E2ERunner, conf config.Config) config.Config {
	// copy contracts from deployer runner
	conf.Contracts.Solana.GatewayProgramID = config.DoubleQuotedString(r.GatewayProgram.String())
	conf.Contracts.Solana.SPLAddr = config.DoubleQuotedString(r.SPLAddr.String())

	conf.Contracts.TON.GatewayAccountID = config.DoubleQuotedString(r.TONGateway.ToRaw())

	if r.SuiGateway != nil {
		conf.Contracts.Sui.GatewayPackageID = config.DoubleQuotedString(r.SuiGateway.PackageID())
		conf.Contracts.Sui.GatewayObjectID = config.DoubleQuotedString(r.SuiGateway.ObjectID())
	}
	conf.Contracts.Sui.GatewayUpgradeCap = config.DoubleQuotedString(r.SuiGatewayUpgradeCap)

	conf.Contracts.Sui.FungibleTokenCoinType = config.DoubleQuotedString(r.SuiTokenCoinType)
	conf.Contracts.Sui.FungibleTokenTreasuryCap = config.DoubleQuotedString(r.SuiTokenTreasuryCap)
	conf.Contracts.Sui.Example = r.SuiExample

	conf.Contracts.EVM.MuseEthAddr = config.DoubleQuotedString(r.MuseEthAddr.Hex())
	conf.Contracts.EVM.ConnectorEthAddr = config.DoubleQuotedString(r.ConnectorEthAddr.Hex())
	conf.Contracts.EVM.CustodyAddr = config.DoubleQuotedString(r.ERC20CustodyAddr.Hex())
	conf.Contracts.EVM.ERC20 = config.DoubleQuotedString(r.ERC20Addr.Hex())
	conf.Contracts.EVM.TestDappAddr = config.DoubleQuotedString(r.EvmTestDAppAddr.Hex())
	conf.Contracts.EVM.Gateway = config.DoubleQuotedString(r.GatewayEVMAddr.Hex())
	conf.Contracts.EVM.TestDAppV2Addr = config.DoubleQuotedString(r.TestDAppV2EVMAddr.Hex())

	conf.Contracts.MEVM.SystemContractAddr = config.DoubleQuotedString(r.SystemContractAddr.Hex())
	conf.Contracts.MEVM.ETHMRC20Addr = config.DoubleQuotedString(r.ETHMRC20Addr.Hex())
	conf.Contracts.MEVM.ERC20MRC20Addr = config.DoubleQuotedString(r.ERC20MRC20Addr.Hex())
	conf.Contracts.MEVM.BTCMRC20Addr = config.DoubleQuotedString(r.BTCMRC20Addr.Hex())
	conf.Contracts.MEVM.SOLMRC20Addr = config.DoubleQuotedString(r.SOLMRC20Addr.Hex())
	conf.Contracts.MEVM.SPLMRC20Addr = config.DoubleQuotedString(r.SPLMRC20Addr.Hex())
	conf.Contracts.MEVM.TONMRC20Addr = config.DoubleQuotedString(r.TONMRC20Addr.Hex())
	conf.Contracts.MEVM.SUIMRC20Addr = config.DoubleQuotedString(r.SUIMRC20Addr.Hex())
	conf.Contracts.MEVM.SuiTokenMRC20Addr = config.DoubleQuotedString(r.SuiTokenMRC20Addr.Hex())

	conf.Contracts.MEVM.UniswapFactoryAddr = config.DoubleQuotedString(r.UniswapV2FactoryAddr.Hex())
	conf.Contracts.MEVM.UniswapRouterAddr = config.DoubleQuotedString(r.UniswapV2RouterAddr.Hex())
	conf.Contracts.MEVM.ConnectorMEVMAddr = config.DoubleQuotedString(r.ConnectorMEVMAddr.Hex())
	conf.Contracts.MEVM.WMuseAddr = config.DoubleQuotedString(r.WMuseAddr.Hex())
	conf.Contracts.MEVM.MEVMSwapAppAddr = config.DoubleQuotedString(r.MEVMSwapAppAddr.Hex())
	conf.Contracts.MEVM.ContextAppAddr = config.DoubleQuotedString(r.ContextAppAddr.Hex())
	conf.Contracts.MEVM.TestDappAddr = config.DoubleQuotedString(r.MevmTestDAppAddr.Hex())
	conf.Contracts.MEVM.Gateway = config.DoubleQuotedString(r.GatewayMEVMAddr.Hex())
	conf.Contracts.MEVM.TestDAppV2Addr = config.DoubleQuotedString(r.TestDAppV2MEVMAddr.Hex())

	return conf
}
