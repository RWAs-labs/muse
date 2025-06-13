package keeper

import (
	"fmt"
	"math/big"

	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/systemcontract.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/wmuse.sol"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/RWAs-labs/muse/pkg/contracts/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

// SetSystemContract set system contract in the store
func (k Keeper) SetSystemContract(ctx sdk.Context, sytemContract types.SystemContract) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SystemContractKey))
	b := k.cdc.MustMarshal(&sytemContract)
	store.Set([]byte{0}, b)
}

// GetSystemContract returns system contract from the store
func (k Keeper) GetSystemContract(ctx sdk.Context) (val types.SystemContract, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SystemContractKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveSystemContract removes system contract from the store
func (k Keeper) RemoveSystemContract(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SystemContractKey))
	store.Delete([]byte{0})
}

// GetSystemContractAddress returns the system contract address
// TODO : wmuseContractAddress and other constant strings , can be declared as a constant string in types
// TODO Remove repetitive code
func (k *Keeper) GetSystemContractAddress(ctx sdk.Context) (ethcommon.Address, error) {
	// set the system contract
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrStateVariableNotFound,
			"failed to get system contract variable",
		)
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	return systemAddress, nil
}

// GetWMuseContractAddress returns the wmuse contract address on MuseChain
func (k *Keeper) GetWMuseContractAddress(ctx sdk.Context) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrStateVariableNotFound,
			"failed to get system contract variable",
		)
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to get system contract abi")
	}

	res, err := k.CallEVM(
		ctx,
		*sysABI,
		types.ModuleAddressEVM,
		systemAddress,
		BigIntZero,
		nil,
		false,
		false,
		"wMuseContractAddress",
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to call wMuseContractAddress (%s)",
			err.Error(),
		)
	}
	type AddressResponse struct {
		Value ethcommon.Address
	}
	var wmuseResponse AddressResponse
	if err := sysABI.UnpackIntoInterface(&wmuseResponse, "wMuseContractAddress", res.Ret); err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrABIUnpack,
			"failed to unpack wMuseContractAddress: %s",
			err.Error(),
		)
	}

	if wmuseResponse.Value == (ethcommon.Address{}) {
		return ethcommon.Address{}, cosmoserrors.Wrapf(types.ErrContractNotFound, "wmuse contract invalid address")
	}
	return wmuseResponse.Value, nil
}

// GetUniswapV2FactoryAddress returns the uniswapv2 factory contract address on MuseChain
func (k *Keeper) GetUniswapV2FactoryAddress(ctx sdk.Context) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrStateVariableNotFound,
			"failed to get system contract variable",
		)
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to get system contract abi")
	}

	res, err := k.CallEVM(
		ctx,
		*sysABI,
		types.ModuleAddressEVM,
		systemAddress,
		BigIntZero,
		nil,
		false,
		false,
		"uniswapv2FactoryAddress",
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to call uniswapv2FactoryAddress (%s)",
			err.Error(),
		)
	}
	type AddressResponse struct {
		Value ethcommon.Address
	}
	var uniswapFactoryResponse AddressResponse
	if err := sysABI.UnpackIntoInterface(&uniswapFactoryResponse, "uniswapv2FactoryAddress", res.Ret); err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrABIUnpack,
			"failed to unpack uniswapv2FactoryAddress: %s",
			err.Error(),
		)
	}

	if uniswapFactoryResponse.Value == (ethcommon.Address{}) {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrContractNotFound,
			"uniswap factory contract invalid address",
		)
	}
	return uniswapFactoryResponse.Value, nil
}

// GetUniswapV2Router02Address returns the uniswapv2 router02 address on MuseChain
func (k *Keeper) GetUniswapV2Router02Address(ctx sdk.Context) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrStateVariableNotFound,
			"failed to get system contract variable",
		)
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to get system contract abi")
	}

	res, err := k.CallEVM(
		ctx,
		*sysABI,
		types.ModuleAddressEVM,
		systemAddress,
		BigIntZero,
		nil,
		false,
		false,
		"uniswapv2Router02Address",
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to call uniswapv2Router02Address (%s)",
			err.Error(),
		)
	}
	type AddressResponse struct {
		Value ethcommon.Address
	}
	var routerResponse AddressResponse
	if err := sysABI.UnpackIntoInterface(&routerResponse, "uniswapv2Router02Address", res.Ret); err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrABIUnpack,
			"failed to unpack uniswapv2Router02Address: %s",
			err.Error(),
		)
	}

	if routerResponse.Value == (ethcommon.Address{}) {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrContractNotFound,
			"uniswap router contract invalid address",
		)
	}
	return routerResponse.Value, nil
}

// CallWMuseDeposit calls the deposit method of the wmuse contract
func (k *Keeper) CallWMuseDeposit(ctx sdk.Context, sender ethcommon.Address, amount *big.Int) error {
	wmuseAddress, err := k.GetWMuseContractAddress(ctx)
	if err != nil {
		return cosmoserrors.Wrapf(err, "failed to get wmuse contract address")
	}
	abi, err := wmuse.WETH9MetaData.GetAbi()
	if err != nil {
		return err
	}
	gasLimit := big.NewInt(70_000) // for some reason, GasEstimate for this contract call is always insufficient

	_, err = k.CallEVM(
		ctx,
		*abi,
		sender,
		wmuseAddress,
		amount,
		gasLimit,
		true,
		false,
		"deposit",
	)
	if err != nil {
		return cosmoserrors.Wrapf(types.ErrContractCall, "failed to call wmuse deposit (%s)", err.Error())
	}
	return nil
}

// QueryWMuseBalanceOf returns the balance of the given address in the wmuse contract
func (k *Keeper) QueryWMuseBalanceOf(ctx sdk.Context, addr ethcommon.Address) (*big.Int, error) {
	wmuseAddress, err := k.GetWMuseContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get wmuse contract address")
	}

	wmuseABI, err := wmuse.WETH9MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get ABI")
	}

	res, err := k.CallEVM(
		ctx,
		*wmuseABI,
		addr,
		wmuseAddress,
		big.NewInt(0),
		nil,
		false,
		false,
		"balanceOf",
		addr,
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrContractCall, "failed to call balanceOf (%s)", err.Error())
	}

	type BigIntResponse struct {
		Value *big.Int
	}
	var balanceResponse BigIntResponse
	if err := wmuseABI.UnpackIntoInterface(&balanceResponse, "balanceOf", res.Ret); err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrABIUnpack, "failed to unpack balanceOf: %s", err.Error())
	}

	return balanceResponse.Value, nil
}

// QuerySystemContractGasCoinMRC20 returns the gas coin mrc20 address for the given chain id
func (k *Keeper) QuerySystemContractGasCoinMRC20(ctx sdk.Context, chainid *big.Int) (ethcommon.Address, error) {
	system, found := k.GetSystemContract(ctx)
	if !found {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrStateVariableNotFound,
			"failed to get system contract variable",
		)
	}
	systemAddress := ethcommon.HexToAddress(system.SystemContract)
	sysABI, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to get system contract abi")
	}

	res, err := k.CallEVM(
		ctx,
		*sysABI,
		types.ModuleAddressEVM,
		systemAddress,
		BigIntZero,
		nil,
		false,
		false,
		"gasCoinMRC20ByChainId",
		chainid,
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to call gasCoinMRC20ByChainId (%s)",
			err.Error(),
		)
	}

	type AddressResponse struct {
		Value ethcommon.Address
	}
	var mrc20Res AddressResponse
	if err := sysABI.UnpackIntoInterface(&mrc20Res, "gasCoinMRC20ByChainId", res.Ret); err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrABIUnpack,
			"failed to unpack gasCoinMRC20ByChainId: %s",
			err.Error(),
		)
	}
	if mrc20Res.Value == (ethcommon.Address{}) {
		return ethcommon.Address{}, cosmoserrors.Wrapf(types.ErrContractNotFound, "gas coin contract invalid address")
	}
	return mrc20Res.Value, nil
}

// CallUniswapV2RouterSwapExactTokensForTokens calls the swapExactTokensForETH method of the uniswapv2 router contract
// to swap tokens to another tokens using wMuse as intermediary
func (k *Keeper) CallUniswapV2RouterSwapExactTokensForTokens(
	ctx sdk.Context,
	sender ethcommon.Address,
	to ethcommon.Address,
	amountIn *big.Int,
	inMRC4,
	outMRC4 ethcommon.Address,
	noEthereumTxEvent bool,
) (ret []*big.Int, err error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}
	wmuseAddr, err := k.GetWMuseContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWMuseContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}

	//function swapExactTokensForTokens(
	//	uint amountIn,
	//	uint amountOutMin,
	//	address[] calldata path,
	//	address to,
	//	uint deadline
	//)
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		sender,
		routerAddress,
		BigIntZero,
		big.NewInt(1_000_000),
		true,
		noEthereumTxEvent,
		"swapExactTokensForTokens",
		amountIn,
		BigIntZero,
		[]ethcommon.Address{inMRC4, wmuseAddr, outMRC4},
		to,
		big.NewInt(1e17),
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to CallEVM method swapExactTokensForTokens (%s)",
			err.Error(),
		)
	}

	amounts := new([3]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "swapExactTokensForTokens", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to UnpackIntoInterface swapExactTokensForTokens")
	}
	return (*amounts)[:], nil
}

// CallUniswapV2RouterSwapExactTokensForETH calls the swapExactTokensForETH method of the uniswapv2 router contract
func (k *Keeper) CallUniswapV2RouterSwapExactTokensForETH(
	ctx sdk.Context,
	sender ethcommon.Address,
	to ethcommon.Address,
	amountIn *big.Int,
	inMRC4 ethcommon.Address,
	noEthereumTxEvent bool,
) (ret []*big.Int, err error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}
	wmuseAddr, err := k.GetWMuseContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWMuseContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}

	//function swapExactTokensForETH(
	//	uint amountIn,
	//	uint amountOutMin,
	//	address[] calldata path,
	//	address to,
	//	uint deadline
	//)
	ctx.Logger().Error("Calling swapExactTokensForETH")
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		sender,
		routerAddress,
		BigIntZero,
		big.NewInt(300_000),
		true,
		noEthereumTxEvent,
		"swapExactTokensForETH",
		amountIn,
		BigIntZero,
		[]ethcommon.Address{inMRC4, wmuseAddr},
		to,
		big.NewInt(1e17),
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to CallEVM method swapExactTokensForETH (%s)",
			err.Error(),
		)
	}

	amounts := new([2]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "swapExactTokensForETH", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to UnpackIntoInterface swapExactTokensForETH")
	}
	return (*amounts)[:], nil
}

// CallUniswapV2RouterSwapExactETHForToken calls the swapExactETHForTokens method of the uniswapv2 router contract
func (k *Keeper) CallUniswapV2RouterSwapExactETHForToken(
	ctx sdk.Context,
	sender ethcommon.Address,
	to ethcommon.Address,
	amountIn *big.Int,
	outMRC4 ethcommon.Address,
	noEthereumTxEvent bool,
) ([]*big.Int, error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}

	wmuseAddr, err := k.GetWMuseContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWMuseContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}

	//function swapExactETHForTokens(uint amountOutMin, address[] calldata path, address to, uint deadline)  external payable
	//returns (uint[] memory amounts);
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		sender,
		routerAddress,
		amountIn,
		big.NewInt(300_000),
		true,
		noEthereumTxEvent,
		"swapExactETHForTokens",
		BigIntZero,
		[]ethcommon.Address{wmuseAddr, outMRC4},
		to,
		big.NewInt(1e17),
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to CallEVM method swapExactETHForTokens (%s)",
			err.Error(),
		)
	}

	amounts := new([2]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "swapExactETHForTokens", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to UnpackIntoInterface swapExactETHForTokens")
	}
	return (*amounts)[:], nil
}

// CallUniswapV2RouterSwapEthForExactToken calls the swapETHForExactTokens method of the uniswapv2 router contract
func (k *Keeper) CallUniswapV2RouterSwapEthForExactToken(
	ctx sdk.Context,
	sender ethcommon.Address,
	to ethcommon.Address,
	maxAmountIn *big.Int,
	amountOut *big.Int,
	outMRC4 ethcommon.Address,
) ([]*big.Int, error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}
	wmuseAddr, err := k.GetWMuseContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWMuseContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}

	//function swapETHForExactTokens(uint amountOut, address[] calldata path, address to, uint deadline)
	//returns (uint[] memory amounts);
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		sender,
		routerAddress,
		maxAmountIn,
		big.NewInt(300_000),
		true,
		false,
		"swapETHForExactTokens",
		amountOut,
		[]ethcommon.Address{wmuseAddr, outMRC4},
		to,
		big.NewInt(1e17),
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(
			types.ErrContractCall,
			"failed to CallEVM method swapETHForExactTokens (%s)",
			err.Error(),
		)
	}

	amounts := new([2]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "swapETHForExactTokens", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to unpack swapETHForExactTokens")
	}
	return (*amounts)[:], nil
}

// QueryUniswapV2RouterGetMuseAmountsIn returns the amount of muse needed to buy the given amount of MRC4 tokens
func (k *Keeper) QueryUniswapV2RouterGetMuseAmountsIn(
	ctx sdk.Context,
	amountOut *big.Int,
	outMRC4 ethcommon.Address,
) (*big.Int, error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}
	wmuseAddr, err := k.GetWMuseContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWMuseContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}

	//function getAmountsIn(uint amountOut, address[] memory path) public view returns (uint[] memory amounts);
	k.Logger(ctx).Info("getAmountsIn", "outMRC20", outMRC4.Hex(), "amountOut", amountOut, "wmuse", wmuseAddr.Hex())
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		types.ModuleAddressEVM,
		routerAddress,
		BigIntZero,
		nil,
		false,
		false,
		"getAmountsIn",
		amountOut,
		[]ethcommon.Address{wmuseAddr, outMRC4},
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, cosmoserrors.Wrap(
			types.ErrContractCall,
			fmt.Sprintf("failed to CallEVM method getAmountsIn (%s)", err.Error()),
		)
	}

	amounts := new([2]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "getAmountsIn", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to unpack getAmountsIn")
	}
	return (*amounts)[0], nil
}

// QueryUniswapV2RouterGetMRC4AmountsIn returns the amount of MRC4 tokens needed to buy the given amount of muse
func (k *Keeper) QueryUniswapV2RouterGetMRC4AmountsIn(
	ctx sdk.Context,
	amountOut *big.Int,
	inMRC4 ethcommon.Address,
) (*big.Int, error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}
	wmuseAddr, err := k.GetWMuseContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWMuseContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}

	//function getAmountsIn(uint amountOut, address[] memory path) public view returns (uint[] memory amounts);
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		types.ModuleAddressEVM,
		routerAddress,
		BigIntZero,
		nil,
		false,
		false,
		"getAmountsIn",
		amountOut,
		[]ethcommon.Address{inMRC4, wmuseAddr},
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrContractCall, "failed to CallEVM method getAmountsIn (%s)", err.Error())
	}

	amounts := new([2]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "getAmountsIn", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to unpack getAmountsIn")
	}
	return (*amounts)[0], nil
}

// QueryUniswapV2RouterGetMRC4ToMRC4AmountsIn returns the amount of MRC4 tokens needed to buy another MRC4 token, it uses the WMuse contract as a bridge
func (k *Keeper) QueryUniswapV2RouterGetMRC4ToMRC4AmountsIn(
	ctx sdk.Context,
	amountOut *big.Int,
	inMRC4, outMRC4 ethcommon.Address,
) (*big.Int, error) {
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to get router abi")
	}
	wmuseAddr, err := k.GetWMuseContractAddress(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetWMuseContractAddress")
	}
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}

	//function getAmountsIn(uint amountOut, address[] memory path) public view returns (uint[] memory amounts);
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		types.ModuleAddressEVM,
		routerAddress,
		BigIntZero,
		nil,
		false,
		false,
		"getAmountsIn",
		amountOut,
		[]ethcommon.Address{inMRC4, wmuseAddr, outMRC4},
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrContractCall, "failed to CallEVM method getAmountsIn (%s)", err.Error())
	}

	amounts := new([3]*big.Int)
	err = routerABI.UnpackIntoInterface(&amounts, "getAmountsIn", res.Ret)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to unpack getAmountsIn")
	}
	return (*amounts)[0], nil
}

// CallMRC20Burn calls the burn method of the mrc20 contract
func (k *Keeper) CallMRC20Burn(
	ctx sdk.Context,
	sender ethcommon.Address,
	mrc20address ethcommon.Address,
	amount *big.Int,
	noEthereumTxEvent bool,
) error {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	if err != nil {
		return cosmoserrors.Wrapf(err, "failed to get mrc20 abi")
	}

	_, err = k.CallEVM(
		ctx,
		*mrc20ABI,
		sender,
		mrc20address,
		big.NewInt(0),
		DefaultGasLimit,
		true,
		noEthereumTxEvent,
		"burn",
		amount,
	)
	if err != nil {
		return cosmoserrors.Wrapf(types.ErrContractCall, "failed to CallEVM method burn (%s)", err.Error())
	}

	return nil
}

// CallMRC20Deposit calls the deposit method of the mrc20 contract
func (k *Keeper) CallMRC20Deposit(
	ctx sdk.Context,
	sender ethcommon.Address,
	mrc20address ethcommon.Address,
	to ethcommon.Address,
	amount *big.Int,
) error {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	if err != nil {
		return cosmoserrors.Wrapf(err, "failed to get mrc20 abi")
	}

	_, err = k.CallEVM(
		ctx,
		*mrc20ABI,
		sender,
		mrc20address,
		big.NewInt(0),
		DefaultGasLimit,
		true,
		false,
		"deposit",
		to,
		amount,
	)
	if err != nil {
		return cosmoserrors.Wrapf(types.ErrContractCall, "failed to CallEVM method burn (%s)", err.Error())
	}
	return nil
}

// CallMRC20Approve calls the approve method of the mrc20 contract
func (k *Keeper) CallMRC20Approve(
	ctx sdk.Context,
	owner ethcommon.Address,
	mrc20address ethcommon.Address,
	spender ethcommon.Address,
	amount *big.Int,
	noEthereumTxEvent bool,
) error {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	if err != nil {
		return cosmoserrors.Wrapf(err, "failed to get mrc20 abi")
	}

	_, err = k.CallEVM(
		ctx,
		*mrc20ABI,
		owner,
		mrc20address,
		BigIntZero,
		DefaultGasLimit,
		true,
		noEthereumTxEvent,
		"approve",
		spender,
		amount,
	)
	if err != nil {
		return cosmoserrors.Wrapf(types.ErrContractCall, "failed to CallEVM method approve (%s)", err.Error())
	}

	return nil
}
