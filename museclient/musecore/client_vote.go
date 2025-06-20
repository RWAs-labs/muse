package musecore

import (
	"context"

	"github.com/RWAs-labs/go-tss/blame"
	"github.com/pkg/errors"

	zctx "github.com/RWAs-labs/muse/museclient/context"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/retry"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	observerclient "github.com/RWAs-labs/muse/x/observer/client/cli"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// PostVoteGasPrice posts a gas price vote. Returns txHash and error.
func (c *Client) PostVoteGasPrice(
	ctx context.Context,
	chain chains.Chain,
	gasPrice uint64, priorityFee, blockNum uint64,
) (string, error) {
	// get gas price multiplier for the chain
	multiplier := GasPriceMultiplier(chain)

	// #nosec G115 always in range
	gasPrice = uint64(float64(gasPrice) * multiplier)
	signerAddress := c.keys.GetOperatorAddress().String()
	msg := types.NewMsgVoteGasPrice(signerAddress, chain.ChainId, gasPrice, priorityFee, blockNum)

	authzMsg, authzSigner, err := WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	hash, err := retry.DoTypedWithRetry(func() (string, error) {
		return c.Broadcast(ctx, PostGasPriceGasLimit, authzMsg, authzSigner)
	})

	if err != nil {
		return "", errors.Wrap(err, "unable to broadcast vote gas price")
	}

	return hash, nil
}

// PostVoteTSS sends message to vote TSS. Returns txHash and error.
func (c *Client) PostVoteTSS(
	ctx context.Context,
	tssPubKey string,
	keyGenMuseHeight int64,
	status chains.ReceiveStatus,
) (string, error) {
	signerAddress := c.keys.GetOperatorAddress().String()
	msg := observertypes.NewMsgVoteTSS(signerAddress, tssPubKey, keyGenMuseHeight, status)

	authzMsg, authzSigner, err := WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	museTxHash, err := retry.DoTypedWithRetry(func() (string, error) {
		return c.Broadcast(ctx, PostTSSGasLimit, authzMsg, authzSigner)
	})

	if err != nil {
		return "", errors.Wrap(err, "unable to broadcast vote for setting tss")
	}

	return museTxHash, nil
}

// PostVoteBlameData posts blame data message to musecore. Returns txHash and error.
func (c *Client) PostVoteBlameData(
	ctx context.Context,
	blame *blame.Blame,
	chainID int64,
	index string,
) (string, error) {
	signerAddress := c.keys.GetOperatorAddress().String()
	museBlame := observertypes.Blame{
		Index:         index,
		FailureReason: blame.FailReason,
		Nodes:         observerclient.ConvertNodes(blame.BlameNodes),
	}
	msg := observertypes.NewMsgVoteBlameMsg(signerAddress, chainID, museBlame)

	authzMsg, authzSigner, err := WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	var gasLimit uint64 = PostBlameDataGasLimit

	museTxHash, err := retry.DoTypedWithRetry(func() (string, error) {
		return c.Broadcast(ctx, gasLimit, authzMsg, authzSigner)
	})

	if err != nil {
		return "", errors.Wrap(err, "unable to broadcast blame data")
	}

	return museTxHash, nil
}

// PostVoteOutbound posts a vote on an observed outbound tx from a MsgVoteOutbound.
// Returns tx hash, ballotIndex, and error.
func (c *Client) PostVoteOutbound(
	ctx context.Context,
	gasLimit, retryGasLimit uint64,
	msg *types.MsgVoteOutbound,
) (string, string, error) {
	authzMsg, authzSigner, err := WrapMessageWithAuthz(msg)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to wrap message with authz")
	}

	// don't post confirmation if it  already voted before
	ballotIndex := msg.Digest()
	hasVoted, err := c.HasVoted(ctx, ballotIndex, msg.Creator)
	if err != nil {
		return "", ballotIndex, errors.Wrapf(
			err,
			"PostVoteOutbound: unable to check if already voted for ballot %s voter %s",
			ballotIndex,
			msg.Creator,
		)
	}
	if hasVoted {
		return "", ballotIndex, nil
	}

	museTxHash, err := retry.DoTypedWithRetry(func() (string, error) {
		return c.Broadcast(ctx, gasLimit, authzMsg, authzSigner)
	})

	if err != nil {
		return "", ballotIndex, errors.Wrap(err, "unable to broadcast vote outbound")
	}

	go func() {
		ctxForWorker := zctx.Copy(ctx, context.Background())

		errMonitor := c.MonitorVoteOutboundResult(ctxForWorker, museTxHash, retryGasLimit, msg)
		if errMonitor != nil {
			c.logger.Error().Err(err).Msg("PostVoteOutbound: failed to monitor vote outbound result")
		}
	}()

	return museTxHash, ballotIndex, nil
}

// PostVoteInbound posts a vote on an observed inbound tx
// retryGasLimit is the gas limit used to resend the tx if it fails because of insufficient gas
// it is used when the ballot is finalized and the inbound tx needs to be processed
func (c *Client) PostVoteInbound(
	ctx context.Context,
	gasLimit, retryGasLimit uint64,
	msg *types.MsgVoteInbound,
) (string, string, error) {
	authzMsg, authzSigner, err := WrapMessageWithAuthz(msg)
	if err != nil {
		return "", "", err
	}

	// don't post send if has already voted before
	ballotIndex := msg.Digest()
	hasVoted, err := c.HasVoted(ctx, ballotIndex, msg.Creator)
	if err != nil {
		return "", ballotIndex, errors.Wrapf(err,
			"PostVoteInbound: unable to check if already voted for ballot %s voter %s",
			ballotIndex,
			msg.Creator,
		)
	}
	if hasVoted {
		return "", ballotIndex, nil
	}

	museTxHash, err := retry.DoTypedWithRetry(func() (string, error) {
		return c.Broadcast(ctx, gasLimit, authzMsg, authzSigner)
	})

	if err != nil {
		return "", ballotIndex, errors.Wrap(err, "unable to broadcast vote inbound")
	}

	go func() {
		ctxForWorker := zctx.Copy(ctx, context.Background())

		errMonitor := c.MonitorVoteInboundResult(ctxForWorker, museTxHash, retryGasLimit, msg)
		if errMonitor != nil {
			c.logger.Error().Err(err).Msg("PostVoteInbound: failed to monitor vote inbound result")
		}
	}()

	return museTxHash, ballotIndex, nil
}
