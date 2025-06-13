package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/RWAs-labs/muse/x/fungible/types"
)

// EnsureGasStabilityPoolAccountCreated ensures the gas stability pool account exists
func (k Keeper) EnsureGasStabilityPoolAccountCreated(ctx sdk.Context) {
	address := types.GasStabilityPoolAddress()

	ak := k.GetAuthKeeper()
	accExists := ak.HasAccount(ctx, address)
	if !accExists {
		ak.SetAccount(ctx, ak.NewAccountWithAddress(ctx, address))
	}
}

// GetGasStabilityPoolBalance returns the balance of the gas stability pool
func (k Keeper) GetGasStabilityPoolBalance(
	ctx sdk.Context,
	chainID int64,
) (*big.Int, error) {
	// get the gas mrc20 contract from the chain
	gasMRC20, err := k.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(chainID))
	if err != nil {
		return nil, err
	}

	return k.BalanceOfMRC4(ctx, gasMRC20, types.GasStabilityPoolAddressEVM())
}

// FundGasStabilityPool mints the MRC20 into a special address called gas stability pool for the chain
func (k Keeper) FundGasStabilityPool(
	ctx sdk.Context,
	chainID int64,
	amount *big.Int,
) error {
	k.EnsureGasStabilityPoolAccountCreated(ctx)

	// get the gas mrc20 contract from the chain
	gasMRC20, err := k.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(chainID))
	if err != nil {
		return err
	}

	// call deposit MRC20 method
	return k.CallMRC20Deposit(
		ctx,
		types.ModuleAddressEVM,
		gasMRC20,
		types.GasStabilityPoolAddressEVM(),
		amount,
	)
}

// WithdrawFromGasStabilityPool burns the MRC20 from the gas stability pool
func (k Keeper) WithdrawFromGasStabilityPool(
	ctx sdk.Context,
	chainID int64,
	amount *big.Int,
) error {
	k.EnsureGasStabilityPoolAccountCreated(ctx)

	// get the gas mrc20 contract from the chain
	gasMRC20, err := k.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(chainID))
	if err != nil {
		return err
	}

	// call burn MRC20 method
	return k.CallMRC20Burn(
		ctx,
		types.GasStabilityPoolAddressEVM(),
		gasMRC20,
		amount,
		false,
	)
}
