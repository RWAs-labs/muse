package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

// DeploySystemContracts deploy new instances of the system contracts
//
// Authorized: admin policy group 2.
func (k msgServer) DeploySystemContracts(
	goCtx context.Context,
	msg *types.MsgDeploySystemContracts,
) (*types.MsgDeploySystemContractsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// uniswap v2 factory
	factory, err := k.DeployUniswapV2Factory(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to deploy UniswapV2Factory")
	}

	// wmuse contract
	wmuse, err := k.DeployWMUSE(ctx)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to DeployWMuseContract")
	}

	// uniswap v2 router
	router, err := k.DeployUniswapV2Router02(ctx, factory, wmuse)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to deploy UniswapV2Router02")
	}

	// connector mevm
	connector, err := k.DeployConnectorMEVM(ctx, wmuse)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to deploy ConnectorMEVM")
	}

	// system contract
	systemContract, err := k.DeploySystemContract(ctx, wmuse, factory, router)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to deploy SystemContract")
	}

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventSystemContractsDeployed{
			MsgTypeUrl:       sdk.MsgTypeURL(&types.MsgDeploySystemContracts{}),
			UniswapV2Factory: factory.Hex(),
			Wmuse:            wmuse.Hex(),
			UniswapV2Router:  router.Hex(),
			ConnectorMevm:    connector.Hex(),
			SystemContract:   systemContract.Hex(),
			Signer:           msg.Creator,
		},
	)
	if err != nil {
		k.Logger(ctx).Error("failed to emit event",
			"event", "EventSystemContractsDeployed",
			"error", err.Error(),
		)
		return nil, cosmoserrors.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	}

	return &types.MsgDeploySystemContractsResponse{
		UniswapV2Factory: factory.Hex(),
		Wmuse:            wmuse.Hex(),
		UniswapV2Router:  router.Hex(),
		ConnectorMEVM:    connector.Hex(),
		SystemContract:   systemContract.Hex(),
	}, nil
}
