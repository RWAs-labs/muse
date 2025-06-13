package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgDeployFungibleCoinMRC20{}, "fungible/DeployFungibleCoinMRC20", nil)
	cdc.RegisterConcrete(&MsgDeploySystemContracts{}, "fungible/MsgDeploySystemContracts", nil)
	cdc.RegisterConcrete(&MsgRemoveForeignCoin{}, "fungible/RemoveForeignCoin", nil)
	cdc.RegisterConcrete(&MsgUpdateSystemContract{}, "fungible/UpdateSystemContract", nil)
	cdc.RegisterConcrete(&MsgUpdateMRC20WithdrawFee{}, "fungible/UpdateMRC20WithdrawFee", nil)
	cdc.RegisterConcrete(&MsgUpdateContractBytecode{}, "fungible/UpdateContractBytecode", nil)
	cdc.RegisterConcrete(&MsgUpdateMRC20LiquidityCap{}, "fungible/UpdateMRC20LiquidityCap", nil)
	cdc.RegisterConcrete(&MsgPauseMRC20{}, "fungible/PauseMRC20", nil)
	cdc.RegisterConcrete(&MsgUnpauseMRC20{}, "fungible/UnpauseMRC20", nil)
	cdc.RegisterConcrete(&MsgUpdateGatewayContract{}, "fungible/UpdateGatewayContract", nil)
	cdc.RegisterConcrete(&MsgUpdateMRC20Name{}, "fungible/UpdateMRC20Name", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeployFungibleCoinMRC20{},
		&MsgDeploySystemContracts{},
		&MsgRemoveForeignCoin{},
		&MsgUpdateSystemContract{},
		&MsgUpdateMRC20WithdrawFee{},
		&MsgUpdateContractBytecode{},
		&MsgUpdateMRC20LiquidityCap{},
		&MsgPauseMRC20{},
		&MsgUnpauseMRC20{},
		&MsgUpdateGatewayContract{},
		&MsgUpdateMRC20Name{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
