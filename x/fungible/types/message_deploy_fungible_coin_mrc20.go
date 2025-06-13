package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/RWAs-labs/muse/pkg/coin"
)

const TypeMsgDeployFungibleCoinMRC20 = "deploy_fungible_coin_mrc_20"

var _ sdk.Msg = &MsgDeployFungibleCoinMRC20{}

func NewMsgDeployFungibleCoinMRC20(
	creator string,
	ERC20 string,
	foreignChainID int64,
	decimals uint32,
	name string,
	symbol string,
	coinType coin.CoinType,
	gasLimit int64,
	liquidityCap *sdkmath.Uint,
) *MsgDeployFungibleCoinMRC20 {
	return &MsgDeployFungibleCoinMRC20{
		Creator:        creator,
		ERC20:          ERC20,
		ForeignChainId: foreignChainID,
		Decimals:       decimals,
		Name:           name,
		Symbol:         symbol,
		CoinType:       coinType,
		GasLimit:       gasLimit,
		LiquidityCap:   liquidityCap,
	}
}

func (msg *MsgDeployFungibleCoinMRC20) Route() string {
	return RouterKey
}

func (msg *MsgDeployFungibleCoinMRC20) Type() string {
	return TypeMsgDeployFungibleCoinMRC20
}

func (msg *MsgDeployFungibleCoinMRC20) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeployFungibleCoinMRC20) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeployFungibleCoinMRC20) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.GasLimit < 0 {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidGasLimit, "invalid gas limit")
	}
	if msg.Decimals > 77 {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "decimals must be less than 78")
	}
	if msg.LiquidityCap != nil && msg.LiquidityCap.IsNil() {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "liquidity cap is nil")
	}

	return nil
}
