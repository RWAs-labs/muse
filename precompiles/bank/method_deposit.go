package bank

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	precompiletypes "github.com/RWAs-labs/muse/precompiles/types"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

// deposit is used to deposit MRC20 into the bank contract, and receive the same amount of cosmos coins in exchange.
// The denomination of the cosmos coin will be "mrc20/MRC20Address", as an example depossiting an arbitrary MRC20 token with
// address 0x12345 will mint cosmos coins with the denomination "mrc20/0x12345".
// The caller cosmos address will be calculated from the EVM caller address. by executing toAddr := sdk.AccAddress(addr.Bytes()).
// This function can be think of a permissionless way of minting cosmos coins.
// This is how deposit works:
// - The caller has to allow the bank precompile address to spend a certain amount MRC20 token coins on its behalf. This is mandatory.
// - Then, the caller calls deposit(MRC20 address, amount), to deposit the amount and receive cosmos coins.
// - The bank will check there's enough balance, the caller is not a blocked address, and the token is a not paused MRC20.
// - Then the cosmos coins "mrc20/0x12345" will be minted and sent to the caller's cosmos address.
// Call this function using solidity with the following signature:
// - From IBank.sol: function deposit(address mrc20, uint256 amount) external returns (bool success);
func (c *Contract) deposit(
	ctx sdk.Context,
	evm *vm.EVM,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) (result []byte, err error) {
	// This function is developed using the Check - Effects - Interactions pattern:
	// 1. Check everything is correct.
	if len(args) != 2 {
		return nil, &(precompiletypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		})
	}

	// Unpack parameters for function deposit.
	// function deposit(address mrc20, uint256 amount) external returns (bool success);
	mrc20Addr, amount, err := unpackDepositArgs(args)
	if err != nil {
		return nil, err
	}

	// Get the correct caller address.
	caller, err := precompiletypes.GetEVMCallerAddress(evm, contract)
	if err != nil {
		return nil, err
	}

	// Get the cosmos address of the caller.
	toAddr, err := precompiletypes.GetCosmosAddress(c.bankKeeper, caller)
	if err != nil {
		return nil, err
	}

	// Check for enough balance.
	// function balanceOf(address account) public view virtual override returns (uint256)
	balance, err := c.fungibleKeeper.MRC20BalanceOf(ctx, mrc20Addr, caller)
	if err != nil {
		return nil, &precompiletypes.ErrUnexpected{
			When: "balanceOf",
			Got:  err.Error(),
		}
	}

	if balance.Cmp(amount) < 0 || balance.Cmp(big.NewInt(0)) <= 0 {
		return nil, &precompiletypes.ErrInvalidAmount{
			Got: balance.String(),
		}
	}

	// The process of creating a new cosmos coin is:
	// - Generate the new coin denom using MRC20 address,
	//   this way we map MRC20 addresses to cosmos denoms "mevm/0x12345".
	// - Mint coins to the fungible module.
	// - Send coins from fungible to the caller.
	coinSet, err := precompiletypes.CreateMRC20CoinSet(mrc20Addr, amount)
	if err != nil {
		return nil, err
	}

	// 2. Effect: subtract balance.
	if err := c.fungibleKeeper.LockMRC20(ctx, mrc20Addr, c.Address(), caller, c.Address(), amount); err != nil {
		return nil, &precompiletypes.ErrUnexpected{
			When: "LockMRC20InBank",
			Got:  err.Error(),
		}
	}

	// 3. Interactions: create cosmos coin and send.
	if err := c.bankKeeper.MintCoins(ctx, types.ModuleName, coinSet); err != nil {
		return nil, &precompiletypes.ErrUnexpected{
			When: "MintCoins",
			Got:  err.Error(),
		}
	}

	err = c.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddr, coinSet)
	if err != nil {
		return nil, &precompiletypes.ErrUnexpected{
			When: "SendCoinsFromModuleToAccount",
			Got:  err.Error(),
		}
	}

	if err := c.addEventLog(ctx, evm.StateDB, DepositEventName, eventData{caller, mrc20Addr, toAddr.String(), coinSet.Denoms()[0], amount}); err != nil {
		return nil, &precompiletypes.ErrUnexpected{
			When: "AddDepositLog",
			Got:  err.Error(),
		}
	}

	return method.Outputs.Pack(true)
}

func unpackDepositArgs(args []interface{}) (mrc20Addr common.Address, amount *big.Int, err error) {
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
