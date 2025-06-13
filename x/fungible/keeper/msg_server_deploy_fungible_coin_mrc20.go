package keeper

import (
	"context"
	"math/big"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/RWAs-labs/muse/pkg/coin"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

// DeployFungibleCoinMRC20 deploys a fungible coin from a connected chains as a MRC20 on MuseChain.
//
// If this is a gas coin, the following happens:
//
// * MRC20 contract for the coin is deployed
// * contract address of MRC20 is set as a token address in the system
// contract
// * MUSE tokens are minted and deposited into the module account
// * setGasMusePool is called on the system contract to add the information
// about the pool to the system contract
// * addLiquidityETH is called to add liquidity to the pool
//
// If this is a non-gas coin, the following happens:
//
// * MRC20 contract for the coin is deployed
// * The coin is added to the list of foreign coins in the module's state
//
// Authorized: admin policy group 2.
func (k msgServer) DeployFungibleCoinMRC20(
	goCtx context.Context,
	msg *types.MsgDeployFungibleCoinMRC20,
) (*types.MsgDeployFungibleCoinMRC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var address common.Address
	var err error

	if err = msg.ValidateBasic(); err != nil {
		return nil, err
	}

	err = k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	if msg.CoinType == coin.CoinType_Gas {
		// #nosec G115 always in range
		address, err = k.SetupChainGasCoinAndPool(
			ctx,
			msg.ForeignChainId,
			msg.Name,
			msg.Symbol,
			uint8(msg.Decimals),
			big.NewInt(msg.GasLimit),
			msg.LiquidityCap,
		)
		if err != nil {
			return nil, cosmoserrors.Wrapf(err, "failed to setupChainGasCoinAndPool")
		}
	} else {
		// #nosec G115 always in range
		address, err = k.DeployMRC20Contract(
			ctx,
			msg.Name,
			msg.Symbol,
			uint8(msg.Decimals),
			msg.ForeignChainId,
			msg.CoinType,
			msg.ERC20,
			big.NewInt(msg.GasLimit),
			msg.LiquidityCap,
		)
		if err != nil {
			return nil, err
		}
	}

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventMRC20Deployed{
			MsgTypeUrl: sdk.MsgTypeURL(&types.MsgDeployFungibleCoinMRC20{}),
			ChainId:    msg.ForeignChainId,
			Contract:   address.String(),
			Name:       msg.Name,
			Symbol:     msg.Symbol,
			// #nosec G115 always in range
			Decimals: int64(msg.Decimals),
			CoinType: msg.CoinType,
			Erc20:    msg.ERC20,
			GasLimit: msg.GasLimit,
		},
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to emit event")
	}

	return &types.MsgDeployFungibleCoinMRC20Response{
		Address: address.Hex(),
	}, nil
}
