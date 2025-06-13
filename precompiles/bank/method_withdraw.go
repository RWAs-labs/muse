package bank

import (
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	precompiletypes "github.com/RWAs-labs/muse/precompiles/types"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

// withdraw is used to withdraw cosmos coins minted using the bank's deposit function.
// The caller has to have enough cosmos coin on its cosmos account balance to withdraw the requested amount.
// After all check pass the bank will burn the cosmos coins and transfer the MRC20 amount to the withdrawer.
// The cosmos coins have the denomination of "mrc20/0x12345" where 0x12345 is the MRC20 address.
// Call this function using solidity with the following signature:
// From IBank.sol: function withdraw(address mrc20, uint256 amount) external returns (bool success);
// The address to be passed to the function is the MRC20 address, like in 0x12345.
func (c *Contract) withdraw(
	ctx sdk.Context,
	evm *vm.EVM,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) (result []byte, err error) {
	// 1. Check everything is correct.
	if len(args) != 2 {
		return nil, &(precompiletypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		})
	}

	// Unpack parameters for function withdraw.
	// function withdraw(address mrc20, uint256 amount) external returns (bool success);
	mrc20Addr, amount, err := unpackWithdrawArgs(args)
	if err != nil {
		return nil, err
	}

	// Get the correct caller address.
	caller, err := precompiletypes.GetEVMCallerAddress(evm, contract)
	if err != nil {
		return nil, err
	}

	// Get the cosmos address of the caller.
	// This address should have enough cosmos coin balance as the requested amount.
	fromAddr, err := precompiletypes.GetCosmosAddress(c.bankKeeper, caller)
	if err != nil {
		return nil, err
	}

	// Safety check: token has to be a non-paused whitelisted MRC20.
	if err := c.fungibleKeeper.IsValidMRC20(ctx, mrc20Addr); err != nil {
		return nil, &precompiletypes.ErrInvalidToken{
			Got:    mrc20Addr.String(),
			Reason: err.Error(),
		}
	}

	// Caller has to have enough cosmos coin balance to withdraw the requested amount.
	coin := c.bankKeeper.GetBalance(ctx, fromAddr, precompiletypes.MRC20ToCosmosDenom(mrc20Addr))
	if !coin.IsValid() {
		return nil, &precompiletypes.ErrInsufficientBalance{
			Requested: amount.String(),
			Got:       "invalid coin",
		}
	}

	if coin.Amount.LT(math.NewIntFromBigInt(amount)) {
		return nil, &precompiletypes.ErrInsufficientBalance{
			Requested: amount.String(),
			Got:       coin.Amount.String(),
		}
	}

	coinSet, err := precompiletypes.CreateMRC20CoinSet(mrc20Addr, amount)
	if err != nil {
		return nil, err
	}

	// Check if bank address has enough MRC20 balance.
	if err := c.fungibleKeeper.CheckMRC20Balance(ctx, mrc20Addr, c.Address(), amount); err != nil {
		return nil, &precompiletypes.ErrInsufficientBalance{
			Requested: amount.String(),
			Got:       err.Error(),
		}
	}

	// 2. Effect: burn cosmos coin balance.
	if err := c.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAddr, types.ModuleName, coinSet); err != nil {
		return nil, &precompiletypes.ErrUnexpected{
			When: "SendCoinsFromAccountToModule",
			Got:  err.Error(),
		}
	}

	if err := c.bankKeeper.BurnCoins(ctx, types.ModuleName, coinSet); err != nil {
		return nil, &precompiletypes.ErrUnexpected{
			When: "BurnCoins",
			Got:  err.Error(),
		}
	}

	// 3. Interactions: send MRC20.
	if err := c.fungibleKeeper.UnlockMRC20(ctx, mrc20Addr, caller, c.Address(), amount); err != nil {
		return nil, &precompiletypes.ErrUnexpected{
			When: "UnlockMRC20InBank",
			Got:  err.Error(),
		}
	}

	if err := c.addEventLog(ctx, evm.StateDB, WithdrawEventName, eventData{caller, mrc20Addr, fromAddr.String(), coinSet.Denoms()[0], amount}); err != nil {
		return nil, &precompiletypes.ErrUnexpected{
			When: "AddWithdrawLog",
			Got:  err.Error(),
		}
	}

	return method.Outputs.Pack(true)
}

func unpackWithdrawArgs(args []interface{}) (mrc20Addr common.Address, amount *big.Int, err error) {
	mrc20Addr, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, nil, &precompiletypes.ErrInvalidAddr{
			Got: mrc20Addr.String(),
		}
	}

	amount, ok = args[1].(*big.Int)
	if !ok || amount == nil || amount.Sign() <= 0 {
		return common.Address{}, nil, &precompiletypes.ErrInvalidAmount{
			Got: amount.String(),
		}
	}

	return mrc20Addr, amount, nil
}
