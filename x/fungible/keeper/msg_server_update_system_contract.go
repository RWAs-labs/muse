package keeper

import (
	"context"
	"math/big"

	cosmoserrors "cosmossdk.io/errors"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/systemcontract.sol"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/RWAs-labs/muse/pkg/coin"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

// UpdateSystemContract updates the system contract
func (k msgServer) UpdateSystemContract(
	goCtx context.Context,
	msg *types.MsgUpdateSystemContract,
) (*types.MsgUpdateSystemContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}
	newSystemContractAddr := ethcommon.HexToAddress(msg.NewSystemContractAddress)
	if newSystemContractAddr == (ethcommon.Address{}) {
		return nil, cosmoserrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid system contract address (%s)",
			msg.NewSystemContractAddress,
		)
	}

	// update contracts
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrABIGet, "failed to get mrc20 abi")
	}
	sysABI, err := systemcontract.SystemContractMetaData.GetAbi()
	if err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrABIGet, "failed to get system contract abi")
	}
	foreignCoins := k.GetAllForeignCoins(ctx)
	tmpCtx, commit := ctx.CacheContext()
	for _, fcoin := range foreignCoins {
		mrc20Addr := ethcommon.HexToAddress(fcoin.Mrc20ContractAddress)
		if mrc20Addr == (ethcommon.Address{}) {
			k.Logger(ctx).Error("invalid mrc20 contract address", "address", fcoin.Mrc20ContractAddress)
			continue
		}
		_, err = k.CallEVM(
			tmpCtx,
			*mrc20ABI,
			types.ModuleAddressEVM,
			mrc20Addr,
			BigIntZero,
			DefaultGasLimit,
			true,
			false,
			"updateSystemContractAddress",
			newSystemContractAddr,
		)
		if err != nil {
			return nil, cosmoserrors.Wrapf(
				types.ErrContractCall,
				"failed to call mrc20 contract method updateSystemContractAddress (%s)",
				err.Error(),
			)
		}
		if fcoin.CoinType == coin.CoinType_Gas {
			_, err = k.CallEVM(
				tmpCtx,
				*sysABI,
				types.ModuleAddressEVM,
				newSystemContractAddr,
				BigIntZero,
				DefaultGasLimit,
				true,
				false,
				"setGasCoinMRC20",
				big.NewInt(fcoin.ForeignChainId),
				mrc20Addr,
			)
			if err != nil {
				return nil, cosmoserrors.Wrapf(
					types.ErrContractCall,
					"failed to call system contract method setGasCoinMRC20 (%s)",
					err.Error(),
				)
			}
			_, err = k.CallEVM(
				tmpCtx,
				*sysABI,
				types.ModuleAddressEVM,
				newSystemContractAddr,
				BigIntZero,
				DefaultGasLimit,
				true,
				false,
				"setGasMusePool",
				big.NewInt(fcoin.ForeignChainId),
				mrc20Addr,
			)
			if err != nil {
				return nil, cosmoserrors.Wrapf(
					types.ErrContractCall,
					"failed to call system contract method setGasMusePool (%s)",
					err.Error(),
				)
			}
		}
	}

	sys, found := k.GetSystemContract(ctx)
	if !found {
		k.Logger(ctx).Error("system contract not found")
	}
	oldSystemContractAddress := sys.SystemContract
	sys.SystemContract = newSystemContractAddr.Hex()
	k.SetSystemContract(ctx, sys)
	err = ctx.EventManager().EmitTypedEvent(
		&types.EventSystemContractUpdated{
			MsgTypeUrl:         sdk.MsgTypeURL(&types.MsgUpdateSystemContract{}),
			NewContractAddress: msg.NewSystemContractAddress,
			OldContractAddress: oldSystemContractAddress,
			Signer:             msg.Creator,
		},
	)
	if err != nil {
		k.Logger(ctx).Error("failed to emit event", "error", err.Error())
		return nil, cosmoserrors.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	}
	commit()
	return &types.MsgUpdateSystemContractResponse{}, nil
}
