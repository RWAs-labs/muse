package maintenance

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/errors"
	"github.com/rs/zerolog"
	"golang.org/x/mod/semver"

	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	"github.com/RWAs-labs/muse/pkg/bg"
	"github.com/RWAs-labs/muse/pkg/constant"
	"github.com/RWAs-labs/muse/pkg/retry"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

const restartListenerTicker = 10 * time.Second

// ShutdownListener is a struct that listens for scheduled shutdown notices via the observer
// operational flags
type ShutdownListener struct {
	client interfaces.MusecoreClient
	logger zerolog.Logger

	lastRestartHeightMissed int64
	// get the current version of museclient
	getVersion func() string
}

// NewShutdownListener creates a new ShutdownListener.
func NewShutdownListener(client interfaces.MusecoreClient, logger zerolog.Logger) *ShutdownListener {
	log := logger.With().Str("module", "shutdown_listener").Logger()
	return &ShutdownListener{
		client:     client,
		logger:     log,
		getVersion: getVersionDefault,
	}
}

// RunPreStartCheck runs any checks that must run before fully starting museclient.
// Specifically this should be run before any TSS P2P is started.
func (o *ShutdownListener) RunPreStartCheck(ctx context.Context) error {
	operationalFlags, err := o.getOperationalFlagsWithRetry(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get initial operational flags")
	}
	return o.checkMinimumVersion(operationalFlags)
}

func (o *ShutdownListener) Listen(ctx context.Context, action func()) {
	var (
		withLogger = bg.WithLogger(o.logger)
		onComplete = bg.OnComplete(action)
	)

	bg.Work(ctx, o.waitForUpdate, bg.WithName("shutdown_listener.wait_for_update"), withLogger, onComplete)
}

func (o *ShutdownListener) waitForUpdate(ctx context.Context) error {
	operationalFlags, err := o.getOperationalFlagsWithRetry(ctx)
	if err != nil {
		return errors.Wrap(err, "get initial operational flags")
	}
	if o.handleNewFlags(ctx, operationalFlags) {
		return nil
	}

	ticker := time.NewTicker(restartListenerTicker)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			operationalFlags, err = o.client.GetOperationalFlags(ctx)
			if err != nil {
				return errors.Wrap(err, "unable to get operational flags")
			}
			if o.handleNewFlags(ctx, operationalFlags) {
				return nil
			}
		case <-ctx.Done():
			o.logger.Info().Msg("waitForUpdate (shutdown listener) stopped")
			return nil
		}
	}
}

func (o *ShutdownListener) getOperationalFlagsWithRetry(ctx context.Context) (observertypes.OperationalFlags, error) {
	return retry.DoTypedWithBackoffAndRetry(
		func() (observertypes.OperationalFlags, error) { return o.client.GetOperationalFlags(ctx) },
		retry.DefaultConstantBackoff(),
	)
}

// handleNewFlags processes the flags and returns true if a shutdown should be signaled
func (o *ShutdownListener) handleNewFlags(ctx context.Context, f observertypes.OperationalFlags) bool {
	if err := o.checkMinimumVersion(f); err != nil {
		o.logger.Error().Err(err).Any("operational_flags", f).Msg("minimum version check")
		return true
	}
	if f.RestartHeight < 1 {
		return false
	}

	currentHeight, err := o.client.GetBlockHeight(ctx)
	if err != nil {
		o.logger.Error().Err(err).Msg("unable to get block height")
		return false
	}

	if f.RestartHeight < currentHeight {
		// only log restart height misseed once
		if o.lastRestartHeightMissed != f.RestartHeight {
			o.logger.Error().
				Int64("restart_height", f.RestartHeight).
				Int64("current_height", currentHeight).
				Msg("restart height missed")
		}
		o.lastRestartHeightMissed = f.RestartHeight
		return false
	}

	o.logger.Warn().
		Int64("restart_height", f.RestartHeight).
		Int64("current_height", currentHeight).
		Msg("restart scheduled")

	newBlockChan, err := o.client.NewBlockSubscriber(ctx)
	if err != nil {
		o.logger.Error().Err(err).Msg("unable to subscribe to new blocks")
		return false
	}
	for {
		select {
		case newBlock := <-newBlockChan:
			if newBlock.Block.Height >= f.RestartHeight {
				o.logger.Warn().
					Int64("restart_height", f.RestartHeight).
					Int64("current_height", newBlock.Block.Height).
					Msg("restart height reached")
				return true
			}
		case <-ctx.Done():
			return false
		}
	}
}

func (o *ShutdownListener) checkMinimumVersion(f observertypes.OperationalFlags) error {
	if f.MinimumVersion != "" {
		currentVersion := o.getVersion()
		if semver.Compare(currentVersion, f.MinimumVersion) == -1 {
			return fmt.Errorf(
				"current version (%s) is less than minimum version (%s)",
				currentVersion,
				f.MinimumVersion,
			)
		}
	}
	return nil
}

func getVersionDefault() string {
	return constant.GetNormalizedVersion()
}
