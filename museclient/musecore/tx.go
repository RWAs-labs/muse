package musecore

import (
	"context"
	"strings"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"

	clientauthz "github.com/RWAs-labs/muse/museclient/authz"
	clientcommon "github.com/RWAs-labs/muse/museclient/common"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// GetInboundVoteMessage returns a new MsgVoteInbound
// TODO(revamp): move to a different file
func GetInboundVoteMessage(
	sender string,
	senderChain int64,
	txOrigin string,
	receiver string,
	receiverChain int64,
	amount math.Uint,
	message string,
	inboundHash string,
	inBlockHeight uint64,
	gasLimit uint64,
	coinType coin.CoinType,
	asset string,
	signerAddress string,
	eventIndex uint64,
	status types.InboundStatus,
) *types.MsgVoteInbound {
	msg := types.NewMsgVoteInbound(
		signerAddress,
		sender,
		senderChain,
		txOrigin,
		receiver,
		receiverChain,
		amount,
		message,
		inboundHash,
		inBlockHeight,
		gasLimit,
		coinType,
		asset,
		eventIndex,
		types.ProtocolContractVersion_V1,
		false, // not relevant for v1
		status,
		types.ConfirmationMode_SAFE,
	)
	return msg
}

// GasPriceMultiplier returns the gas price multiplier for the given chain
func GasPriceMultiplier(chain chains.Chain) float64 {
	switch chain.Consensus {
	case chains.Consensus_ethereum:
		return clientcommon.EVMOutboundGasPriceMultiplier
	case chains.Consensus_bitcoin:
		return clientcommon.BTCOutboundGasPriceMultiplier
	default:
		return clientcommon.DefaultGasPriceMultiplier
	}
}

// WrapMessageWithAuthz wraps a message with an authz message
// used since a hotkey is used to broadcast the transactions, instead of the operator
func WrapMessageWithAuthz(msg sdk.Msg) (sdk.Msg, clientauthz.Signer, error) {
	msgURL := sdk.MsgTypeURL(msg)

	// verify message validity
	if m, ok := msg.(sdk.HasValidateBasic); ok {
		if err := m.ValidateBasic(); err != nil {
			return nil, clientauthz.Signer{}, errors.Wrapf(err, "invalid message %q", msgURL)
		}
	}

	authzSigner := clientauthz.GetSigner(msgURL)
	authzMessage := authz.NewMsgExec(authzSigner.GranteeAddress, []sdk.Msg{msg})

	return &authzMessage, authzSigner, nil
}

// PostOutboundTracker adds an outbound tracker
func (c *Client) PostOutboundTracker(ctx context.Context, chainID int64, nonce uint64, txHash string) (string, error) {
	// don't report if the tracker already contains the txHash
	tracker, err := c.GetOutboundTracker(ctx, chains.Chain{
		ChainId: chainID,
	}, nonce)
	if err == nil {
		for _, hash := range tracker.HashList {
			if strings.EqualFold(hash.TxHash, txHash) {
				return "", nil
			}
		}
	}

	signerAddress := c.keys.GetOperatorAddress().String()
	msg := types.NewMsgAddOutboundTracker(signerAddress, chainID, nonce, txHash)

	authzMsg, authzSigner, err := WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	museTxHash, err := c.Broadcast(ctx, AddOutboundTrackerGasLimit, authzMsg, authzSigner)
	if err != nil {
		return "", err
	}

	return museTxHash, nil
}
