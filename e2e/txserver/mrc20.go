package txserver

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/RWAs-labs/muse/e2e/utils"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
)

// UpdateMRC20GasLimit updates the gas limit for given MRC20 token
func (zts MuseTxServer) UpdateMRC20GasLimit(
	mrc20Addr ethcommon.Address,
	newGasLimit math.Uint,
) (*sdktypes.TxResponse, error) {
	// create msg
	msg := fungibletypes.NewMsgUpdateMRC20WithdrawFee(
		zts.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		mrc20Addr.Hex(),
		math.ZeroUint(), // 0 flat fee
		newGasLimit,
	)

	// broadcast tx
	res, err := zts.BroadcastTx(utils.OperationalPolicyName, msg)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to update MRC20 gas limit for %s", mrc20Addr)
	}

	return res, nil
}
