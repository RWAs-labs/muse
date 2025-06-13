package keeper_test

import (
	"github.com/RWAs-labs/muse/e2e/contracts/dapp"
	"github.com/RWAs-labs/muse/e2e/contracts/dappreverter"
	"math/big"
	"testing"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/RWAs-labs/ethermint/x/evm/statedb"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestKeeper_MEVMDepositAndCallContract(t *testing.T) {
	t.Run("successfully call MUSEDepositAndCallContract on connector contract ", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		dAppContract, err := k.DeployContract(ctx, dapp.DappMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)

		museTxSender := sample.EthAddress()
		museTxReceiver := dAppContract
		inboundSenderChainID := int64(1)
		inboundAmount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		_, err = k.MUSEDepositAndCallContract(
			ctx,
			museTxSender,
			museTxReceiver,
			inboundSenderChainID,
			inboundAmount,
			data,
			cctxIndexBytes,
		)
		require.NoError(t, err)

		dappAbi, err := dapp.DappMetaData.GetAbi()
		require.NoError(t, err)
		res, err := k.CallEVM(
			ctx,
			*dappAbi,
			types.ModuleAddressEVM,
			dAppContract,
			big.NewInt(0),
			nil,
			false,
			false,
			"museTxSenderAddress",
		)
		require.NoError(t, err)
		unpacked, err := dappAbi.Unpack("museTxSenderAddress", res.Ret)
		require.NoError(t, err)
		require.NotZero(t, len(unpacked))
		valSenderAddress, ok := unpacked[0].([]byte)
		require.True(t, ok)
		require.Equal(t, museTxSender.Bytes(), valSenderAddress)
	})

	t.Run("successfully deposit coin if account is not a contract", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		museTxSender := sample.EthAddress()
		museTxReceiver := sample.EthAddress()
		inboundSenderChainID := int64(1)
		inboundAmount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		err := sdkk.EvmKeeper.SetAccount(ctx, museTxReceiver, statedb.Account{
			Nonce:    0,
			Balance:  uint256.NewInt(0),
			CodeHash: crypto.Keccak256(nil),
		})
		require.NoError(t, err)

		_, err = k.MUSEDepositAndCallContract(
			ctx,
			museTxSender,
			museTxReceiver,
			inboundSenderChainID,
			inboundAmount,
			data,
			cctxIndexBytes,
		)
		require.NoError(t, err)
		b := sdkk.BankKeeper.GetBalance(ctx, sdk.AccAddress(museTxReceiver.Bytes()), config.BaseDenom)
		require.Equal(t, inboundAmount.Int64(), b.Amount.Int64())
	})

	t.Run("automatically deposit coin  if account not found", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		museTxSender := sample.EthAddress()
		museTxReceiver := sample.EthAddress()
		inboundSenderChainID := int64(1)
		inboundAmount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		_, err := k.MUSEDepositAndCallContract(
			ctx,
			museTxSender,
			museTxReceiver,
			inboundSenderChainID,
			inboundAmount,
			data,
			cctxIndexBytes,
		)
		require.NoError(t, err)
		b := sdkk.BankKeeper.GetBalance(ctx, sdk.AccAddress(museTxReceiver.Bytes()), config.BaseDenom)
		require.Equal(t, inboundAmount.Int64(), b.Amount.Int64())
	})

	t.Run("fail MUSEDepositAndCallContract if Deposit Fails", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{UseBankMock: true})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		bankMock := keepertest.GetFungibleBankMock(t, k)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		museTxSender := sample.EthAddress()
		museTxReceiver := sample.EthAddress()
		inboundSenderChainID := int64(1)
		inboundAmount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		err := sdkk.EvmKeeper.SetAccount(ctx, museTxReceiver, statedb.Account{
			Nonce:    0,
			Balance:  uint256.NewInt(0),
			CodeHash: crypto.Keccak256(nil),
		})
		require.NoError(t, err)
		errorMint := errors.New("", 10, "error minting coins")
		bankMock.On("GetSupply", ctx, mock.Anything, mock.Anything).
			Return(sdk.NewCoin(config.BaseDenom, sdkmath.NewInt(0))).
			Once()
		bankMock.On("MintCoins", ctx, types.ModuleName, mock.Anything).Return(errorMint).Once()

		_, err = k.MUSEDepositAndCallContract(
			ctx,
			museTxSender,
			museTxReceiver,
			inboundSenderChainID,
			inboundAmount,
			data,
			cctxIndexBytes,
		)
		require.ErrorContains(t, err, errorMint.Error())
	})
}

func TestKeeper_MEVMRevertAndCallContract(t *testing.T) {
	t.Run("successfully call MUSERevertAndCallContract if receiver is a contract", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		dAppContract, err := k.DeployContract(ctx, dapp.DappMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)

		museTxSender := dAppContract
		senderChainID := big.NewInt(1)
		destinationChainID := big.NewInt(2)
		museTxReceiver := sample.EthAddress()
		amount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		_, err = k.MUSERevertAndCallContract(
			ctx,
			museTxSender,
			museTxReceiver,
			senderChainID.Int64(),
			destinationChainID.Int64(),
			amount,
			data,
			cctxIndexBytes,
		)
		require.NoError(t, err)

		dappAbi, err := dapp.DappMetaData.GetAbi()
		require.NoError(t, err)
		res, err := k.CallEVM(
			ctx,
			*dappAbi,
			types.ModuleAddressEVM,
			dAppContract,
			big.NewInt(0),
			nil,
			false,
			false,
			"museTxSenderAddress",
		)
		require.NoError(t, err)
		unpacked, err := dappAbi.Unpack("museTxSenderAddress", res.Ret)
		require.NoError(t, err)
		require.NotZero(t, len(unpacked))
		valSenderAddress, ok := unpacked[0].([]byte)
		require.True(t, ok)
		require.Equal(t, museTxSender.Bytes(), valSenderAddress)
	})

	t.Run("successfully deposit coin if account is not a contract", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		museTxSender := sample.EthAddress()
		museTxReceiver := sample.EthAddress()
		senderChainID := big.NewInt(1)
		destinationChainID := big.NewInt(2)
		amount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		err := sdkk.EvmKeeper.SetAccount(ctx, museTxSender, statedb.Account{
			Nonce:    0,
			Balance:  uint256.NewInt(0),
			CodeHash: crypto.Keccak256(nil),
		})
		require.NoError(t, err)

		_, err = k.MUSERevertAndCallContract(
			ctx,
			museTxSender,
			museTxReceiver,
			senderChainID.Int64(),
			destinationChainID.Int64(),
			amount,
			data,
			cctxIndexBytes,
		)
		require.NoError(t, err)
		b := sdkk.BankKeeper.GetBalance(ctx, sdk.AccAddress(museTxSender.Bytes()), config.BaseDenom)
		require.Equal(t, amount.Int64(), b.Amount.Int64())
	})

	t.Run("automatically deposit coin if account not found", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		museTxSender := sample.EthAddress()
		museTxReceiver := sample.EthAddress()
		senderChainID := big.NewInt(1)
		destinationChainID := big.NewInt(2)
		amount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		_, err := k.MUSERevertAndCallContract(
			ctx,
			museTxSender,
			museTxReceiver,
			senderChainID.Int64(),
			destinationChainID.Int64(),
			amount,
			data,
			cctxIndexBytes,
		)
		require.NoError(t, err)
		b := sdkk.BankKeeper.GetBalance(ctx, sdk.AccAddress(museTxSender.Bytes()), config.BaseDenom)
		require.Equal(t, amount.Int64(), b.Amount.Int64())
	})

	t.Run("fail MUSERevertAndCallContract if Deposit Fails", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{UseBankMock: true})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		bankMock := keepertest.GetFungibleBankMock(t, k)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		museTxSender := sample.EthAddress()
		museTxReceiver := sample.EthAddress()
		senderChainID := big.NewInt(1)
		destinationChainID := big.NewInt(2)
		amount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		err := sdkk.EvmKeeper.SetAccount(ctx, museTxSender, statedb.Account{
			Nonce:    0,
			Balance:  uint256.NewInt(0),
			CodeHash: crypto.Keccak256(nil),
		})
		require.NoError(t, err)
		errorMint := errors.New("", 101, "error minting coins")
		bankMock.On("GetSupply", ctx, mock.Anything, mock.Anything).
			Return(sdk.NewCoin(config.BaseDenom, sdkmath.NewInt(0))).
			Once()
		bankMock.On("MintCoins", ctx, types.ModuleName, mock.Anything).Return(errorMint).Once()

		_, err = k.MUSERevertAndCallContract(
			ctx,
			museTxSender,
			museTxReceiver,
			senderChainID.Int64(),
			destinationChainID.Int64(),
			amount,
			data,
			cctxIndexBytes,
		)
		require.ErrorIs(t, err, errorMint)
	})

	t.Run("fail MUSERevertAndCallContract if MevmOnRevert fails", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		dAppContract, err := k.DeployContract(ctx, dappreverter.DappReverterMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)

		museTxSender := dAppContract
		museTxReceiver := sample.EthAddress()
		senderChainID := big.NewInt(1)
		destinationChainID := big.NewInt(2)
		amount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		_, err = k.MUSERevertAndCallContract(
			ctx,
			museTxSender,
			museTxReceiver,
			senderChainID.Int64(),
			destinationChainID.Int64(),
			amount,
			data,
			cctxIndexBytes,
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
		require.ErrorContains(t, err, "GetSystemContract address not found")
	})
}
