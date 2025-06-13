package keeper_test

import (
	"errors"
	"fmt"
	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	"github.com/RWAs-labs/muse/pkg/chains"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/pkg/coin"
	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
)

func TestKeeper_ProcessAbort(t *testing.T) {
	t.Run("set abort status without processing abort if not v2", func(t *testing.T) {
		// arrange
		// define fungible mock to assess processAbort is not called
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		cctx := sample.CrossChainTx(t, "index")
		cctx.ProtocolContractVersion = types.ProtocolContractVersion_V1

		// act
		k.ProcessAbort(ctx, cctx, types.StatusMessages{})

		// assert
		require.EqualValues(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.False(t, cctx.CctxStatus.IsAbortRefunded)
	})

	t.Run("set abort status without processing abort if no abort address", func(t *testing.T) {
		// arrange
		// define fungible mock to assess processAbort is not called
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		cctx := sample.CrossChainTx(t, "index")
		cctx.ProtocolContractVersion = types.ProtocolContractVersion_V2
		cctx.RevertOptions.AbortAddress = ""

		// act
		k.ProcessAbort(ctx, cctx, types.StatusMessages{})

		// assert
		require.EqualValues(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.False(t, cctx.CctxStatus.IsAbortRefunded)
	})

	t.Run("fail abort with abort error message if connected chain ID can't be retrieved", func(t *testing.T) {
		// arrange
		// define fungible mock to assess processAbort is not called
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		cctx := sample.CrossChainTx(t, "index")
		cctx.ProtocolContractVersion = types.ProtocolContractVersion_V2
		cctx.RevertOptions.AbortAddress = sample.EthAddress().Hex()

		// set inbound to nil to make GetConnectedChainID fail
		cctx.InboundParams = nil

		// act
		k.ProcessAbort(ctx, cctx, types.StatusMessages{})

		// assert
		require.EqualValues(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.False(t, cctx.CctxStatus.IsAbortRefunded)
		require.Contains(t, cctx.CctxStatus.ErrorMessageAbort, "failed to get connected chain ID")
	})

	t.Run("process abort by calling process abort of fungible", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On(
			"ProcessAbort",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, nil)

		cctx := sample.CrossChainTx(t, "index")
		cctx.ProtocolContractVersion = types.ProtocolContractVersion_V2
		cctx.RevertOptions.AbortAddress = sample.EthAddress().Hex()

		// act
		k.ProcessAbort(ctx, cctx, types.StatusMessages{})

		// assert
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Empty(t, cctx.CctxStatus.ErrorMessageAbort)
		require.True(t, cctx.CctxStatus.IsAbortRefunded)
	})

	t.Run(
		"fail abort with abort error message if process abort fails, status not set to refunded if error other than onAbort failure",
		func(t *testing.T) {
			// Arrange
			k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
				UseFungibleMock: true,
			})

			fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
			fungibleMock.On(
				"ProcessAbort",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(&evmtypes.MsgEthereumTxResponse{}, errors.New("process abort failed"))

			cctx := sample.CrossChainTx(t, "index")
			cctx.ProtocolContractVersion = types.ProtocolContractVersion_V2
			cctx.RevertOptions.AbortAddress = sample.EthAddress().Hex()

			// act
			k.ProcessAbort(ctx, cctx, types.StatusMessages{})

			// assert
			require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
			require.Contains(t, cctx.CctxStatus.ErrorMessageAbort, "process abort failed")
			require.False(t, cctx.CctxStatus.IsAbortRefunded)
		},
	)

	t.Run(
		"fail abort with abort error message if process abort fails, status set to refunded if error is onAbort failure",
		func(t *testing.T) {
			// Arrange
			k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
				UseFungibleMock: true,
			})

			fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
			fungibleMock.On(
				"ProcessAbort",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(&evmtypes.MsgEthereumTxResponse{}, errors.Join(errors.New("process abort failed"), fungibletypes.ErrOnAbortFailed))

			cctx := sample.CrossChainTx(t, "index")
			cctx.ProtocolContractVersion = types.ProtocolContractVersion_V2
			cctx.RevertOptions.AbortAddress = sample.EthAddress().Hex()

			// act
			k.ProcessAbort(ctx, cctx, types.StatusMessages{})

			// assert
			require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
			require.Contains(t, cctx.CctxStatus.ErrorMessageAbort, "process abort failed")
			require.True(t, cctx.CctxStatus.IsAbortRefunded)
		},
	)
}

func TestKeeper_RefundAmountOnMuseChainGas(t *testing.T) {
	t.Run("should refund amount mrc20 gas on muse chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID()
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		err := k.LegacyRefundAbortedAmountOnMuseChainGas(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.NewUint(20),
			},
			OutboundParams: []*types.OutboundParams{{
				Amount: math.NewUint(42),
			}},
		},
			sender,
		)
		require.NoError(t, err)
		balance, err := zk.FungibleKeeper.BalanceOfMRC4(ctx, mrc20, sender)
		require.NoError(t, err)
		require.Equal(t, uint64(42), balance.Uint64())
	})

	t.Run("should error if mrc20 address empty", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID()

		fungibleMock.On("GetGasCoinForForeignCoin", mock.Anything, mock.Anything).Return(fungibletypes.ForeignCoins{
			Mrc20ContractAddress: "0x",
		}, true)

		err := k.LegacyRefundAbortedAmountOnMuseChainGas(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.NewUint(20),
			},
			OutboundParams: []*types.OutboundParams{{
				Amount: math.NewUint(42),
			}},
		},
			sender,
		)
		require.Error(t, err)
	})

	t.Run("should error if deposit mrc20 fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID()

		fungibleMock.On("GetGasCoinForForeignCoin", mock.Anything, mock.Anything).Return(fungibletypes.ForeignCoins{
			Mrc20ContractAddress: sample.EthAddress().Hex(),
		}, true)

		fungibleMock.On("DepositMRC20", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New(""))

		err := k.LegacyRefundAbortedAmountOnMuseChainGas(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.NewUint(20),
			},
			OutboundParams: []*types.OutboundParams{{
				Amount: math.NewUint(42),
			}},
		},
			sender,
		)
		require.Error(t, err)
	})

	t.Run("should refund inbound amount mrc20 gas on muse chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID()
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		err := k.LegacyRefundAbortedAmountOnMuseChainGas(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.NewUint(20),
			},
		},
			sender,
		)
		require.NoError(t, err)
		balance, err := zk.FungibleKeeper.BalanceOfMRC4(ctx, mrc20, sender)
		require.NoError(t, err)
		require.Equal(t, uint64(20), balance.Uint64())
	})
	t.Run("failed refund mrc20 gas on muse chain if gas coin not found", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID()
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		err := k.LegacyRefundAbortedAmountOnMuseChainGas(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.NewUint(20),
			},
			OutboundParams: []*types.OutboundParams{{
				Amount: math.NewUint(42),
			}},
		},

			sender,
		)
		require.ErrorContains(t, err, types.ErrForeignCoinNotFound.Error())
	})
	t.Run("failed refund amount mrc20 gas on muse chain if amount is 0", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID()
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		_ = setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		err := k.LegacyRefundAbortedAmountOnMuseChainGas(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.ZeroUint(),
			},
			OutboundParams: []*types.OutboundParams{{
				Amount: math.ZeroUint(),
			}},
		},
			sender,
		)
		require.ErrorContains(t, err, "no amount to refund")
	})

}

func TestKeeper_RefundAmountOnMuseChainMuse(t *testing.T) {
	t.Run("should refund amount on muse chain", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID()

		err := k.LegacyRefundAbortedAmountOnMuseChainMuse(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.NewUint(20),
			},
			OutboundParams: []*types.OutboundParams{{
				Amount: math.NewUint(42),
			}},
		},
			sender,
		)
		require.NoError(t, err)
		coin := sdkk.BankKeeper.GetBalance(ctx, sender.Bytes(), config.BaseDenom)
		fmt.Println(coin.Amount.String())
		require.Equal(t, "42", coin.Amount.String())
	})

	t.Run("should error if non evm chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidBtcChainID()

		err := k.LegacyRefundAbortedAmountOnMuseChainMuse(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.NewUint(20),
			},
			OutboundParams: []*types.OutboundParams{{
				Amount: math.NewUint(42),
			}},
		},
			sender,
		)
		require.Error(t, err)
	})

	t.Run("should error if deposit coin muse fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("DepositCoinMuse", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("err"))
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID()

		err := k.LegacyRefundAbortedAmountOnMuseChainMuse(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.NewUint(20),
			},
			OutboundParams: []*types.OutboundParams{{
				Amount: math.NewUint(42),
			}},
		},
			sender,
		)
		require.Error(t, err)
	})

	t.Run("should refund inbound amount on muse chain if outbound is not present", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID()

		err := k.LegacyRefundAbortedAmountOnMuseChainMuse(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.NewUint(20),
			},
		},
			sender,
		)
		require.NoError(t, err)
		coin := sdkk.BankKeeper.GetBalance(ctx, sdk.AccAddress(sender.Bytes()), config.BaseDenom)
		require.Equal(t, "20", coin.Amount.String())
	})
	t.Run("failed refund amount on muse chain amount is 0", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		sender := sample.EthAddress()
		chainID := getValidEthChainID()

		err := k.LegacyRefundAbortedAmountOnMuseChainMuse(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_Gas,
				SenderChainId: chainID,
				Sender:        sender.String(),
				TxOrigin:      sender.String(),
				Amount:        math.ZeroUint(),
			},
			OutboundParams: []*types.OutboundParams{{
				Amount: math.ZeroUint(),
			}},
		},
			sender,
		)
		require.ErrorContains(t, err, "no amount to refund")
	})
}

func TestKeeper_RefundAmountOnMuseChainERC20(t *testing.T) {
	t.Run("should refund amount on muse chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		asset := sample.EthAddress().String()
		sender := sample.EthAddress()
		chainID := getValidEthChainID()

		// deploy mrc20
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		mrc20Addr := deployMRC20(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.EvmKeeper,
			chainID,
			"bar",
			asset,
			"bar",
		)

		err := k.LegacyRefundAbortedAmountOnMuseChainERC20(ctx, types.CrossChainTx{

			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_ERC20,
				SenderChainId: chainID,
				Sender:        sender.String(),
				Asset:         asset,
				Amount:        math.NewUint(42),
			},
			OutboundParams: []*types.OutboundParams{{
				Amount: math.NewUint(42),
			}},
		},
			sender,
		)
		require.NoError(t, err)

		// check amount deposited in balance
		balance, err := zk.FungibleKeeper.BalanceOfMRC4(ctx, mrc20Addr, sender)
		require.NoError(t, err)
		require.Equal(t, uint64(42), balance.Uint64())

		// can refund again
		err = k.LegacyRefundAbortedAmountOnMuseChainERC20(ctx, types.CrossChainTx{

			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_ERC20,
				SenderChainId: chainID,
				Sender:        sender.String(),
				Asset:         asset,
				Amount:        math.NewUint(42),
			}},
			sender,
		)
		require.NoError(t, err)
		balance, err = zk.FungibleKeeper.BalanceOfMRC4(ctx, mrc20Addr, sender)
		require.NoError(t, err)
		require.Equal(t, uint64(84), balance.Uint64())
	})

	t.Run("should refund amount on muse chain for outgoing cctx", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		asset := sample.EthAddress().String()
		sender := sample.EthAddress()
		chainID := getValidEthChainID()

		// deploy mrc20
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		mrc20Addr := deployMRC20(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.EvmKeeper,
			chainID,
			"bar",
			asset,
			"bar",
		)

		err := k.LegacyRefundAbortedAmountOnMuseChainERC20(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_ERC20,
				SenderChainId: chains.MuseChainPrivnet.ChainId,
				Sender:        sender.String(),
				Asset:         asset,
				Amount:        math.NewUint(42),
			},
			OutboundParams: []*types.OutboundParams{{
				ReceiverChainId: chainID,
				Amount:          math.NewUint(42),
			}},
		},
			sender,
		)
		require.NoError(t, err)

		// check amount deposited in balance
		balance, err := zk.FungibleKeeper.BalanceOfMRC4(ctx, mrc20Addr, sender)
		require.NoError(t, err)
		require.Equal(t, uint64(42), balance.Uint64())
	})

	t.Run("should error if mrc20 address empty", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, mock.Anything, mock.Anything).
			Return(fungibletypes.ForeignCoins{
				Mrc20ContractAddress: "0x",
			}, true)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		asset := sample.EthAddress().String()
		sender := sample.EthAddress()
		chainID := getValidEthChainID()

		err := k.LegacyRefundAbortedAmountOnMuseChainERC20(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_ERC20,
				SenderChainId: chainID,
				Sender:        sender.String(),
				Asset:         asset,
				Amount:        math.NewUint(42),
			},
			OutboundParams: []*types.OutboundParams{{
				Amount: math.NewUint(42),
			}},
		},
			sender,
		)
		require.Error(t, err)
	})

	t.Run("should error if deposit mrc20 fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		asset := sample.EthAddress().String()
		sender := sample.EthAddress()
		chainID := getValidEthChainID()

		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, mock.Anything, mock.Anything).
			Return(fungibletypes.ForeignCoins{
				Mrc20ContractAddress: sample.EthAddress().Hex(),
			}, true)

		fungibleMock.On("DepositMRC20", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New(""))

		err := k.LegacyRefundAbortedAmountOnMuseChainERC20(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_ERC20,
				SenderChainId: chainID,
				Sender:        sender.String(),
				Asset:         asset,
				Amount:        math.NewUint(42),
			},
			OutboundParams: []*types.OutboundParams{{
				Amount: math.NewUint(42),
			}},
		},
			sender,
		)
		require.Error(t, err)
	})

	t.Run("should fail with invalid cctx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		err := k.LegacyRefundAbortedAmountOnMuseChainERC20(ctx, types.CrossChainTx{

			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_Muse,
				Amount:   math.NewUint(42),
			}},
			sample.EthAddress(),
		)
		require.ErrorContains(t, err, "unsupported coin type")

		err = k.LegacyRefundAbortedAmountOnMuseChainERC20(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_Gas,
			}},
			sample.EthAddress(),
		)
		require.ErrorContains(t, err, "unsupported coin type")

		err = k.LegacyRefundAbortedAmountOnMuseChainERC20(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_ERC20,
				SenderChainId: 999999,
				Amount:        math.NewUint(42),
			}},
			sample.EthAddress(),
		)
		require.ErrorContains(t, err, "only EVM chains are supported")

		err = k.LegacyRefundAbortedAmountOnMuseChainERC20(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_ERC20,
				SenderChainId: getValidEthChainID(),
				Sender:        sample.EthAddress().String(),
				Amount:        math.Uint{},
			}},
			sample.EthAddress(),
		)
		require.ErrorContains(t, err, "no amount to refund")

		err = k.LegacyRefundAbortedAmountOnMuseChainERC20(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_ERC20,
				SenderChainId: getValidEthChainID(),
				Sender:        sample.EthAddress().String(),
				Amount:        math.ZeroUint(),
			}},
			sample.EthAddress(),
		)
		require.ErrorContains(t, err, "no amount to refund")

		// the foreign coin has not been set
		err = k.LegacyRefundAbortedAmountOnMuseChainERC20(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_ERC20,
				SenderChainId: getValidEthChainID(),
				Sender:        sample.EthAddress().String(),
				Asset:         sample.EthAddress().String(),
				Amount:        math.NewUint(42),
			}},
			sample.EthAddress(),
		)
		require.ErrorContains(t, err, "mrc not found")
	})
}

func TestKeeper_RefundAbortedAmountOnMuseChain_FailsForUnsupportedCoinType(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)

	cctx := sample.CrossChainTx(t, "index")
	cctx.InboundParams.CoinType = coin.CoinType_Cmd
	err := k.LegacyRefundAbortedAmountOnMuseChain(ctx, *cctx, common.Address{})
	require.ErrorContains(t, err, "unsupported coin type for refund on MuseChain")
}
