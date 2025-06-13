package keeper

import (
	"fmt"
	"math/big"

	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/x/fungible/types"
)

// MUSEDepositAndCallContract deposits native MUSE to the to address if its an account or if the account does not exist yet
// If it's not an account it calls onReceive function of the connector contract and provides the address as the destinationAddress .The amount of tokens is minted to the fungible module account, wrapped and sent to the contract
func (k Keeper) MUSEDepositAndCallContract(ctx sdk.Context,
	sender ethcommon.Address,
	to ethcommon.Address,
	inboundSenderChainID int64,
	inboundAmount *big.Int,
	data []byte,
	indexBytes [32]byte) (*evmtypes.MsgEthereumTxResponse, error) {
	acc := k.evmKeeper.GetAccount(ctx, to)

	if acc == nil || !acc.IsContract() {
		err := k.DepositCoinMuse(ctx, to, inboundAmount)
		if err != nil {
			return nil, errors.Wrap(
				types.ErrDepositMuseToEvmAccount,
				fmt.Sprintf("to: %s, amount: %s err %s", to.String(), inboundAmount.String(), err.Error()),
			)
		}
		return nil, nil
	}
	// Call onReceive function of the connector contract. The connector contract will then call the onReceive function of the destination contract which is the to address
	return k.CallOnReceiveMevmConnector(
		ctx,
		sender.Bytes(),
		big.NewInt(inboundSenderChainID),
		to,
		inboundAmount,
		data,
		indexBytes,
	)
}

// MUSERevertAndCallContract deposits native MUSE to the sender address if its account or if the account does not exist yet
// If it's not an account it calls onRevert function of the connector contract and provides the sender address as the museTxSenderAddress.The amount of tokens is minted to the fungible module account, wrapped and sent to the contract
func (k Keeper) MUSERevertAndCallContract(ctx sdk.Context,
	sender ethcommon.Address,
	to ethcommon.Address,
	inboundSenderChainID int64,
	destinationChainID int64,
	remainingAmount *big.Int,
	data []byte,
	indexBytes [32]byte) (*evmtypes.MsgEthereumTxResponse, error) {
	acc := k.evmKeeper.GetAccount(ctx, sender)
	if acc == nil || !acc.IsContract() {
		err := k.DepositCoinMuse(ctx, sender, remainingAmount)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	// Call onRevert function of the connector contract. The connector contract will then call the onRevert function of the museTxSender contract which is the sender address
	return k.CallOnRevertMevmConnector(
		ctx,
		sender,
		big.NewInt(inboundSenderChainID),
		to.Bytes(),
		big.NewInt(destinationChainID),
		remainingAmount,
		data,
		indexBytes,
	)
}
