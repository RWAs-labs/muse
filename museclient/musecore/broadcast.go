package musecore

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"

	"github.com/RWAs-labs/muse/app/ante"
	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/museclient/authz"
	"github.com/RWAs-labs/muse/museclient/logs"
)

// paying 50% more than the current base gas price to buffer for potential block-by-block
// gas price increase due to EIP1559 feemarket on MuseChain
var bufferMultiplier = sdkmath.LegacyMustNewDecFromStr("1.5")

// Broadcast Broadcasts tx to MuseChain. Returns txHash and error
func (c *Client) Broadcast(
	ctx context.Context,
	gasLimit uint64,
	authzWrappedMsg sdktypes.Msg,
	authzSigner authz.Signer,
) (string, error) {
	blockHeight, err := c.GetBlockHeight(ctx)
	if err != nil {
		return "", errors.Wrap(err, "unable to get block height")
	}

	baseGasPrice, err := c.GetBaseGasPrice(ctx)
	if err != nil {
		return "", errors.Wrap(err, "unable to get base gas price")
	}

	// shouldn't happen, but just in case
	if baseGasPrice == 0 {
		baseGasPrice = DefaultBaseGasPrice
	}

	reductionRate := sdkmath.LegacyMustNewDecFromStr(ante.GasPriceReductionRate)

	// multiply gas price by the system tx reduction rate
	adjustedBaseGasPrice := sdkmath.LegacyNewDec(baseGasPrice).Mul(reductionRate).Mul(bufferMultiplier)

	c.mu.Lock()
	defer c.mu.Unlock()

	if blockHeight > c.blockHeight {
		c.blockHeight = blockHeight
		accountNumber, seqNumber, err := c.GetAccountNumberAndSequenceNumber(authzSigner.KeyType)
		if err != nil {
			return "", err
		}

		c.accountNumber[authzSigner.KeyType] = accountNumber

		if c.seqNumber[authzSigner.KeyType] < seqNumber {
			c.seqNumber[authzSigner.KeyType] = seqNumber
		}
	}

	flags := flag.NewFlagSet("museclient", 0)

	factory, err := clienttx.NewFactoryCLI(c.cosmosClientContext, flags)
	if err != nil {
		return "", err
	}

	factory = factory.WithAccountNumber(c.accountNumber[authzSigner.KeyType])
	factory = factory.WithSequence(c.seqNumber[authzSigner.KeyType])
	factory = factory.WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
	builder, err := factory.BuildUnsignedTx(authzWrappedMsg)
	if err != nil {
		return "", errors.Wrap(err, "unable to build unsigned tx")
	}

	builder.SetGasLimit(gasLimit)

	// #nosec G115 always in range
	fee := sdktypes.NewCoins(sdktypes.NewCoin(
		config.BaseDenom,
		sdkmath.NewInt(int64(gasLimit)).Mul(adjustedBaseGasPrice.Ceil().RoundInt()),
	))
	builder.SetFeeAmount(fee)

	err = c.SignTx(factory, c.cosmosClientContext.GetFromName(), builder, true)
	if err != nil {
		return "", errors.Wrap(err, "unable to sign tx")
	}

	txBytes, err := c.cosmosClientContext.TxConfig.TxEncoder()(builder.GetTx())
	if err != nil {
		return "", errors.Wrap(err, "unable to encode tx")
	}

	// broadcast to a Tendermint node
	commit, err := c.cosmosClientContext.BroadcastTxSync(txBytes)
	if err != nil {
		return "", errors.Wrap(err, "fail to broadcast tx sync")
	}

	// Code will be the tendermint ABICode , it start at 1 , so if it is an error , code will not be zero
	if commit.Code > 0 {
		if commit.Code == 32 {
			errMsg := commit.RawLog
			re := regexp.MustCompile(`account sequence mismatch, expected ([0-9]*), got ([0-9]*)`)
			matches := re.FindStringSubmatch(errMsg)
			if len(matches) != 3 {
				return "", err
			}
			expectedSeq, err := strconv.ParseUint(matches[1], 10, 64)
			if err != nil {
				c.logger.Warn().Msgf("cannot parse expected seq %s", matches[1])
				return "", err
			}
			gotSeq, err := strconv.Atoi(matches[2])
			if err != nil {
				c.logger.Warn().Msgf("cannot parse got seq %s", matches[2])
				return "", err
			}
			c.seqNumber[authzSigner.KeyType] = expectedSeq
			c.logger.Warn().
				Msgf("Reset seq number to %d (from err msg) from %d", c.seqNumber[authzSigner.KeyType], gotSeq)
		}
		return commit.TxHash, fmt.Errorf("fail to broadcast to musechain,code:%d, log:%s", commit.Code, commit.RawLog)
	}

	// increment seqNum
	c.seqNumber[authzSigner.KeyType] = c.seqNumber[authzSigner.KeyType] + 1

	return commit.TxHash, nil
}

// SignTx signs a tx with the given name
func (c *Client) SignTx(
	txf clienttx.Factory,
	name string,
	txBuilder client.TxBuilder,
	overwriteSig bool,
) error {
	return clienttx.Sign(context.TODO(), txf, name, txBuilder, overwriteSig)
}

// QueryTxResult query the result of a tx
func (c *Client) QueryTxResult(hash string) (*sdktypes.TxResponse, error) {
	return authtx.QueryTx(c.cosmosClientContext, hash)
}

// HandleBroadcastError returns whether to retry in a few seconds, and whether to report via AddOutboundTracker
// returns (bool retry, bool report)
func HandleBroadcastError(err error, nonce uint64, toChain int64, outboundHash string) (bool, bool) {
	if err == nil {
		return false, false
	}

	msg := err.Error()
	evt := log.Warn().Err(err).
		Str(logs.FieldMethod, "HandleBroadcastError").
		Int64(logs.FieldChain, toChain).
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, outboundHash)

	switch {
	// From the literal meaning of the error message, the tx with this 'nonce' has already been processed,
	// and the latest TSS account nonce has already been incremented.
	// Theoretically, this the tx hash should not be posted to the tracker in this case, but we've already
	// encountered missed outbound tracker caused by unknown reasons (may or may not be false positive).
	//
	// To prevent missed potential outbound tracker, now we pass this hash to tracker reporter in this case.
	// The overhead is:
	// 	- It is uncertain whether this tx hash was the very FIRST accepted tx with THIS 'nonce', it might be the second...
	//  - Once decided to report this tx hash, we need to spawn extra goroutines and making extra RPC queries for monitoring.
	case strings.Contains(msg, "nonce too low"):
		const m = "nonce too low! this might be an unnecessary key-sign. increase retry interval and awaits outbound confirmation"
		evt.Msg(m)
		return false, true

	case strings.Contains(msg, "replacement transaction underpriced"):
		evt.Msg("Broadcast replacement")
		return false, false

	case strings.Contains(msg, "already known"):
		// report to tracker, because there's possibilities a successful broadcast gets this error code
		evt.Msg("Broadcast duplicates")
		return false, true

	default:
		evt.Msg("Broadcast error. Retrying...")
		return true, false
	}
}
