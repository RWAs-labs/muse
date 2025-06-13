package keeper

import (
	"math/big"

	cosmoserrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/systemcontract.sol"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/contracts/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	"github.com/RWAs-labs/muse/x/fungible/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// SetupChainGasCoinAndPool setup gas MRC20, and MUSE/gas pool for a chain
// add 0.1gas/0.1wmuse to the pool
// FIXME: add cointype and use proper gas limit based on cointype/chain
func (k Keeper) SetupChainGasCoinAndPool(
	ctx sdk.Context,
	chainID int64,
	gasAssetName string,
	symbol string,
	decimals uint8,
	gasLimit *big.Int,
	liquidityCap *sdkmath.Uint,
) (ethcommon.Address, error) {
	// additional on-chain static chain information
	additionalChains := k.GetAuthorityKeeper().GetAdditionalChainList(ctx)

	chain, found := chains.GetChainFromChainID(chainID, additionalChains)
	if !found {
		return ethcommon.Address{}, observertypes.ErrSupportedChains
	}

	transferGasLimit := gasLimit

	// Check if gas coin already exists
	_, found = k.GetGasCoinForForeignCoin(ctx, chainID)
	if found {
		return ethcommon.Address{}, types.ErrForeignCoinAlreadyExist
	}

	// default values
	if transferGasLimit == nil {
		transferGasLimit = big.NewInt(21_000)
		if chains.IsBitcoinChain(chain.ChainId, additionalChains) {
			transferGasLimit = big.NewInt(100) // 100B for a typical tx
		}
	}

	mrc20Addr, err := k.DeployMRC20Contract(
		ctx,
		gasAssetName,
		symbol,
		decimals,
		chain.ChainId,
		coin.CoinType_Gas,
		"",
		transferGasLimit,
		liquidityCap,
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to DeployMRC20Contract")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(gasAssetName, mrc20Addr.String()),
		),
	)
	err = k.SetGasCoin(ctx, big.NewInt(chain.ChainId), mrc20Addr)
	if err != nil {
		return ethcommon.Address{}, err
	}
	amount := big.NewInt(10)
	// #nosec G115 always in range
	amount.Exp(amount, big.NewInt(int64(decimals-1)), nil)
	amountAMuse := big.NewInt(1e17)

	_, err = k.DepositMRC20(ctx, mrc20Addr, types.ModuleAddressEVM, amount)
	if err != nil {
		return ethcommon.Address{}, err
	}
	err = k.bankKeeper.MintCoins(
		ctx,
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin("amuse", sdkmath.NewIntFromBigInt(amountAMuse))),
	)
	if err != nil {
		return ethcommon.Address{}, err
	}
	systemContractAddress, err := k.GetSystemContractAddress(ctx)
	if err != nil || systemContractAddress == (ethcommon.Address{}) {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			types.ErrContractNotFound,
			"system contract address invalid: %s",
			systemContractAddress,
		)
	}
	systemABI, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to get system contract abi")
	}
	_, err = k.CallEVM(
		ctx,
		*systemABI,
		types.ModuleAddressEVM,
		systemContractAddress,
		BigIntZero,
		DefaultGasLimit,
		true,
		false,
		"setGasMusePool",
		big.NewInt(chain.ChainId),
		mrc20Addr,
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			err,
			"failed to CallEVM method setGasMusePool(%d, %s)",
			chain.ChainId,
			mrc20Addr.String(),
		)
	}

	// setup uniswap v2 pools gas/muse
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to get uniswap router abi")
	}
	MRC20ABI, err := mrc20.MRC20MetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to GetAbi mrc20")
	}
	_, err = k.CallEVM(
		ctx,
		*MRC20ABI,
		types.ModuleAddressEVM,
		mrc20Addr,
		BigIntZero,
		DefaultGasLimit,
		true,
		false,
		"approve",
		routerAddress,
		amount,
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			err,
			"failed to CallEVM method approve(%s, %d)",
			routerAddress.String(),
			amount,
		)
	}

	//function addLiquidityETH(
	//	address token,
	//	uint amountTokenDesired,
	//	uint amountTokenMin,
	//	uint amountETHMin,
	//	address to,
	//	uint deadline
	//) external payable returns (uint amountToken, uint amountETH, uint liquidity);
	res, err := k.CallEVM(
		ctx,
		*routerABI,
		types.ModuleAddressEVM,
		routerAddress,
		amountAMuse,
		big.NewInt(5_000_000),
		true,
		false,
		"addLiquidityETH",
		mrc20Addr,
		amount,
		BigIntZero,
		BigIntZero,
		types.ModuleAddressEVM,
		amountAMuse,
	)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(
			err,
			"failed to CallEVM method addLiquidityETH(%s, %s)",
			mrc20Addr.String(),
			amountAMuse.String(),
		)
	}
	AmountToken := new(*big.Int)
	AmountETH := new(*big.Int)
	Liquidity := new(*big.Int)
	err = routerABI.UnpackIntoInterface(&[]interface{}{AmountToken, AmountETH, Liquidity}, "addLiquidityETH", res.Ret)
	if err != nil {
		return ethcommon.Address{}, cosmoserrors.Wrapf(err, "failed to UnpackIntoInterface addLiquidityETH")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("function", "addLiquidityETH"),
			sdk.NewAttribute("amountToken", (*AmountToken).String()),
			sdk.NewAttribute("amountETH", (*AmountETH).String()),
			sdk.NewAttribute("liquidity", (*Liquidity).String()),
		),
	)
	return mrc20Addr, nil
}
