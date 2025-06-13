package tss

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	tsscommon "github.com/RWAs-labs/go-tss/common"
	"github.com/RWAs-labs/go-tss/keygen"
	"github.com/RWAs-labs/go-tss/keysign"
	"github.com/RWAs-labs/go-tss/tss"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/sha3"

	"github.com/RWAs-labs/muse/museclient/logs"
	"github.com/RWAs-labs/muse/museclient/metrics"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/ticker"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

const (
	receiveSuccess = chains.ReceiveStatus_success
	receiveFailed  = chains.ReceiveStatus_failed
)

type keygenCeremony struct {
	tss           *tss.Server
	musecore      Musecore
	lastSeenBlock int64
	iterations    int
	logger        zerolog.Logger
}

// KeygenCeremony runs TSS keygen ceremony as a blocking thread.
// Most likely the keygen is already generated, so this function will be a noop.
// Returns the TSS key if generated, or error.
func KeygenCeremony(
	ctx context.Context,
	server *tss.Server,
	zc Musecore,
	logger zerolog.Logger,
) (observertypes.TSS, error) {
	const interval = time.Second

	ceremony := keygenCeremony{
		tss:      server,
		musecore: zc,
		logger:   logger.With().Str(logs.FieldModule, "tss_keygen").Logger(),
	}

	task := func(ctx context.Context, t *ticker.Ticker) error {
		shouldRetry, err := ceremony.iteration(ctx)
		switch {
		case shouldRetry:
			if err != nil && !errors.Is(err, context.Canceled) {
				logger.Error().Err(err).Msg("Keygen error. Retrying...")
			}

			// continue the ticker
			return nil
		case err != nil:
			return errors.Wrap(err, "keygen ceremony failed")
		default:
			// keygen ceremony is complete (or noop)
			t.Stop()
			return nil
		}
	}

	err := ticker.Run(ctx, interval, task, ticker.WithLogger(logger, "tss_keygen"))
	if err != nil {
		return observertypes.TSS{}, err
	}

	// If there was only a single iteration, most likely the TSS is already generated,
	// Otherwise, we need to wait for the next block to ensure TSS is set by internal keepers.
	if ceremony.iterations > 1 {
		if err = ceremony.waitForBlock(ctx); err != nil {
			return observertypes.TSS{}, errors.Wrap(err, "error waiting for the next block")
		}
	}

	return zc.GetTSS(ctx)
}

// iteration runs ceremony iteration every time interval.
// - Get the keygen task from musecore
// - If the keygen is already generated, return (false, nil) => ceremony is complete
// - If the keygen is pending, ensure we're on the right block
// - Iteration also ensured that the logic is invoked ONLY once per block (regardless of the interval)
func (k *keygenCeremony) iteration(ctx context.Context) (shouldRetry bool, err error) {
	k.iterations++

	keygenTask, err := k.musecore.GetKeyGen(ctx)
	switch {
	case err != nil:
		return true, errors.Wrap(err, "unable to get keygen via RPC")
	case keygenTask.Status == observertypes.KeygenStatus_KeyGenSuccess:
		// all good, tss key is already generated
		return false, nil
	case keygenTask.Status == observertypes.KeygenStatus_KeyGenFailed:
		// come back later to try again (musecore will make status=pending)
		return true, nil
	case keygenTask.Status == observertypes.KeygenStatus_PendingKeygen:
		// okay, let's try to generate the TSS key
	default:
		return false, fmt.Errorf("unexpected keygen status %q", keygenTask.Status.String())
	}

	keygenHeight := keygenTask.BlockNumber

	museHeight, err := k.musecore.GetBlockHeight(ctx)
	switch {
	case err != nil:
		return true, errors.Wrap(err, "unable to get muse height")
	case k.blockThrottled(museHeight):
		return true, nil
	case museHeight < keygenHeight:
		k.logger.Info().
			Int64("keygen.height", keygenHeight).
			Int64("muse_height", museHeight).
			Msgf("Waiting for keygen block to arrive or new keygen block to be set")
		return true, nil
	case museHeight > keygenHeight:
		k.logger.Info().
			Int64("keygen.height", keygenHeight).
			Int64("muse_height", museHeight).
			Msgf("Waiting for keygen finalization")
		return true, nil
	}

	// Now we know that the keygen status is PENDING, and we are the KEYGEN block.
	// Let's perform TSS Keygen and then post successful/failed vote to musecore
	newPubKey, err := k.performKeygen(ctx, keygenTask)
	if err != nil {
		k.logger.Error().Err(err).Msg("Keygen failed. Broadcasting failed TSS vote")

		// Vote for failure
		failedVoteHash, err := k.musecore.PostVoteTSS(ctx, "", keygenTask.BlockNumber, receiveFailed)
		if err != nil {
			return false, errors.Wrap(err, "failed to broadcast failed TSS vote")
		}

		k.logger.Info().
			Str("keygen.failed_vote_tx_hash", failedVoteHash).
			Msg("Broadcasted failed TSS keygen vote")

		return true, nil
	}

	successVoteHash, err := k.musecore.PostVoteTSS(ctx, newPubKey, keygenTask.BlockNumber, receiveSuccess)
	if err != nil {
		return false, errors.Wrap(err, "failed to broadcast successful TSS vote")
	}

	k.logger.Info().
		Str("keygen.success_vote_tx_hash", successVoteHash).
		Msg("Broadcasted successful TSS keygen vote")

	k.logger.Info().Msg("Performing TSS key-sign test")

	if err = TestKeySign(k.tss, newPubKey, k.logger); err != nil {
		k.logger.Error().Err(err).Msg("Failed to test TSS keygen")
		// signing can fail even if tss keygen is successful
	}

	return false, nil
}

// performKeygen performs TSS keygen flow via go-tss server. Returns the new TSS public key or error.
// If fails, then it will post blame data to musecore and return an error.
func (k *keygenCeremony) performKeygen(ctx context.Context, keygenTask observertypes.Keygen) (string, error) {
	k.logger.Warn().
		Int64("keygen.block", keygenTask.BlockNumber).
		Strs("keygen.tss_signers", keygenTask.GranteePubkeys).
		Msg("Performing a keygen!")

	req := keygen.NewRequest(keygenTask.GranteePubkeys, keygenTask.BlockNumber, Version, Algo)

	res, err := k.tss.Keygen(req)
	switch {
	case err != nil:
		// returns error on network failure or other non-recoverable errors
		// if the keygen is unsuccessful, the error will be nil
		return "", errors.Wrap(err, "unable to perform keygen")
	case res.Status == tsscommon.Success && res.PubKey != "":
		// desired outcome
		k.logger.Info().
			Interface("keygen.response", res).
			Interface("keygen.tss_public_key", res.PubKey).
			Msg("Keygen successfully generated!")
		return res.PubKey, nil
	}

	// Something went wrong, let's post blame results and then FAIL
	k.logger.Error().
		Str("keygen.blame_round", res.Blame.Round).
		Str("keygen.fail_reason", res.Blame.FailReason).
		Interface("keygen.blame_nodes", res.Blame.BlameNodes).
		Msg("Keygen failed! Sending blame data to musecore")

	// increment blame counter
	for _, node := range res.Blame.BlameNodes {
		metrics.TSSNodeBlamePerPubKey.WithLabelValues(node.Pubkey).Inc()
	}

	blameDigest, err := digestReq(req)
	if err != nil {
		return "", errors.Wrap(err, "unable to create digest")
	}

	blameIndex := fmt.Sprintf("keygen-%s-%d", blameDigest, keygenTask.BlockNumber)
	chainID := k.musecore.Chain().ChainId

	museHash, err := k.musecore.PostVoteBlameData(ctx, &res.Blame, chainID, blameIndex)
	if err != nil {
		return "", errors.Wrap(err, "unable to post blame data to musecore")
	}

	k.logger.Info().Str("keygen.blame_tx_hash", museHash).Msg("Posted blame data to musecore")

	return "", errors.Errorf("keygen failed: %s", res.Blame.FailReason)
}

// returns true if the block is throttled i.e. we should wait for the next block.
func (k *keygenCeremony) blockThrottled(currentBlock int64) bool {
	switch {
	case currentBlock == 0:
		return false
	case k.lastSeenBlock == currentBlock:
		return true
	default:
		k.lastSeenBlock = currentBlock
		return false
	}
}

func (k *keygenCeremony) waitForBlock(ctx context.Context) error {
	height, err := k.musecore.GetBlockHeight(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get block height (initial)")
	}

	for {
		k.logger.Info().Msg("Waiting for the next block to arrive")
		newHeight, err := k.musecore.GetBlockHeight(ctx)
		switch {
		case err != nil:
			return errors.Wrap(err, "unable to get block height")
		case newHeight > height:
			return nil
		default:
			time.Sleep(time.Second)
		}
	}
}

func digestReq(req keygen.Request) (string, error) {
	bytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(bytes)
	digest := hex.EncodeToString(hasher.Sum(nil))

	return digest, nil
}

var testKeySignData = []byte("hello meta")

// TestKeySign performs a TSS key-sign test of sample data.
func TestKeySign(keySigner KeySigner, tssPubKeyBec32 string, logger zerolog.Logger) error {
	logger = logger.With().Str(logs.FieldModule, "tss_keysign").Logger()

	tssPubKey, err := NewPubKeyFromBech32(tssPubKeyBec32)
	if err != nil {
		return errors.Wrap(err, "unable to parse TSS public key")
	}

	hashedData := crypto.Keccak256Hash(testKeySignData)

	logger.Info().
		Str("keysign.test_data", string(testKeySignData)).
		Str("keysign.test_data_hash", hashedData.String()).
		Msg("Performing TSS key-sign test")

	req := keysign.NewRequest(
		tssPubKey.Bech32String(),
		[]string{base64.StdEncoding.EncodeToString(hashedData.Bytes())},
		10,
		nil,
		Version,
	)

	res, err := keySigner.KeySign(req)
	switch {
	case err != nil:
		return errors.Wrap(err, "key signing request error")
	case res.Status != tsscommon.Success:
		logger.Error().Interface("keysign.fail_blame", res.Blame).Msg("Keysign failed")
		return errors.Wrapf(err, "key signing is not successful (status %d)", res.Status)
	case len(res.Signatures) == 0:
		return errors.New("signatures list is empty")
	}

	// 32B msg hash, 32B R, 32B S, 1B RC
	signature := res.Signatures[0]

	logger.Info().Interface("keysign.signature", signature).Msg("Received signature from TSS")

	if _, err = VerifySignature(signature, tssPubKey, hashedData.Bytes()); err != nil {
		return errors.Wrap(err, "signature verification failed")
	}

	logger.Info().Msg("TSS key-sign test passed")

	return nil
}
