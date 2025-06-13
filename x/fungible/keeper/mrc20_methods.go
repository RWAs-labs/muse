package keeper

import (
	"fmt"
	"math/big"

	"cosmossdk.io/errors"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/RWAs-labs/muse/pkg/crypto"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
)

const (
	setName      = "setName"
	setSymbol    = "setSymbol"
	name         = "name"
	symbol       = "symbol"
	allowance    = "allowance"
	balanceOf    = "balanceOf"
	totalSupply  = "totalSupply"
	transfer     = "transfer"
	transferFrom = "transferFrom"
)

// MRC20SetName updates the name of a MRC20 token
func (k Keeper) MRC20SetName(
	ctx sdk.Context,
	mrc20Address common.Address,
	newName string,
) error {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	if err != nil {
		return err
	}

	// function setName(string memory name)
	res, err := k.CallEVM(
		ctx,
		*mrc20ABI,
		fungibletypes.ModuleAddressEVM,
		mrc20Address,
		big.NewInt(0),
		DefaultGasLimit,
		true,
		true,
		setName,
		newName,
	)
	if err != nil {
		return errors.Wrap(err, "EVM error calling MRC20 setName function")
	}
	if res.VmError != "" {
		return fmt.Errorf("EVM execution error calling allowance: %s", res.VmError)
	}

	return nil
}

// MRC20SetSymbol updates the symbol of a MRC20 token
func (k Keeper) MRC20SetSymbol(
	ctx sdk.Context,
	mrc20Address common.Address,
	newSymbol string,
) error {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	if err != nil {
		return err
	}

	// function setSymbol(string memory symbol)
	res, err := k.CallEVM(
		ctx,
		*mrc20ABI,
		fungibletypes.ModuleAddressEVM,
		mrc20Address,
		big.NewInt(0),
		DefaultGasLimit,
		true,
		true,
		setSymbol,
		newSymbol,
	)
	if err != nil {
		return errors.Wrap(err, "EVM error calling MRC20 setSymbol function")
	}
	if res.VmError != "" {
		return fmt.Errorf("EVM execution error calling allowance: %s", res.VmError)
	}

	return nil
}

// MRC20Name returns the name of a MRC20 token
func (k Keeper) MRC20Name(
	ctx sdk.Context,
	mrc20Address common.Address,
) (string, error) {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	if err != nil {
		return "", err
	}

	// function name() public view virtual override returns (string memory)
	res, err := k.CallEVM(
		ctx,
		*mrc20ABI,
		fungibletypes.ModuleAddressEVM,
		mrc20Address,
		big.NewInt(0),
		nil,
		false,
		true,
		name,
	)
	if err != nil {
		return "", errors.Wrap(err, "EVM error calling MRC20 name function")
	}

	if res.VmError != "" {
		return "", fmt.Errorf("EVM execution error calling name: %s", res.VmError)
	}

	ret, err := mrc20ABI.Methods[name].Outputs.Unpack(res.Ret)
	if err != nil {
		return "", errors.Wrap(err, "failed to unpack MRC20 name return value")
	}

	if len(ret) == 0 {
		return "", fmt.Errorf("no data returned from 'name' method")
	}

	name, ok := ret[0].(string)
	if !ok {
		return "", fmt.Errorf("MRC20 name returned an unexpected type")
	}

	return name, nil
}

// MRC20Symbol returns the symbol of a MRC20 token
func (k Keeper) MRC20Symbol(
	ctx sdk.Context,
	mrc20Address common.Address,
) (string, error) {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	if err != nil {
		return "", err
	}

	// function symbol() public view virtual override returns (string memory)
	res, err := k.CallEVM(
		ctx,
		*mrc20ABI,
		fungibletypes.ModuleAddressEVM,
		mrc20Address,
		big.NewInt(0),
		nil,
		false,
		true,
		symbol,
	)
	if err != nil {
		return "", errors.Wrap(err, "EVM error calling MRC20 symbol function")
	}

	if res.VmError != "" {
		return "", fmt.Errorf("EVM execution error calling symbol: %s", res.VmError)
	}

	ret, err := mrc20ABI.Methods[symbol].Outputs.Unpack(res.Ret)
	if err != nil {
		return "", errors.Wrap(err, "failed to unpack MRC20 symbol return value")
	}

	if len(ret) == 0 {
		return "", fmt.Errorf("no data returned from 'symbol' method")
	}

	symbol, ok := ret[0].(string)
	if !ok {
		return "", fmt.Errorf("MRC20 symbol returned an unexpected type")
	}

	return symbol, nil
}

// MRC20Allowance returns the MRC20 allowance for a given spender.
// The allowance has to be previously approved by the MRC20 tokens owner.
func (k Keeper) MRC20Allowance(
	ctx sdk.Context,
	mrc20Address, owner, spender common.Address,
) (*big.Int, error) {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	if crypto.IsEmptyAddress(owner) || crypto.IsEmptyAddress(spender) {
		return nil, fungibletypes.ErrZeroAddress
	}

	if err := k.IsValidMRC20(ctx, mrc20Address); err != nil {
		return nil, err
	}

	// function allowance(address owner, address spender)
	args := []interface{}{owner, spender}
	res, err := k.CallEVM(
		ctx,
		*mrc20ABI,
		fungibletypes.ModuleAddressEVM,
		mrc20Address,
		big.NewInt(0),
		nil,
		false,
		true,
		allowance,
		args...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "EVM error calling MRC20 allowance function")
	}

	if res.VmError != "" {
		return nil, fmt.Errorf("EVM execution error calling allowance: %s", res.VmError)
	}

	ret, err := mrc20ABI.Methods[allowance].Outputs.Unpack(res.Ret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unpack MRC20 allowance return value")
	}

	if len(ret) == 0 {
		return nil, fmt.Errorf("no data returned from 'allowance' method")
	}

	allowanceValue, ok := ret[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("MRC20 allowance returned an unexpected type")
	}

	return allowanceValue, nil
}

// MRC20BalanceOf checks the MRC20 balance of a given EOA.
func (k Keeper) MRC20BalanceOf(
	ctx sdk.Context,
	mrc20Address, owner common.Address,
) (*big.Int, error) {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	if crypto.IsEmptyAddress(owner) {
		return nil, fungibletypes.ErrZeroAddress
	}

	if err := k.IsValidMRC20(ctx, mrc20Address); err != nil {
		return nil, err
	}

	// function balanceOf(address account)
	res, err := k.CallEVM(
		ctx,
		*mrc20ABI,
		fungibletypes.ModuleAddressEVM,
		mrc20Address,
		big.NewInt(0),
		nil,
		false,
		true,
		balanceOf,
		owner,
	)
	if err != nil {
		return nil, errors.Wrap(err, "EVM error calling MRC20 balanceOf function")
	}

	if res.VmError != "" {
		return nil, fmt.Errorf("EVM execution error calling balanceOf: %s", res.VmError)
	}

	ret, err := mrc20ABI.Methods[balanceOf].Outputs.Unpack(res.Ret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unpack MRC20 balanceOf return value")
	}

	if len(ret) == 0 {
		return nil, fmt.Errorf("no data returned from 'balanceOf' method")
	}

	balance, ok := ret[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("MRC20 balanceOf returned an unexpected type")
	}

	return balance, nil
}

// MRC20TotalSupply returns the total supply of a MRC20 token.
func (k Keeper) MRC20TotalSupply(
	ctx sdk.Context,
	mrc20Address common.Address,
) (*big.Int, error) {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	if err := k.IsValidMRC20(ctx, mrc20Address); err != nil {
		return nil, err
	}

	// function totalSupply() public view virtual override returns (uint256)
	res, err := k.CallEVM(
		ctx,
		*mrc20ABI,
		fungibletypes.ModuleAddressEVM,
		mrc20Address,
		big.NewInt(0),
		nil,
		false,
		true,
		totalSupply,
	)
	if err != nil {
		return nil, errors.Wrap(err, "EVM error calling MRC20 totalSupply function")
	}

	if res.VmError != "" {
		return nil, fmt.Errorf("EVM execution error calling totalSupply: %s", res.VmError)
	}

	ret, err := mrc20ABI.Methods[totalSupply].Outputs.Unpack(res.Ret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unpack MRC20 totalSupply return value")
	}

	if len(ret) == 0 {
		return nil, fmt.Errorf("no data returned from 'totalSupply' method")
	}

	totalSupply, ok := ret[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("MRC20 totalSupply returned an unexpected type")
	}

	return totalSupply, nil
}

// MRC20Transfer transfers MRC20 tokens from the sender to the recipient.
func (k Keeper) MRC20Transfer(
	ctx sdk.Context,
	mrc20Address, from, to common.Address,
	amount *big.Int,
) (bool, error) {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	if err != nil {
		return false, err
	}

	if crypto.IsEmptyAddress(from) || crypto.IsEmptyAddress(to) {
		return false, fungibletypes.ErrZeroAddress
	}

	if err := k.IsValidMRC20(ctx, mrc20Address); err != nil {
		return false, err
	}

	// function transfer(address recipient, uint256 amount)
	args := []interface{}{to, amount}
	res, err := k.CallEVM(
		ctx,
		*mrc20ABI,
		from,
		mrc20Address,
		big.NewInt(0),
		DefaultGasLimit,
		true,
		true,
		transfer,
		args...,
	)
	if err != nil {
		return false, errors.Wrap(err, "EVM error calling MRC20 transfer function")
	}

	if res.VmError != "" {
		return false, fmt.Errorf("EVM execution error in transfer: %s", res.VmError)
	}

	ret, err := mrc20ABI.Methods[transfer].Outputs.Unpack(res.Ret)
	if err != nil {
		return false, errors.Wrap(err, "failed to unpack MRC20 transfer return value")
	}

	if len(ret) == 0 {
		return false, fmt.Errorf("no data returned from 'transfer' method")
	}

	transferred, ok := ret[0].(bool)
	if !ok {
		return false, fmt.Errorf("transfer returned an unexpected value")
	}

	return transferred, nil
}

// MRC20TransferFrom transfers MRC20 tokens "from" to the EOA "to".
// The transaction is started by the spender.
// Requisite: the original EOA must have approved the spender to spend the tokens.
func (k Keeper) MRC20TransferFrom(
	ctx sdk.Context,
	mrc20Address, spender, from, to common.Address,
	amount *big.Int,
) (bool, error) {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	if err != nil {
		return false, err
	}

	if crypto.IsEmptyAddress(from) || crypto.IsEmptyAddress(to) || crypto.IsEmptyAddress(spender) {
		return false, fungibletypes.ErrZeroAddress
	}

	if err := k.IsValidMRC20(ctx, mrc20Address); err != nil {
		return false, err
	}

	// function transferFrom(address sender, address recipient, uint256 amount)
	args := []interface{}{from, to, amount}
	res, err := k.CallEVM(
		ctx,
		*mrc20ABI,
		spender,
		mrc20Address,
		big.NewInt(0),
		DefaultGasLimit,
		true,
		true,
		transferFrom,
		args...,
	)
	if err != nil {
		return false, errors.Wrap(err, "EVM error calling MRC20 transferFrom function")
	}

	if res.VmError != "" {
		return false, fmt.Errorf("EVM execution error in transferFrom: %s", res.VmError)
	}

	ret, err := mrc20ABI.Methods[transferFrom].Outputs.Unpack(res.Ret)
	if err != nil {
		return false, errors.Wrap(err, "failed to unpack MRC20 transferFrom return value")
	}

	if len(ret) == 0 {
		return false, fmt.Errorf("no data returned from 'transferFrom' method")
	}

	transferred, ok := ret[0].(bool)
	if !ok {
		return false, fmt.Errorf("transferFrom returned an unexpected value")
	}

	return transferred, nil
}
