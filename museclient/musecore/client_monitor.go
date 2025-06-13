package musecore

import (
	"context"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/pkg/retry"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// MonitorVoteOutboundResult monitors the result of a vote outbound tx
// retryGasLimit is the gas limit used to resend the tx if it fails because of insufficient gas
// if retryGasLimit is 0, the tx is not resent
func (c *Client) MonitorVoteOutboundResult(
	ctx context.Context,
	museTxHash string,
	retryGasLimit uint64,
	msg *types.MsgVoteOutbound,
) error {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error().
				Interface("panic", r).
				Str("outbound.hash", museTxHash).
				Msg("monitorVoteOutboundResult: recovered from panic")
		}
	}()

	call := func() error {
		return retry.Retry(c.monitorVoteOutboundResult(ctx, museTxHash, retryGasLimit, msg))
	}

	err := retryWithBackoff(call, monitorRetryCount, monitorInterval/2, monitorInterval)
	if err != nil {
		c.logger.Error().Err(err).
			Str("outbound.hash", museTxHash).
			Msg("monitorVoteOutboundResult: unable to query tx result")

		return err
	}

	return nil
}

func (c *Client) monitorVoteOutboundResult(
	ctx context.Context,
	museTxHash string,
	retryGasLimit uint64,
	msg *types.MsgVoteOutbound,
) error {
	// query tx result from MuseChain
	txResult, err := c.QueryTxResult(museTxHash)
	if err != nil {
		return errors.Wrap(err, "failed to query tx result")
	}

	logFields := map[string]any{
		"outbound.hash":    museTxHash,
		"outbound.raw_log": txResult.RawLog,
	}

	switch {
	case strings.Contains(txResult.RawLog, "failed to execute message"):
		// the inbound vote tx shouldn't fail to execute
		// this shouldn't happen
		c.logger.Error().Fields(logFields).Msg("monitorVoteOutboundResult: failed to execute vote")
	case strings.Contains(txResult.RawLog, "out of gas"):
		// if the tx fails with an out of gas error, resend the tx with more gas if retryGasLimit > 0
		c.logger.Debug().Fields(logFields).Msg("monitorVoteOutboundResult: out of gas")

		if retryGasLimit > 0 {
			// new retryGasLimit set to 0 to prevent reentering this function
			if _, _, err := c.PostVoteOutbound(ctx, retryGasLimit, 0, msg); err != nil {
				c.logger.Error().Err(err).Fields(logFields).Msg("monitorVoteOutboundResult: failed to resend tx")
			} else {
				c.logger.Info().Fields(logFields).Msg("monitorVoteOutboundResult: successfully resent tx")
			}
		}
	default:
		c.logger.Debug().Fields(logFields).Msg("monitorVoteOutboundResult: successful")
	}

	return nil
}

// MonitorVoteInboundResult monitors the result of a vote inbound tx
// retryGasLimit is the gas limit used to resend the tx if it fails because of insufficient gas
// if retryGasLimit is 0, the tx is not resent
func (c *Client) MonitorVoteInboundResult(
	ctx context.Context,
	museTxHash string,
	retryGasLimit uint64,
	msg *types.MsgVoteInbound,
) error {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error().
				Interface("panic", r).
				Str("inbound.hash", museTxHash).
				Msg("monitorVoteInboundResult: recovered from panic")
		}
	}()

	call := func() error {
		return retry.Retry(c.monitorVoteInboundResult(ctx, museTxHash, retryGasLimit, msg))
	}

	err := retryWithBackoff(call, monitorRetryCount, monitorInterval/2, monitorInterval)
	if err != nil {
		c.logger.Error().Err(err).
			Str("inbound.hash", museTxHash).
			Msg("monitorVoteInboundResult: unable to query tx result")

		return err
	}

	return nil
}

func (c *Client) monitorVoteInboundResult(
	ctx context.Context,
	museTxHash string,
	retryGasLimit uint64,
	msg *types.MsgVoteInbound,
) error {
	// query tx result from MuseChain
	txResult, err := c.QueryTxResult(museTxHash)
	if err != nil {
		return errors.Wrap(err, "failed to query tx result")
	}

	logFields := map[string]any{
		"inbound.hash":    museTxHash,
		"inbound.raw_log": txResult.RawLog,
	}

	switch {
	case strings.Contains(txResult.RawLog, "failed to execute message"):
		// the inbound vote tx shouldn't fail to execute
		// this shouldn't happen
		c.logger.Error().Fields(logFields).Msg("monitorVoteInboundResult: failed to execute vote")

	case strings.Contains(txResult.RawLog, "out of gas"):
		// if the tx fails with an out of gas error, resend the tx with more gas if retryGasLimit > 0
		c.logger.Debug().Fields(logFields).Msg("monitorVoteInboundResult: out of gas")

		if retryGasLimit > 0 {
			// new retryGasLimit set to 0 to prevent reentering this function
			if resentTxHash, _, err := c.PostVoteInbound(ctx, retryGasLimit, 0, msg); err != nil {
				c.logger.Error().Err(err).Fields(logFields).Msg("monitorVoteInboundResult: failed to resend tx")
			} else {
				logFields["inbound.resent_hash"] = resentTxHash
				c.logger.Info().Fields(logFields).Msgf("monitorVoteInboundResult: successfully resent tx")
			}
		}

	default:
		c.logger.Debug().Fields(logFields).Msgf("monitorVoteInboundResult: successful")
	}

	return nil
}

func retryWithBackoff(call func() error, attempts int, minInternal, maxInterval time.Duration) error {
	if attempts < 1 {
		return errors.New("attempts must be positive")
	}
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = minInternal
	bo.MaxInterval = maxInterval

	backoffWithRetry := backoff.WithMaxRetries(bo, uint64(attempts))

	return retry.DoWithBackoff(call, backoffWithRetry)
}
