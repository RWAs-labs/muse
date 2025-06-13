package keeper_test

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"cosmossdk.io/math"
	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/pkg/coin"
	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
)

func TestMsgServer_HandleEVMDeposit(t *testing.T) {
	t.Run("can process Muse deposit calling fungible method", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		sender := sample.EthAddress()
		senderChainId := int64(0)

		// expect DepositCoinMuse to be called
		fungibleMock.On("MUSEDepositAndCallContract", ctx, ethcommon.HexToAddress(sender.String()), receiver, senderChainId, amount, mock.Anything, mock.Anything).
			Return(nil, nil)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutboundParam().Receiver = receiver.String()
		cctx.GetInboundParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundParams().CoinType = coin.CoinType_Muse
		cctx.GetInboundParams().SenderChainId = senderChainId
		cctx.InboundParams.Sender = sender.String()
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.NoError(t, err)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should return error with non-reverted if deposit Muse fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		sender := sample.EthAddress()
		senderChainId := int64(0)
		amount := big.NewInt(42)
		cctx := sample.CrossChainTx(t, "foo")
		// expect DepositCoinMuse to be called
		errDeposit := errors.New("deposit failed")
		fungibleMock.On("MUSEDepositAndCallContract", ctx, ethcommon.HexToAddress(sender.String()), receiver, senderChainId, amount, mock.Anything, mock.Anything).
			Return(nil, errDeposit)

		// call HandleEVMDeposit

		cctx.InboundParams.Sender = sender.String()
		cctx.GetCurrentOutboundParam().Receiver = receiver.String()
		cctx.GetInboundParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundParams().CoinType = coin.CoinType_Muse
		cctx.GetInboundParams().SenderChainId = senderChainId
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.ErrorIs(t, err, errDeposit)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("can process ERC20 deposit calling fungible method", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID()

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinMuse to be called
		// MRC20DepositAndCallContract(ctx, from, to, msg.Amount.BigInt(), senderChain, msg.Message, contract, data, msg.FungibleTokenCoinType, msg.Asset)
		fungibleMock.On(
			"MRC20DepositAndCallContract",
			mock.Anything,
			mock.Anything,
			receiver,
			amount,
			senderChain,
			mock.Anything,
			coin.CoinType_ERC20,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, nil)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutboundParam().Receiver = receiver.String()
		cctx.GetInboundParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundParams().CoinType = coin.CoinType_ERC20
		cctx.GetInboundParams().Sender = sample.EthAddress().String()
		cctx.GetInboundParams().SenderChainId = senderChain
		cctx.RelayedMessage = ""
		cctx.GetInboundParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.NoError(t, err)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run(
		"should error on processing ERC20 deposit calling fungible method for contract call if process logs fails",
		func(t *testing.T) {
			k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
				UseFungibleMock: true,
			})

			senderChain := getValidEthChainID()

			fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
			receiver := sample.EthAddress()
			amount := big.NewInt(42)

			// expect DepositCoinMuse to be called
			// MRC20DepositAndCallContract(ctx, from, to, msg.Amount.BigInt(), senderChain, msg.Message, contract, data, msg.FungibleTokenCoinType, msg.Asset)
			fungibleMock.On(
				"MRC20DepositAndCallContract",
				mock.Anything,
				mock.Anything,
				receiver,
				amount,
				senderChain,
				mock.Anything,
				coin.CoinType_ERC20,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(&evmtypes.MsgEthereumTxResponse{
				Logs: []*evmtypes.Log{
					{
						Address:     receiver.Hex(),
						Topics:      []string{},
						Data:        []byte{},
						BlockNumber: uint64(ctx.BlockHeight()),
						TxHash:      sample.Hash().Hex(),
						TxIndex:     1,
						BlockHash:   sample.Hash().Hex(),
						Index:       1,
					},
				},
			}, true, nil)

			fungibleMock.On("GetSystemContract", mock.Anything).Return(fungibletypes.SystemContract{}, false)

			// call HandleEVMDeposit
			cctx := sample.CrossChainTx(t, "foo")
			cctx.InboundParams.TxOrigin = ""
			cctx.GetCurrentOutboundParam().Receiver = receiver.String()
			cctx.GetInboundParams().Amount = math.NewUintFromBigInt(amount)
			cctx.GetInboundParams().CoinType = coin.CoinType_ERC20
			cctx.GetInboundParams().Sender = sample.EthAddress().String()
			cctx.GetInboundParams().SenderChainId = senderChain
			cctx.RelayedMessage = ""
			cctx.GetInboundParams().Asset = ""
			reverted, err := k.HandleEVMDeposit(
				ctx,
				cctx,
			)
			require.Error(t, err)
			require.True(t, reverted)
			fungibleMock.AssertExpectations(t)
		},
	)

	t.Run(
		"can process ERC20 deposit calling fungible method for contract call if process logs doesnt fail",
		func(t *testing.T) {
			k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
				UseFungibleMock: true,
			})

			senderChain := getValidEthChainID()

			fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
			receiver := sample.EthAddress()
			amount := big.NewInt(42)

			// expect DepositCoinMuse to be called
			// MRC20DepositAndCallContract(ctx, from, to, msg.Amount.BigInt(), senderChain, msg.Message, contract, data, msg.FungibleTokenCoinType, msg.Asset)
			fungibleMock.On(
				"MRC20DepositAndCallContract",
				mock.Anything,
				mock.Anything,
				receiver,
				amount,
				senderChain,
				mock.Anything,
				coin.CoinType_ERC20,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(&evmtypes.MsgEthereumTxResponse{
				Logs: []*evmtypes.Log{
					{
						Address:     receiver.Hex(),
						Topics:      []string{},
						Data:        []byte{},
						BlockNumber: uint64(ctx.BlockHeight()),
						TxHash:      sample.Hash().Hex(),
						TxIndex:     1,
						BlockHash:   sample.Hash().Hex(),
						Index:       1,
					},
				},
			}, true, nil)

			fungibleMock.On("GetSystemContract", mock.Anything).Return(fungibletypes.SystemContract{
				ConnectorMevm: sample.EthAddress().Hex(),
			}, true)

			// call HandleEVMDeposit
			cctx := sample.CrossChainTx(t, "foo")
			cctx.InboundParams.TxOrigin = ""
			cctx.GetCurrentOutboundParam().Receiver = receiver.String()
			cctx.GetInboundParams().Amount = math.NewUintFromBigInt(amount)
			cctx.GetInboundParams().CoinType = coin.CoinType_ERC20
			cctx.GetInboundParams().Sender = sample.EthAddress().String()
			cctx.GetInboundParams().SenderChainId = senderChain
			cctx.RelayedMessage = ""
			cctx.GetInboundParams().Asset = ""
			reverted, err := k.HandleEVMDeposit(
				ctx,
				cctx,
			)
			require.NoError(t, err)
			require.False(t, reverted)
			fungibleMock.AssertExpectations(t)
		},
	)

	t.Run("should error if invalid sender", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.InboundParams.TxOrigin = ""
		cctx.GetCurrentOutboundParam().Receiver = receiver.String()
		cctx.GetInboundParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundParams().CoinType = coin.CoinType_ERC20
		cctx.GetInboundParams().Sender = "invalid"
		cctx.GetInboundParams().SenderChainId = 987
		cctx.RelayedMessage = ""
		cctx.GetInboundParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.Error(t, err)
		require.False(t, reverted)
	})

	t.Run("should return error with non-reverted if deposit ERC20 fails with tx non-failed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID()

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinMuse to be called
		// MRC20DepositAndCallContract(ctx, from, to, msg.Amount.BigInt(), senderChain, msg.Message, contract, data, msg.FungibleTokenCoinType, msg.Asset)
		errDeposit := errors.New("deposit failed")
		fungibleMock.On(
			"MRC20DepositAndCallContract",
			mock.Anything,
			mock.Anything,
			receiver,
			amount,
			senderChain,
			mock.Anything,
			coin.CoinType_ERC20,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, errDeposit)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutboundParam().Receiver = receiver.String()
		cctx.GetInboundParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundParams().CoinType = coin.CoinType_ERC20
		cctx.GetInboundParams().Sender = sample.EthAddress().String()
		cctx.GetInboundParams().SenderChainId = senderChain
		cctx.RelayedMessage = ""
		cctx.GetInboundParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.ErrorIs(t, err, errDeposit)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should return error with reverted if deposit ERC20 fails with tx failed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID()

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinMuse to be called
		// MRC20DepositAndCallContract(ctx, from, to, msg.Amount.BigInt(), senderChain, msg.Message, contract, data, msg.FungibleTokenCoinType, msg.Asset)
		errDeposit := errors.New("deposit failed")
		fungibleMock.On(
			"MRC20DepositAndCallContract",
			mock.Anything,
			mock.Anything,
			receiver,
			amount,
			senderChain,
			mock.Anything,
			coin.CoinType_ERC20,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{VmError: "reverted"}, false, errDeposit)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutboundParam().Receiver = receiver.String()
		cctx.GetInboundParams().Amount = math.NewUintFromBigInt(amount)
		cctx.InboundParams.CoinType = coin.CoinType_ERC20
		cctx.GetInboundParams().Sender = sample.EthAddress().String()
		cctx.GetInboundParams().SenderChainId = senderChain
		cctx.RelayedMessage = ""
		cctx.GetInboundParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.ErrorIs(t, err, errDeposit)
		require.True(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should return error with reverted if deposit ERC20 fails with liquidity cap reached", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID()

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinMuse to be called
		// MRC20DepositAndCallContract(ctx, from, to, msg.Amount.BigInt(), senderChain, msg.Message, contract, data, msg.FungibleTokenCoinType, msg.Asset)
		fungibleMock.On(
			"MRC20DepositAndCallContract",
			mock.Anything,
			mock.Anything,
			receiver,
			amount,
			senderChain,
			mock.Anything,
			coin.CoinType_ERC20,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, fungibletypes.ErrForeignCoinCapReached)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutboundParam().Receiver = receiver.String()
		cctx.GetInboundParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundParams().CoinType = coin.CoinType_ERC20
		cctx.GetInboundParams().Sender = sample.EthAddress().String()
		cctx.GetInboundParams().SenderChainId = senderChain
		cctx.RelayedMessage = ""
		cctx.GetInboundParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.ErrorIs(t, err, fungibletypes.ErrForeignCoinCapReached)
		require.True(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run("should return error with reverted if deposit ERC20 fails with mrc20 paused", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID()

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinMuse to be called
		// MRC20DepositAndCallContract(ctx, from, to, msg.Amount.BigInt(), senderChain, msg.Message, contract, data, msg.FungibleTokenCoinType, msg.Asset)
		fungibleMock.On(
			"MRC20DepositAndCallContract",
			mock.Anything,
			mock.Anything,
			receiver,
			amount,
			senderChain,
			mock.Anything,
			coin.CoinType_ERC20,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, fungibletypes.ErrPausedMRC20)

		// call HandleEVMDeposit
		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutboundParam().Receiver = receiver.String()
		cctx.GetInboundParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundParams().CoinType = coin.CoinType_ERC20
		cctx.GetInboundParams().Sender = sample.EthAddress().String()
		cctx.GetInboundParams().SenderChainId = senderChain
		cctx.RelayedMessage = ""
		cctx.GetInboundParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.ErrorIs(t, err, fungibletypes.ErrPausedMRC20)
		require.True(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	t.Run(
		"should return error with reverted if deposit ERC20 fails with calling a non-contract address",
		func(t *testing.T) {
			k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
				UseFungibleMock: true,
			})

			senderChain := getValidEthChainID()

			fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
			receiver := sample.EthAddress()
			amount := big.NewInt(42)

			fungibleMock.On(
				"MRC20DepositAndCallContract",
				mock.Anything,
				mock.Anything,
				receiver,
				amount,
				senderChain,
				mock.Anything,
				coin.CoinType_ERC20,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(&evmtypes.MsgEthereumTxResponse{}, false, fungibletypes.ErrCallNonContract)

			// call HandleEVMDeposit
			cctx := sample.CrossChainTx(t, "foo")
			cctx.GetCurrentOutboundParam().Receiver = receiver.String()
			cctx.GetInboundParams().Amount = math.NewUintFromBigInt(amount)
			cctx.GetInboundParams().CoinType = coin.CoinType_ERC20
			cctx.GetInboundParams().Sender = sample.EthAddress().String()
			cctx.GetInboundParams().SenderChainId = senderChain
			cctx.RelayedMessage = ""
			cctx.GetInboundParams().Asset = ""
			reverted, err := k.HandleEVMDeposit(
				ctx,
				cctx,
			)
			require.ErrorIs(t, err, fungibletypes.ErrCallNonContract)
			require.True(t, reverted)
			fungibleMock.AssertExpectations(t)
		},
	)

	t.Run("should fail if can't parse address and data", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		senderChain := getValidEthChainID()

		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutboundParam().Receiver = sample.EthAddress().String()
		cctx.GetInboundParams().Amount = math.NewUint(42)
		cctx.GetInboundParams().CoinType = coin.CoinType_Gas
		cctx.GetInboundParams().Sender = sample.EthAddress().String()
		cctx.GetInboundParams().SenderChainId = senderChain
		cctx.RelayedMessage = "not_hex"
		cctx.GetInboundParams().Asset = ""
		_, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.ErrorIs(t, err, types.ErrUnableToParseAddress)
	})

	t.Run("should deposit into address if address is parsed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID()

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		data, err := hex.DecodeString("DEADBEEF")
		require.NoError(t, err)
		cctx := sample.CrossChainTx(t, "foo")
		b, err := cctx.Marshal()
		require.NoError(t, err)
		ctx = ctx.WithTxBytes(b)
		fungibleMock.On(
			"MRC20DepositAndCallContract",
			mock.Anything,
			mock.Anything,
			receiver,
			amount,
			senderChain,
			data,
			coin.CoinType_ERC20,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, nil)

		cctx.GetCurrentOutboundParam().Receiver = sample.EthAddress().String()
		cctx.GetInboundParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundParams().CoinType = coin.CoinType_ERC20
		cctx.GetInboundParams().Sender = sample.EthAddress().String()
		cctx.GetInboundParams().SenderChainId = senderChain
		cctx.RelayedMessage = receiver.Hex()[2:] + "DEADBEEF"
		cctx.GetInboundParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.NoError(t, err)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
		require.Equal(t, uint64(ctx.BlockHeight()), cctx.GetCurrentOutboundParam().ObservedExternalHeight)
	})

	t.Run("should deposit into receiver with specified data if no address parsed with data", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		senderChain := getValidEthChainID()

		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		data, err := hex.DecodeString("DEADBEEF")
		require.NoError(t, err)
		fungibleMock.On(
			"MRC20DepositAndCallContract",
			mock.Anything,
			mock.Anything,
			receiver,
			amount,
			senderChain,
			data,
			coin.CoinType_ERC20,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&evmtypes.MsgEthereumTxResponse{}, false, nil)

		cctx := sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutboundParam().Receiver = receiver.String()
		cctx.GetInboundParams().Amount = math.NewUintFromBigInt(amount)
		cctx.GetInboundParams().CoinType = coin.CoinType_ERC20
		cctx.GetInboundParams().Sender = sample.EthAddress().String()
		cctx.GetInboundParams().SenderChainId = senderChain
		cctx.RelayedMessage = "DEADBEEF"
		cctx.GetInboundParams().Asset = ""
		reverted, err := k.HandleEVMDeposit(
			ctx,
			cctx,
		)
		require.NoError(t, err)
		require.False(t, reverted)
		fungibleMock.AssertExpectations(t)
	})

	// TODO: add test cases for testing logs process
	// https://github.com/RWAs-labs/muse/issues/1207
}
