package txserver

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/RWAs-labs/muse/e2e/utils"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
)

// SetMRC20LiquidityCap sets the liquidity cap for given MRC20 token
func (zts MuseTxServer) SetMRC20LiquidityCap(
	mrc20Addr ethcommon.Address,
	liquidityCap math.Uint,
) (*sdktypes.TxResponse, error) {
	return zts.updateMRC20LiquidityCap(mrc20Addr, liquidityCap)
}

// RemoveMRC20LiquidityCap removes the liquidity cap for given MRC20 token
func (zts MuseTxServer) RemoveMRC20LiquidityCap(mrc20Addr ethcommon.Address) (*sdktypes.TxResponse, error) {
	return zts.updateMRC20LiquidityCap(mrc20Addr, math.ZeroUint())
}

// updateMRC20LiquidityCap updates the liquidity cap for given MRC20 token
func (zts MuseTxServer) updateMRC20LiquidityCap(
	mrc20Addr ethcommon.Address,
	liquidityCap math.Uint,
) (*sdktypes.TxResponse, error) {
	// create msg
	msg := fungibletypes.NewMsgUpdateMRC20LiquidityCap(
		zts.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		mrc20Addr.Hex(),
		liquidityCap,
	)

	// broadcast tx
	res, err := zts.BroadcastTx(utils.OperationalPolicyName, msg)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to set MRC20 liquidity cap for %s", mrc20Addr)
	}

	return res, nil
}
