package keeper

import (
	"context"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"

	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/ptr"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
)

// WhitelistERC20 deploys a new mrc20, create a foreign coin object for the ERC20
// and emit a crosschain tx to whitelist the ERC20 on the external chain
//
// Authorized: admin policy group 1.
func (k msgServer) WhitelistERC20(
	goCtx context.Context,
	msg *types.MsgWhitelistERC20,
) (*types.MsgWhitelistERC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if authorized
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	chain, found := k.museObserverKeeper.GetSupportedChainFromChainID(ctx, msg.ChainId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrInvalidChainID, "chain id (%d) not supported", msg.ChainId)
	}

	switch {
	case chain.IsEVMChain():
		erc20Addr := ethcommon.HexToAddress(msg.Erc20Address)
		if erc20Addr == (ethcommon.Address{}) {
			return nil, errorsmod.Wrapf(
				sdkerrors.ErrInvalidAddress,
				"invalid ERC20 contract address (%s)",
				msg.Erc20Address,
			)
		}

	case chain.IsSolanaChain():
		_, err := solana.PublicKeyFromBase58(msg.Erc20Address)
		if err != nil {
			return nil, errorsmod.Wrapf(
				sdkerrors.ErrInvalidAddress,
				"invalid solana contract address (%s)",
				msg.Erc20Address,
			)
		}

	default:
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidChainID,
			"whitelist for chain id (%d) not supported",
			msg.ChainId,
		)
	}

	// check if the asset is already whitelisted
	foreignCoins := k.fungibleKeeper.GetAllForeignCoins(ctx)
	for _, fCoin := range foreignCoins {
		if fCoin.Asset == msg.Erc20Address && fCoin.ForeignChainId == msg.ChainId {
			return nil, errorsmod.Wrapf(
				fungibletypes.ErrForeignCoinAlreadyExist,
				"ERC20 contract address (%s) already whitelisted on chain (%d)",
				msg.Erc20Address,
				msg.ChainId,
			)
		}
	}

	tss, found := k.museObserverKeeper.GetTSS(ctx)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrCannotFindTSSKeys, "Cannot create new admin cmd of type whitelistERC20")
	}

	// use a temporary context for the mrc20 deployment
	tmpCtx, commit := ctx.CacheContext()

	// add to the foreign coins. Deploy MRC20 contract for it.
	mrc20Addr, err := k.fungibleKeeper.DeployMRC20Contract(
		tmpCtx,
		msg.Name,
		msg.Symbol,
		// #nosec G115 always in range
		uint8(msg.Decimals),
		chain.ChainId,
		coin.CoinType_ERC20,
		msg.Erc20Address,
		big.NewInt(msg.GasLimit),
		ptr.Ptr(msg.LiquidityCap),
	)
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrDeployContract,
			"failed to deploy MRC20 contract for ERC20 contract address (%s) on chain (%d)",
			msg.Erc20Address,
			msg.ChainId,
		)
	}

	if mrc20Addr == (ethcommon.Address{}) {
		return nil, errorsmod.Wrapf(
			types.ErrDeployContract,
			"deployed MRC20 return 0 address for ERC20 contract address (%s) on chain (%d)",
			msg.Erc20Address,
			msg.ChainId,
		)
	}

	// get necessary parameters to create the cctx
	params, found := k.museObserverKeeper.GetChainParamsByChainID(ctx, msg.ChainId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrInvalidChainID, "chain params not found for chain id (%d)", msg.ChainId)
	}
	medianGasPrice, priorityFee, isFound := k.GetMedianGasValues(ctx, msg.ChainId)
	if !isFound {
		return nil, errorsmod.Wrapf(
			types.ErrUnableToGetGasPrice,
			"median gas price not found for chain id (%d)",
			msg.ChainId,
		)
	}

	// overpays gas price by 2x
	medianGasPrice = medianGasPrice.MulUint64(types.ERC20CustodyWhitelistGasMultiplierEVM)
	priorityFee = priorityFee.MulUint64(types.ERC20CustodyWhitelistGasMultiplierEVM)

	// should not happen
	if priorityFee.GT(medianGasPrice) {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidGasAmount,
			"priorityFee %s is greater than median gasPrice %s",
			priorityFee.String(),
			medianGasPrice.String(),
		)
	}

	// create the cctx
	cctx := types.WhitelistERC20CmdCCTX(
		msg.Creator,
		mrc20Addr,
		msg.Erc20Address,
		params.Erc20CustodyContractAddress,
		msg.ChainId,
		medianGasPrice.String(),
		priorityFee.String(),
		tss.TssPubkey,
	)

	err = k.SetObserverOutboundInfo(ctx, msg.ChainId, &cctx)
	if err != nil {
		return nil, err
	}

	// add to the foreign coins
	foreignCoin := fungibletypes.ForeignCoins{
		Mrc20ContractAddress: mrc20Addr.Hex(),
		Asset:                msg.Erc20Address,
		ForeignChainId:       msg.ChainId,
		Decimals:             msg.Decimals,
		Name:                 msg.Name,
		Symbol:               msg.Symbol,
		CoinType:             coin.CoinType_ERC20,
		// #nosec G115 always positive
		GasLimit:     uint64(msg.GasLimit),
		LiquidityCap: msg.LiquidityCap,
	}
	k.fungibleKeeper.SetForeignCoins(ctx, foreignCoin)
	k.SaveCCTXUpdate(ctx, cctx, tss.TssPubkey)

	commit()

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventERC20Whitelist{
			Mrc20Address:       mrc20Addr.Hex(),
			WhitelistCctxIndex: cctx.Index,
		},
	)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "failed to emit event")
	}

	return &types.MsgWhitelistERC20Response{
		Mrc20Address: mrc20Addr.Hex(),
		CctxIndex:    cctx.Index,
	}, nil
}
