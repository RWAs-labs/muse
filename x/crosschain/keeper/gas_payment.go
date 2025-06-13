package keeper

import (
	"fmt"
	"math/big"

	cosmoserrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// ChainGasParams represents the parameters to calculate the fees for gas for a chain.
type ChainGasParams struct {
	GasMRC20 ethcommon.Address

	GasLimit sdkmath.Uint
	GasPrice sdkmath.Uint

	// PriorityFee (optional for EIP-1559)
	PriorityFee sdkmath.Uint

	ProtocolFlatFee sdkmath.Uint
}

// PayGasAndUpdateCctx updates the outbound tx with the new amount after paying the gas fee
// **Caller should feed temporary ctx into this function**
// chainID is the outbound chain chain id , this can be receiver chain for regular transactions and sender-chain to reverted transactions
func (k Keeper) PayGasAndUpdateCctx(
	ctx sdk.Context,
	chainID int64,
	cctx *types.CrossChainTx,
	inputAmount sdkmath.Uint,
	noEthereumTxEvent bool,
) error {
	// Dispatch to the correct function based on the coin type
	switch cctx.InboundParams.CoinType {
	case coin.CoinType_Muse:
		return k.PayGasInMuseAndUpdateCctx(ctx, chainID, cctx, inputAmount, noEthereumTxEvent)
	case coin.CoinType_Gas:
		return k.PayGasNativeAndUpdateCctx(ctx, chainID, cctx, inputAmount)
	case coin.CoinType_ERC20:
		return k.PayGasInERC20AndUpdateCctx(ctx, chainID, cctx, inputAmount, noEthereumTxEvent)
	default:
		// can't pay gas with coin type
		return fmt.Errorf("can't pay gas with coin type %s", cctx.InboundParams.CoinType.String())
	}
}

// ChainGasParams returns the params to calculates the fees for gas for a chain
// tha gas address, the gas limit, gas price and protocol flat fee are returned
func (k Keeper) ChainGasParams(ctx sdk.Context, chainID int64) (ChainGasParams, error) {
	gasMRC20, err := k.fungibleKeeper.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(chainID))
	if err != nil {
		return ChainGasParams{}, errors.Wrap(err, "unable to get system contract gas coin")
	}

	// get the gas limit
	gasLimit, err := k.fungibleKeeper.QueryGasLimit(ctx, gasMRC20)
	switch {
	case err != nil:
		return ChainGasParams{}, errors.Wrap(err, "unable to get gas limit")
	case gasLimit == nil:
		return ChainGasParams{}, types.ErrInvalidGasLimit
	}

	// get the protocol flat fee
	protocolFlatFee, err := k.fungibleKeeper.QueryProtocolFlatFee(ctx, gasMRC20)
	switch {
	case err != nil:
		return ChainGasParams{}, errors.Wrap(err, "unable to get protocol flat fee")
	case protocolFlatFee == nil:
		return ChainGasParams{}, cosmoserrors.Wrap(types.ErrInvalidGasAmount, "protocol flat fee is nil")
	}

	// get the gas price
	gasPrice, priorityFee, isFound := k.GetMedianGasValues(ctx, chainID)
	if !isFound {
		return ChainGasParams{}, types.ErrUnableToGetGasPrice
	}

	return ChainGasParams{
		GasMRC20:        gasMRC20,
		GasLimit:        sdkmath.NewUintFromBigInt(gasLimit),
		GasPrice:        gasPrice,
		PriorityFee:     priorityFee,
		ProtocolFlatFee: sdkmath.NewUintFromBigInt(protocolFlatFee),
	}, nil
}

// PayGasNativeAndUpdateCctx updates the outbound tx with the new amount subtracting the gas fee
// **Caller should feed temporary ctx into this function**
func (k Keeper) PayGasNativeAndUpdateCctx(
	ctx sdk.Context,
	chainID int64,
	cctx *types.CrossChainTx,
	inputAmount sdkmath.Uint,
) error {
	// preliminary checks
	if cctx.InboundParams.CoinType != coin.CoinType_Gas {
		return cosmoserrors.Wrapf(
			types.ErrInvalidCoinType,
			"can't pay gas in native gas with %s",
			cctx.InboundParams.CoinType.String(),
		)
	}
	if _, found := k.museObserverKeeper.GetSupportedChainFromChainID(ctx, chainID); !found {
		return observertypes.ErrSupportedChains
	}

	// get gas params
	gas, err := k.ChainGasParams(ctx, chainID)
	if err != nil {
		return cosmoserrors.Wrap(types.ErrCannotFindGasParams, err.Error())
	}

	// with V2 protocol, reverts on connected chains can eventually call a onRevert function which can require a higher gas limit
	if cctx.ProtocolContractVersion == types.ProtocolContractVersion_V2 && cctx.RevertOptions.CallOnRevert &&
		!cctx.RevertOptions.RevertGasLimit.IsZero() {
		gas.GasLimit = cctx.RevertOptions.RevertGasLimit
	}

	// calculate the final gas fee
	outTxGasFee := gas.GasLimit.Mul(gas.GasPrice).Add(gas.ProtocolFlatFee)

	// subtract the withdraw fee from the input amount
	if outTxGasFee.GT(inputAmount) {
		return cosmoserrors.Wrap(
			types.ErrNotEnoughGas,
			fmt.Sprintf(
				"unable to pay for outbound tx using gas token, outbound chain: %d, required: %s, available: %s",
				chainID,
				outTxGasFee,
				inputAmount,
			),
		)
	}
	ctx.Logger().Info("Subtracting amount from inbound tx", "amount", inputAmount.String(), "fee", outTxGasFee.String())
	newAmount := inputAmount.Sub(outTxGasFee)

	// update cctx
	cctx.GetCurrentOutboundParam().Amount = newAmount
	cctx.GetCurrentOutboundParam().CallOptions.GasLimit = gas.GasLimit.Uint64()
	cctx.GetCurrentOutboundParam().GasPrice = gas.GasPrice.String()
	cctx.GetCurrentOutboundParam().GasPriorityFee = gas.PriorityFee.String()

	return nil
}

// PayGasInERC20AndUpdateCctx updates parameter cctx amount subtracting the gas fee
// the gas fee in ERC20 is calculated by swapping ERC20 -> Muse -> Gas
// if the route is not available, the gas payment will fail
// **Caller should feed temporary ctx into this function**
func (k Keeper) PayGasInERC20AndUpdateCctx(
	ctx sdk.Context,
	chainID int64,
	cctx *types.CrossChainTx,
	inputAmount sdkmath.Uint,
	noEthereumTxEvent bool,
) error {
	// preliminary checks
	if cctx.InboundParams.CoinType != coin.CoinType_ERC20 {
		return cosmoserrors.Wrapf(
			types.ErrInvalidCoinType,
			"can't pay gas in erc20 with %s",
			cctx.InboundParams.CoinType.String(),
		)
	}

	if _, found := k.museObserverKeeper.GetSupportedChainFromChainID(ctx, chainID); !found {
		return observertypes.ErrSupportedChains
	}

	// get gas params
	gas, err := k.ChainGasParams(ctx, chainID)
	if err != nil {
		return cosmoserrors.Wrap(types.ErrCannotFindGasParams, err.Error())
	}

	// with V2 protocol, reverts on connected chains can eventually call a onRevert function which can require a higher gas limit
	if cctx.ProtocolContractVersion == types.ProtocolContractVersion_V2 && cctx.RevertOptions.CallOnRevert &&
		!cctx.RevertOptions.RevertGasLimit.IsZero() {
		gas.GasLimit = cctx.RevertOptions.RevertGasLimit
	}

	outTxGasFee := gas.GasLimit.Mul(gas.GasPrice).Add(gas.ProtocolFlatFee)
	// get address of the mrc20
	fc, found := k.fungibleKeeper.GetForeignCoinFromAsset(ctx, cctx.InboundParams.Asset, chainID)
	if !found {
		return cosmoserrors.Wrapf(
			types.ErrForeignCoinNotFound,
			"mrc20 from asset %s not found",
			cctx.InboundParams.Asset,
		)
	}
	mrc20 := ethcommon.HexToAddress(fc.Mrc20ContractAddress)
	if mrc20 == (ethcommon.Address{}) {
		return cosmoserrors.Wrapf(
			types.ErrForeignCoinNotFound,
			"mrc20 from asset %s invalid address",
			cctx.InboundParams.Asset,
		)
	}

	// get the necessary ERC20 amount for gas
	feeInMRC20, err := k.fungibleKeeper.QueryUniswapV2RouterGetMRC4ToMRC4AmountsIn(
		ctx,
		outTxGasFee.BigInt(),
		mrc20,
		gas.GasMRC20,
	)
	if err != nil {
		// NOTE: this is the first method that fails when a liquidity pool is not set for the gas MRC20, so we return a specific error
		return cosmoserrors.Wrap(types.ErrNoLiquidityPool, err.Error())
	}

	// subtract the withdraw fee from the input amount
	if sdkmath.NewUintFromBigInt(feeInMRC20).GT(inputAmount) {
		return cosmoserrors.Wrap(
			types.ErrNotEnoughGas,
			fmt.Sprintf(
				"unable to pay for outbound tx using mrc20 token, outbound chain: %d, required: %s, available: %s",
				chainID,
				outTxGasFee,
				inputAmount,
			),
		)
	}
	newAmount := inputAmount.Sub(sdkmath.NewUintFromBigInt(feeInMRC20))

	// mint the amount of ERC20 to be burnt as gas fee
	_, err = k.fungibleKeeper.DepositMRC20(ctx, mrc20, types.ModuleAddressEVM, feeInMRC20)
	if err != nil {
		return cosmoserrors.Wrap(fungibletypes.ErrContractCall, err.Error())
	}
	ctx.Logger().Info("Minted ERC20 for gas fee",
		"mrc20", mrc20.Hex(),
		"amount", feeInMRC20,
	)

	// approve the uniswapv2 router to spend the ERC20
	routerAddress, err := k.fungibleKeeper.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return cosmoserrors.Wrap(fungibletypes.ErrContractCall, err.Error())
	}
	err = k.fungibleKeeper.CallMRC20Approve(
		ctx,
		types.ModuleAddressEVM,
		mrc20,
		routerAddress,
		feeInMRC20,
		noEthereumTxEvent,
	)
	if err != nil {
		return cosmoserrors.Wrap(fungibletypes.ErrContractCall, err.Error())
	}

	// swap the fee in ERC20 into gas passing through Muse and burn the gas MRC20
	amounts, err := k.fungibleKeeper.CallUniswapV2RouterSwapExactTokensForTokens(
		ctx,
		types.ModuleAddressEVM,
		types.ModuleAddressEVM,
		feeInMRC20,
		mrc20,
		gas.GasMRC20,
		noEthereumTxEvent,
	)
	if err != nil {
		return cosmoserrors.Wrap(fungibletypes.ErrContractCall, err.Error())
	}
	ctx.Logger().Info("CallUniswapV2RouterSwapExactTokensForTokens",
		"mrc20AmountIn", amounts[0],
		"gasAmountOut", amounts[2],
	)
	gasObtained := amounts[2]

	// FIXME: investigate small mismatches between gasObtained and outTxGasFee
	// https://github.com/RWAs-labs/muse/issues/1303
	// check if the final gas received after swap matches the gas fee defined
	// if not there might be issues with the pool liquidity and it is safer from an accounting perspective to return an error
	if gasObtained.Cmp(outTxGasFee.BigInt()) == -1 {
		return cosmoserrors.Wrapf(
			types.ErrInvalidGasAmount,
			"gas obtained for burn (%s) is lower than gas fee(%s)",
			gasObtained,
			outTxGasFee,
		)
	}

	// burn the gas MRC20
	err = k.fungibleKeeper.CallMRC20Burn(ctx, types.ModuleAddressEVM, gas.GasMRC20, gasObtained, noEthereumTxEvent)
	if err != nil {
		return cosmoserrors.Wrap(fungibletypes.ErrContractCall, err.Error())
	}
	ctx.Logger().Info("Burning gas MRC20",
		"mrc20", gas.GasMRC20.Hex(),
		"amount", gasObtained,
	)

	// update cctx
	cctx.GetCurrentOutboundParam().Amount = newAmount
	cctx.GetCurrentOutboundParam().CallOptions.GasLimit = gas.GasLimit.Uint64()
	cctx.GetCurrentOutboundParam().GasPrice = gas.GasPrice.String()
	cctx.GetCurrentOutboundParam().GasPriorityFee = gas.PriorityFee.String()

	return nil
}

// PayGasInMuseAndUpdateCctx updates parameter cctx with the gas price and gas fee for the outbound tx;
// it also makes a trade to fulfill the outbound tx gas fee in MUSE by swapping MUSE for some gas MRC20 balances
// The gas MRC20 balance is subsequently burned to account for the expense of TSS address gas fee payment in the outbound tx.
// museBurnt represents the amount of Muse that has been burnt for the tx, the final amount for the tx is museBurnt - gasFee
// **Caller should feed temporary ctx into this function**
func (k Keeper) PayGasInMuseAndUpdateCctx(
	ctx sdk.Context,
	chainID int64,
	cctx *types.CrossChainTx,
	museBurnt sdkmath.Uint,
	noEthereumTxEvent bool,
) error {
	// preliminary checks
	if cctx.InboundParams.CoinType != coin.CoinType_Muse {
		return cosmoserrors.Wrapf(
			types.ErrInvalidCoinType,
			"can't pay gas in muse with %s",
			cctx.InboundParams.CoinType.String(),
		)
	}

	if _, found := k.museObserverKeeper.GetSupportedChainFromChainID(ctx, chainID); !found {
		return observertypes.ErrSupportedChains
	}

	gasMRC20, err := k.fungibleKeeper.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(chainID))
	if err != nil {
		return cosmoserrors.Wrapf(
			err,
			"PayGasInMuseAndUpdateCctx: unable to get system contract gas coin, chaind ID %d",
			chainID,
		)
	}

	// get the gas price
	gasPrice, priorityFee, isFound := k.GetMedianGasValues(ctx, chainID)
	if !isFound {
		return cosmoserrors.Wrapf(types.ErrUnableToGetGasPrice,
			"chain %d",
			chainID,
		)
	}
	// overpays gas price
	const multiplier = 2
	gasPrice = gasPrice.MulUint64(multiplier)
	priorityFee = priorityFee.MulUint64(multiplier)

	// should not happen
	if priorityFee.GT(gasPrice) {
		return cosmoserrors.Wrapf(
			types.ErrInvalidGasAmount,
			"priorityFee %s is greater than gasPrice %s",
			priorityFee.String(),
			gasPrice.String(),
		)
	}

	// get the gas fee in gas token
	gasLimit := sdkmath.NewUint(cctx.GetCurrentOutboundParam().CallOptions.GasLimit)
	outTxGasFee := gasLimit.Mul(gasPrice)

	// get the gas fee in Muse using system uniswapv2 pool wmuse/gasMRC20 and adding the protocol fee
	outTxGasFeeInMuse, err := k.fungibleKeeper.QueryUniswapV2RouterGetMuseAmountsIn(ctx, outTxGasFee.BigInt(), gasMRC20)
	if err != nil {
		return cosmoserrors.Wrap(err, "PayGasInMuseAndUpdateCctx: unable to QueryUniswapV2RouterGetMuseAmountsIn")
	}
	feeInMuse := types.GetProtocolFee().Add(sdkmath.NewUintFromBigInt(outTxGasFeeInMuse))

	// reduce the amount of the outbound tx
	if feeInMuse.GT(museBurnt) {
		return cosmoserrors.Wrap(
			types.ErrNotEnoughGas,
			fmt.Sprintf(
				"unable to pay for outbound tx using muse token, outbound chain: %d, required: %s, available: %s",
				chainID,
				outTxGasFee,
				museBurnt,
			),
		)
	}
	ctx.Logger().Info("Subtracting amount from inbound tx",
		"amount", museBurnt.String(),
		"feeInMuse", feeInMuse.String(),
	)
	newAmount := museBurnt.Sub(feeInMuse)

	// ** The following logic converts the outTxGasFeeInMuse into gasMRC20 and burns it **
	// swap the outTxGasFeeInMuse portion of muse to the real gas MRC20 and burn it, in a temporary context.
	{
		coins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdkmath.NewIntFromBigInt(feeInMuse.BigInt())))
		err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
		if err != nil {
			return cosmoserrors.Wrap(err, "PayGasInMuseAndUpdateCctx: unable to mint coins")
		}

		amounts, err := k.fungibleKeeper.CallUniswapV2RouterSwapExactETHForToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			outTxGasFeeInMuse,
			gasMRC20,
			noEthereumTxEvent,
		)
		if err != nil {
			return cosmoserrors.Wrap(
				err,
				"PayGasInMuseAndUpdateCctx: unable to CallUniswapv2RouterSwapExactETHForToken",
			)
		}

		ctx.Logger().Info("gas fee", "outTxGasFee", outTxGasFee, "outTxGasFeeInMuse", outTxGasFeeInMuse)
		ctx.Logger().Info("CallUniswapv2RouterSwapExactETHForToken",
			"museAmountIn", amounts[0],
			"mrc20AmountOut", amounts[1],
		)

		// FIXME: investigate small mismatches between amounts[1] and outTxGasFee
		// https://github.com/RWAs-labs/muse/issues/1303
		err = k.fungibleKeeper.CallMRC20Burn(ctx, types.ModuleAddressEVM, gasMRC20, amounts[1], noEthereumTxEvent)
		if err != nil {
			return cosmoserrors.Wrap(err, "PayGasInMuseAndUpdateCctx: unable to CallMRC20Burn")
		}
	}

	// Update the cctx
	cctx.GetCurrentOutboundParam().GasPrice = gasPrice.String()
	cctx.GetCurrentOutboundParam().GasPriorityFee = priorityFee.String()
	cctx.GetCurrentOutboundParam().Amount = newAmount
	if cctx.MuseFees.IsNil() {
		cctx.MuseFees = feeInMuse
	} else {
		cctx.MuseFees = cctx.MuseFees.Add(feeInMuse)
	}

	return nil
}
