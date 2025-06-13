package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	precompiletypes "github.com/RWAs-labs/muse/precompiles/types"
)

// balanceOf returns the balance of cosmos coins minted by the bank's deposit function,
// for a given cosmos account calculated with toAddr := sdk.AccAddress(addr.Bytes()).
// The denomination of the cosmos coin will be "mrc20/0x12345" where 0x12345 is the MRC20 address.
// Call this function using solidity with the following signature:
// From IBank.sol: function balanceOf(address mrc20, address user) external view returns (uint256 balance);
func (c *Contract) balanceOf(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) (result []byte, err error) {
	if len(args) != 2 {
		return nil, &(precompiletypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		})
	}

	// function balanceOf(address mrc20, address user) external view returns (uint256 balance);
	mrc20Addr, addr, err := unpackBalanceOfArgs(args)
	if err != nil {
		return nil, err
	}

	// Get the counterpart cosmos address.
	toAddr, err := precompiletypes.GetCosmosAddress(c.bankKeeper, addr)
	if err != nil {
		return nil, err
	}

	// Safety check: token has to be a valid whitelisted MRC20 and not be paused.
	// Do not check for t.Paused, as the balance is read only the EOA won't be able to operate.
	_, found := c.fungibleKeeper.GetForeignCoins(ctx, mrc20Addr.String())
	if !found {
		return nil, &precompiletypes.ErrInvalidToken{
			Got:    mrc20Addr.String(),
			Reason: "token is not a whitelisted MRC20",
		}
	}

	// Bank Keeper GetBalance returns the specified Cosmos coin balance for a given address.
	// Check explicitly the balance is a non-negative non-nil value.
	coin := c.bankKeeper.GetBalance(ctx, toAddr, precompiletypes.MRC20ToCosmosDenom(mrc20Addr))
	if !coin.IsValid() {
		return nil, &precompiletypes.ErrInvalidCoin{
			Got:      coin.GetDenom(),
			Negative: coin.IsNegative(),
			Nil:      coin.IsNil(),
		}
	}

	return method.Outputs.Pack(coin.Amount.BigInt())
}

func unpackBalanceOfArgs(args []interface{}) (mrc20Addr common.Address, addr common.Address, err error) {
	mrc20Addr, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, &precompiletypes.ErrInvalidAddr{
			Got: mrc20Addr.String(),
		}
	}

	addr, ok = args[1].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, &precompiletypes.ErrInvalidAddr{
			Got: addr.String(),
		}
	}

	return mrc20Addr, addr, nil
}
