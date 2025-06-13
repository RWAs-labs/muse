package config

import (
	"fmt"

	"github.com/RWAs-labs/protocol-contracts/pkg/erc20custody.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	museeth "github.com/RWAs-labs/protocol-contracts/pkg/muse.eth.sol"
	museconnectoreth "github.com/RWAs-labs/protocol-contracts/pkg/museconnector.eth.sol"
	connectormevm "github.com/RWAs-labs/protocol-contracts/pkg/museconnectormevm.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/systemcontract.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/wmuse.sol"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/ton"

	"github.com/RWAs-labs/muse/e2e/config"
	"github.com/RWAs-labs/muse/e2e/contracts/contextapp"
	"github.com/RWAs-labs/muse/e2e/contracts/erc20"
	"github.com/RWAs-labs/muse/e2e/contracts/mevmswap"
	"github.com/RWAs-labs/muse/e2e/contracts/testdappv2"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/contracts/sui"
	"github.com/RWAs-labs/muse/pkg/contracts/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/RWAs-labs/muse/pkg/contracts/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	"github.com/RWAs-labs/muse/x/observer/types"
)

func chainParamsBySelector(
	chainParams []*types.ChainParams,
	selector func(chainID int64, additionalChains []chains.Chain) bool,
) *types.ChainParams {
	for _, chainParam := range chainParams {
		if selector(chainParam.ChainId, nil) {
			return chainParam
		}
	}
	return nil
}

func chainParamsByChainID(chainParams []*types.ChainParams, id int64) *types.ChainParams {
	for _, chainParam := range chainParams {
		if chainParam.ChainId == id {
			return chainParam
		}
	}
	return nil
}

func foreignCoinByChainID(
	foreignCoins []fungibletypes.ForeignCoins,
	id int64,
	coinType coin.CoinType,
) *fungibletypes.ForeignCoins {
	for _, fCoin := range foreignCoins {
		if fCoin.ForeignChainId == id && fCoin.CoinType == coinType {
			return &fCoin
		}
	}
	return nil
}

func setContractsGatewayEVM(r *runner.E2ERunner, params *types.ChainParams) error {
	r.GatewayEVMAddr = common.HexToAddress(params.GatewayAddress)
	if r.GatewayEVMAddr == (common.Address{}) {
		return nil
	}
	gatewayCode, err := r.EVMClient.CodeAt(r.Ctx, r.GatewayEVMAddr, nil)
	if err != nil || len(gatewayCode) == 0 {
		r.Logger.Print("❓ no code at EVM gateway address (%s)", r.GatewayEVMAddr)
		return nil
	}
	r.GatewayEVM, err = gatewayevm.NewGatewayEVM(r.GatewayEVMAddr, r.EVMClient)
	if err != nil {
		return err
	}
	r.MuseEthAddr = common.HexToAddress(params.MuseTokenContractAddress)
	r.MuseEth, err = museeth.NewMuseEth(r.MuseEthAddr, r.EVMClient)
	if err != nil {
		return err
	}

	r.ConnectorEthAddr = common.HexToAddress(params.ConnectorContractAddress)
	r.ConnectorEth, err = museconnectoreth.NewMuseConnectorEth(r.ConnectorEthAddr, r.EVMClient)
	if err != nil {
		return err
	}
	r.ERC20CustodyAddr = common.HexToAddress(params.Erc20CustodyContractAddress)
	r.ERC20Custody, err = erc20custody.NewERC20Custody(r.ERC20CustodyAddr, r.EVMClient)
	if err != nil {
		return err
	}
	return nil
}

// setContractsFromConfig get EVM contracts from config
func setContractsFromConfig(r *runner.E2ERunner, conf config.Config) error {
	var err error

	chainParams, err := r.Clients.Musecore.GetChainParams(r.Ctx)
	require.NoError(r, err, "get chain params")

	solChainParams := chainParamsBySelector(chainParams, chains.IsSolanaChain)

	// set Solana contracts
	if c := conf.Contracts.Solana.GatewayProgramID; c != "" {
		r.GatewayProgram = solana.MustPublicKeyFromBase58(c.String())
	} else if solChainParams != nil && solChainParams.GatewayAddress != "" {
		r.GatewayProgram = solana.MustPublicKeyFromBase58(solChainParams.GatewayAddress)
	}

	if c := conf.Contracts.Solana.SPLAddr; c != "" {
		r.SPLAddr = solana.MustPublicKeyFromBase58(c.String())
	}

	if c := conf.Contracts.TON.GatewayAccountID; c != "" {
		r.TONGateway = ton.MustParseAccountID(c.String())
	}

	// set Sui contracts
	suiPackageID := conf.Contracts.Sui.GatewayPackageID
	suiGatewayID := conf.Contracts.Sui.GatewayObjectID

	if suiPackageID != "" && suiGatewayID != "" {
		r.SuiGateway = sui.NewGateway(suiPackageID.String(), suiGatewayID.String())
	}
	if c := conf.Contracts.Sui.GatewayUpgradeCap; c != "" {
		r.SuiGatewayUpgradeCap = c.String()
	}
	if c := conf.Contracts.Sui.FungibleTokenCoinType; c != "" {
		r.SuiTokenCoinType = c.String()
	}
	if c := conf.Contracts.Sui.FungibleTokenTreasuryCap; c != "" {
		r.SuiTokenTreasuryCap = c.String()
	}
	r.SuiExample = conf.Contracts.Sui.Example

	// set EVM contracts
	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err, "get evm chain ID")
	evmChainParams := chainParamsByChainID(chainParams, evmChainID.Int64())

	if evmChainParams == nil {
		return fmt.Errorf("no EVM chain params found for chain ID %d", evmChainID.Int64())
	}

	err = setContractsGatewayEVM(r, evmChainParams)
	if err != nil {
		return fmt.Errorf("setContractsGatewayEVM: %w", err)
	}

	if c := conf.Contracts.EVM.ERC20; c != "" {
		r.ERC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ERC20: %w", err)
		}
		r.ERC20, err = erc20.NewERC20(r.ERC20Addr, r.EVMClient)
		if err != nil {
			return err
		}
	}

	// set MEVM contracts
	foreignCoins, err := r.Clients.Musecore.Fungible.ForeignCoinsAll(
		r.Ctx,
		&fungibletypes.QueryAllForeignCoinsRequest{},
	)
	if err != nil {
		return err
	}

	ethForeignCoin := foreignCoinByChainID(foreignCoins.ForeignCoins, evmChainID.Int64(), coin.CoinType_Gas)
	if ethForeignCoin != nil {
		r.ETHMRC20Addr = common.HexToAddress(ethForeignCoin.Mrc20ContractAddress)
		r.ETHMRC20, err = mrc20.NewMRC20(r.ETHMRC20Addr, r.MEVMClient)
		if err != nil {
			return err
		}
	}
	if c := conf.Contracts.MEVM.SystemContractAddr; c != "" {
		r.SystemContractAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid SystemContractAddr: %w", err)
		}
		r.SystemContract, err = systemcontract.NewSystemContract(r.SystemContractAddr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.ERC20MRC20Addr; c != "" {
		r.ERC20MRC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ERC20MRC20Addr: %w", err)
		}
		r.ERC20MRC20, err = mrc20.NewMRC20(r.ERC20MRC20Addr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.BTCMRC20Addr; c != "" {
		r.BTCMRC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid BTCMRC20Addr: %w", err)
		}
		r.BTCMRC20, err = mrc20.NewMRC20(r.BTCMRC20Addr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.SOLMRC20Addr; c != "" {
		r.SOLMRC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid SOLMRC20Addr: %w", err)
		}
		r.SOLMRC20, err = mrc20.NewMRC20(r.SOLMRC20Addr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.SPLMRC20Addr; c != "" {
		r.SPLMRC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid SPLMRC20Addr: %w", err)
		}
		r.SPLMRC20, err = mrc20.NewMRC20(r.SPLMRC20Addr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.TONMRC20Addr; c != "" {
		r.TONMRC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid TONMRC20Addr: %w", err)
		}
		r.TONMRC20, err = mrc20.NewMRC20(r.TONMRC20Addr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.SUIMRC20Addr; c != "" {
		r.SUIMRC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid SUIMRC20Addr: %w", err)
		}
		r.SUIMRC20, err = mrc20.NewMRC20(r.SUIMRC20Addr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.SuiTokenMRC20Addr; c != "" {
		r.SuiTokenMRC20Addr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid SuiTokenMRC20Addr: %w", err)
		}
		r.SuiTokenMRC20, err = mrc20.NewMRC20(r.SuiTokenMRC20Addr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.UniswapFactoryAddr; c != "" {
		r.UniswapV2FactoryAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid UniswapFactoryAddr: %w", err)
		}
		r.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(r.UniswapV2FactoryAddr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.UniswapRouterAddr; c != "" {
		r.UniswapV2RouterAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid UniswapRouterAddr: %w", err)
		}
		r.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(r.UniswapV2RouterAddr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.ConnectorMEVMAddr; c != "" {
		r.ConnectorMEVMAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ConnectorMEVMAddr: %w", err)
		}
		r.ConnectorMEVM, err = connectormevm.NewMuseConnectorMEVM(r.ConnectorMEVMAddr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.WMuseAddr; c != "" {
		r.WMuseAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid WMuseAddr: %w", err)
		}
		r.WMuse, err = wmuse.NewWETH9(r.WMuseAddr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.MEVMSwapAppAddr; c != "" {
		r.MEVMSwapAppAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid MEVMSwapAppAddr: %w", err)
		}
		r.MEVMSwapApp, err = mevmswap.NewMEVMSwapApp(r.MEVMSwapAppAddr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.ContextAppAddr; c != "" {
		r.ContextAppAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid ContextAppAddr: %w", err)
		}
		r.ContextApp, err = contextapp.NewContextApp(r.ContextAppAddr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.TestDappAddr; c != "" {
		r.MevmTestDAppAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid MevmTestDappAddr: %w", err)
		}
	}

	if c := conf.Contracts.EVM.TestDappAddr; c != "" {
		r.EvmTestDAppAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid EvmTestDappAddr: %w", err)
		}
	}

	if c := conf.Contracts.EVM.TestDAppV2Addr; c != "" {
		r.TestDAppV2EVMAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid TestDAppV2Addr: %w", err)
		}
		r.TestDAppV2EVM, err = testdappv2.NewTestDAppV2(r.TestDAppV2EVMAddr, r.EVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.Gateway; c != "" {
		r.GatewayMEVMAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid GatewayAddr: %w", err)
		}
		r.GatewayMEVM, err = gatewaymevm.NewGatewayMEVM(r.GatewayMEVMAddr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	if c := conf.Contracts.MEVM.TestDAppV2Addr; c != "" {
		r.TestDAppV2MEVMAddr, err = c.AsEVMAddress()
		if err != nil {
			return fmt.Errorf("invalid TestDAppV2Addr: %w", err)
		}
		r.TestDAppV2MEVM, err = testdappv2.NewTestDAppV2(r.TestDAppV2MEVMAddr, r.MEVMClient)
		if err != nil {
			return err
		}
	}

	return nil
}
