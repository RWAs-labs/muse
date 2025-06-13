package keeper_test

import (
	"math/big"
	"testing"

	"github.com/RWAs-labs/muse/testutil/sample"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestKeeper_MRC20SetName(t *testing.T) {
	ts := setupChain(t)

	t.Run("should update name", func(t *testing.T) {
		err := ts.fungibleKeeper.MRC20SetName(ts.ctx, ts.mrc20Address, "NewName")
		require.NoError(t, err)

		name, err := ts.fungibleKeeper.MRC20Name(ts.ctx, ts.mrc20Address)
		require.NoError(t, err)

		require.Equal(t, "NewName", name)
	})
}

func TestKeeper_MRC20SetSymbol(t *testing.T) {
	ts := setupChain(t)

	t.Run("should update symbol", func(t *testing.T) {
		err := ts.fungibleKeeper.MRC20SetSymbol(ts.ctx, ts.mrc20Address, "SYM")
		require.NoError(t, err)

		symbol, err := ts.fungibleKeeper.MRC20Symbol(ts.ctx, ts.mrc20Address)
		require.NoError(t, err)

		require.Equal(t, "SYM", symbol)
	})
}

func TestKeeper_MRC20Allowance(t *testing.T) {
	ts := setupChain(t)

	t.Run("should fail when owner is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.MRC20Allowance(
			ts.ctx,
			ts.mrc20Address,
			common.Address{},
			sample.EthAddress(),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when spender is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.MRC20Allowance(
			ts.ctx,
			ts.mrc20Address,
			sample.EthAddress(),
			common.Address{},
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when mrc20 address is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.MRC20Allowance(
			ts.ctx,
			common.Address{},
			sample.EthAddress(),
			fungibletypes.ModuleAddressEVM,
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrMRC20ZeroAddress)
	})

	t.Run("should pass with correct input", func(t *testing.T) {
		allowance, err := ts.fungibleKeeper.MRC20Allowance(
			ts.ctx,
			ts.mrc20Address,
			fungibletypes.ModuleAddressEVM,
			sample.EthAddress(),
		)
		require.NoError(t, err)
		require.Equal(t, uint64(0), allowance.Uint64())
	})
}

func TestKeeper_MRC20BalanceOf(t *testing.T) {
	ts := setupChain(t)

	t.Run("should fail when owner is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.MRC20BalanceOf(ts.ctx, ts.mrc20Address, common.Address{})
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when mrc20 address is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.MRC20BalanceOf(ts.ctx, common.Address{}, sample.EthAddress())
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrMRC20ZeroAddress)
	})

	t.Run("should pass with correct input", func(t *testing.T) {
		balance, err := ts.fungibleKeeper.MRC20BalanceOf(
			ts.ctx,
			ts.mrc20Address,
			fungibletypes.ModuleAddressEVM,
		)
		require.NoError(t, err)
		require.Equal(t, uint64(0), balance.Uint64())
	})
}

func TestKeeper_MRC20TotalSupply(t *testing.T) {
	ts := setupChain(t)

	t.Run("should fail when mrc20 address is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.MRC20TotalSupply(ts.ctx, common.Address{})
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrMRC20ZeroAddress)
	})

	t.Run("should pass with correct input", func(t *testing.T) {
		totalSupply, err := ts.fungibleKeeper.MRC20TotalSupply(ts.ctx, ts.mrc20Address)
		require.NoError(t, err)
		require.Equal(t, uint64(10000000), totalSupply.Uint64())
	})
}

func TestKeeper_MRC20Transfer(t *testing.T) {
	ts := setupChain(t)

	// Make sure sample.EthAddress() exists as an ethermint account in state.
	accAddress := sdk.AccAddress(sample.EthAddress().Bytes())
	acc := ts.fungibleKeeper.GetAuthKeeper().NewAccountWithAddress(ts.ctx, accAddress)
	ts.fungibleKeeper.GetAuthKeeper().SetAccount(ts.ctx, acc)

	t.Run("should fail when owner is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.MRC20Transfer(
			ts.ctx,
			ts.mrc20Address,
			common.Address{},
			sample.EthAddress(),
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when spender is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.MRC20Transfer(
			ts.ctx,
			ts.mrc20Address,
			sample.EthAddress(),
			common.Address{},
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when mrc20 address is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.MRC20Transfer(
			ts.ctx,
			common.Address{},
			sample.EthAddress(),
			fungibletypes.ModuleAddressEVM,
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrMRC20ZeroAddress)
	})

	t.Run("should pass with correct input", func(t *testing.T) {
		ts.fungibleKeeper.DepositMRC20(ts.ctx, ts.mrc20Address, fungibletypes.ModuleAddressEVM, big.NewInt(10))
		transferred, err := ts.fungibleKeeper.MRC20Transfer(
			ts.ctx,
			ts.mrc20Address,
			fungibletypes.ModuleAddressEVM,
			sample.EthAddress(),
			big.NewInt(10),
		)
		require.NoError(t, err)
		require.True(t, transferred)
	})
}

func TestKeeper_MRC20TransferFrom(t *testing.T) {
	// Instantiate the MRC20 ABI only one time.
	// This avoids instantiating it every time deposit or withdraw are called.
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	require.NoError(t, err)

	ts := setupChain(t)

	// Make sure sample.EthAddress() exists as an ethermint account in state.
	accAddress := sdk.AccAddress(sample.EthAddress().Bytes())
	acc := ts.fungibleKeeper.GetAuthKeeper().NewAccountWithAddress(ts.ctx, accAddress)
	ts.fungibleKeeper.GetAuthKeeper().SetAccount(ts.ctx, acc)

	t.Run("should fail when from is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.MRC20TransferFrom(
			ts.ctx,
			ts.mrc20Address,
			sample.EthAddress(),
			common.Address{},
			sample.EthAddress(),
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when to is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.MRC20TransferFrom(
			ts.ctx,
			ts.mrc20Address,
			sample.EthAddress(),
			sample.EthAddress(),
			common.Address{},
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when spender is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.MRC20TransferFrom(
			ts.ctx,
			ts.mrc20Address,
			common.Address{},
			sample.EthAddress(),
			sample.EthAddress(),
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when mrc20 address is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.MRC20TransferFrom(
			ts.ctx,
			common.Address{},
			sample.EthAddress(),
			sample.EthAddress(),
			fungibletypes.ModuleAddressEVM,
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrMRC20ZeroAddress)
	})

	t.Run("should fail without an allowance approval", func(t *testing.T) {
		// Deposit MRC20 into fungible EOA.
		ts.fungibleKeeper.DepositMRC20(ts.ctx, ts.mrc20Address, fungibletypes.ModuleAddressEVM, big.NewInt(1000))

		// Transferring the tokens with transferFrom without approval should fail.
		_, err = ts.fungibleKeeper.MRC20TransferFrom(
			ts.ctx,
			ts.mrc20Address,
			fungibletypes.ModuleAddressEVM,
			sample.EthAddress(),
			fungibletypes.ModuleAddressEVM,
			big.NewInt(10),
		)
		require.Error(t, err)
	})

	t.Run("should success with an allowance approval", func(t *testing.T) {
		// Deposit MRC20 into fungible EOA.
		ts.fungibleKeeper.DepositMRC20(ts.ctx, ts.mrc20Address, fungibletypes.ModuleAddressEVM, big.NewInt(1000))

		// Approve allowance to sample.EthAddress() to spend 10 MRC20 tokens.
		approveAllowance(t, ts, mrc20ABI, fungibletypes.ModuleAddressEVM, sample.EthAddress(), big.NewInt(10))

		// Transferring the tokens with transferFrom without approval should fail.
		_, err = ts.fungibleKeeper.MRC20TransferFrom(
			ts.ctx,
			ts.mrc20Address,
			fungibletypes.ModuleAddressEVM,
			sample.EthAddress(),
			fungibletypes.ModuleAddressEVM,
			big.NewInt(10),
		)
		require.Error(t, err)
	})
}
