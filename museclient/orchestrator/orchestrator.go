// Package orchestrator is responsible for (de)provisioning, running, and monitoring various observer-signer instances.
// It also updates app context with data from musecore (eg chain parameters).
package orchestrator

import (
	"context"
	"math"
	"sync"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/RWAs-labs/muse/museclient/chains/base"
	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	zctx "github.com/RWAs-labs/muse/museclient/context"
	"github.com/RWAs-labs/muse/museclient/logs"
	"github.com/RWAs-labs/muse/museclient/metrics"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/constant"
	"github.com/RWAs-labs/muse/pkg/scheduler"
	"github.com/RWAs-labs/muse/pkg/ticker"
)

// Orchestrator chain orchestrator.
type Orchestrator struct {
	deps      *Dependencies
	scheduler *scheduler.Scheduler

	chains map[int64]ObserverSigner
	mu     sync.RWMutex

	operatorBalance sdkmath.Int

	logger loggers
}

type loggers struct {
	zerolog.Logger
	sampled zerolog.Logger
	base    base.Logger
}

const schedulerGroup = scheduler.Group("orchestrator")

type ObserverSigner interface {
	Chain() chains.Chain
	Start(ctx context.Context) error
	Stop()
}

type Dependencies struct {
	Musecore  interfaces.MusecoreClient
	TSS       interfaces.TSSSigner
	DBPath    string
	Telemetry *metrics.TelemetryServer
}

func New(scheduler *scheduler.Scheduler, deps *Dependencies, logger base.Logger) (*Orchestrator, error) {
	if err := validateConstructor(scheduler, deps); err != nil {
		return nil, errors.Wrap(err, "invalid args")
	}

	return &Orchestrator{
		scheduler: scheduler,
		deps:      deps,
		chains:    make(map[int64]ObserverSigner),
		logger:    newLoggers(logger),
	}, nil
}

func (oc *Orchestrator) Start(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	newBlocksChan, err := oc.deps.Musecore.NewBlockSubscriber(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to subscribe to new block")
	}

	// syntax sugar
	opts := func(name string, opts ...scheduler.Opt) []scheduler.Opt {
		return append(opts, scheduler.GroupName(schedulerGroup), scheduler.Name(name))
	}

	contextInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(app.Config().ConfigUpdateTicker)
	})

	// every other block, regardless of block events from musecore
	syncInterval := scheduler.Interval(2 * constant.MuseBlockTime)

	blocksTicker := scheduler.BlockTicker(newBlocksChan)

	// refresh preflight metrics in a lazy manner
	preflightTicker := scheduler.Interval(1 * time.Minute)

	oc.scheduler.Register(ctx, oc.UpdateContext, opts("update_context", contextInterval)...)
	oc.scheduler.Register(ctx, oc.SyncChains, opts("sync_chains", syncInterval)...)
	oc.scheduler.Register(ctx, oc.updateMetrics, opts("update_metrics", blocksTicker)...)
	oc.scheduler.Register(ctx, oc.reportPreflightMetrics, opts("report_preflight_metrics", preflightTicker)...)

	return nil
}

func (oc *Orchestrator) Stop() {
	oc.logger.Info().Msg("Stopping orchestrator")

	// stops *all* scheduler tasks
	oc.scheduler.Stop()
}

func (oc *Orchestrator) UpdateContext(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	err = UpdateAppContext(ctx, app, oc.deps.Musecore, oc.logger.Logger)

	switch {
	case errors.Is(err, ErrUpgradeRequired):
		const msg = "Upgrade detected. Kill the process, " +
			"replace the binary with upgraded version, and restart museclientd"

		oc.logger.Warn().Str("upgrade", err.Error()).Msg(msg)

		// stop the orchestrator
		go oc.Stop()

		return nil
	case err != nil:
		return errors.Wrap(err, "unable to update app context")
	default:
		return nil
	}
}

var errSkipChain = errors.New("skip chain")

func (oc *Orchestrator) SyncChains(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	var (
		added, removed  int
		presentChainIDs = make([]int64, 0)
	)

	for _, chain := range app.ListChains() {
		// skip musechain
		if chain.IsMuse() {
			continue
		}

		presentChainIDs = append(presentChainIDs, chain.ID())

		// skip existing chain
		if oc.hasChain(chain.ID()) {
			continue
		}

		var observerSigner ObserverSigner

		switch {
		case chain.IsBitcoin():
			observerSigner, err = oc.bootstrapBitcoin(ctx, chain)
		case chain.IsEVM():
			observerSigner, err = oc.bootstrapEVM(ctx, chain)
		case chain.IsSolana():
			observerSigner, err = oc.bootstrapSolana(ctx, chain)
		case chain.IsSui():
			observerSigner, err = oc.bootstrapSui(ctx, chain)
		case chain.IsTON():
			observerSigner, err = oc.bootstrapTON(ctx, chain)
		}

		switch {
		case errors.Is(err, errSkipChain):
			// TODO use throttled logger instead of sampled one.
			// https://github.com/RWAs-labs/muse/issues/3336
			oc.logger.sampled.Warn().Err(err).Fields(chain.LogFields()).Msg("Skipping observer-signer")
			continue
		case err != nil:
			oc.logger.Error().Err(err).Fields(chain.LogFields()).Msg("Failed to bootstrap observer-signer")
			continue
		case observerSigner == nil:
			// should not happen
			oc.logger.Error().Fields(chain.LogFields()).Msg("Nil observer-signer")
			continue
		}

		if err = observerSigner.Start(ctx); err != nil {
			oc.logger.Error().Err(err).Fields(chain.LogFields()).Msg("Failed to start observer-signer")
			continue
		}

		oc.addChain(observerSigner)
		added++
	}

	removed = oc.removeMissingChains(presentChainIDs)

	if (added + removed) > 0 {
		oc.logger.Info().
			Int("chains.added", added).
			Int("chains.removed", removed).
			Msg("Synced observer-signers")
	}

	return nil
}

var (
	zero   = sdkmath.NewInt(0)
	maxInt = sdkmath.NewInt(math.MaxInt64)
)

func (oc *Orchestrator) updateMetrics(ctx context.Context) error {
	block, sleepDuration, err := scheduler.BlockFromContextWithDelay(ctx)
	if err != nil {
		return errors.Wrap(err, "unable get block from context")
	}

	musecore := oc.deps.Musecore
	ts := oc.deps.Telemetry

	museBlockHeight := block.Block.Height

	// 0. Set block metrics
	metrics.CoreBlockLatency.Set(time.Since(block.Block.Time).Seconds())
	metrics.CoreBlockLatencySleep.Set(sleepDuration.Seconds())

	ts.SetCoreBlockNumber(museBlockHeight)

	// 1. Fetch hot key balance
	balance, err := musecore.GetMuseHotKeyBalance(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get hot key balance")
	}

	// 2. Set it within orchestrator
	oc.operatorBalance = balance

	// 3. Update telemetry
	diff := oc.operatorBalance.Sub(balance)
	if diff.GT(zero) && diff.LT(maxInt) {
		ts.AddFeeEntry(museBlockHeight, diff.Int64())
	}

	// 4. Update metrics
	burnRate := ts.HotKeyBurnRate.GetBurnRate().Int64()
	metrics.HotKeyBurnRate.Set(float64(burnRate))

	return nil
}

func (oc *Orchestrator) reportPreflightMetrics(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	return ReportPreflightMetrics(ctx, app, oc.deps.Musecore, oc.logger.Logger)
}

func (oc *Orchestrator) hasChain(chainID int64) bool {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	_, ok := oc.chains[chainID]
	return ok
}

func (oc *Orchestrator) chainIDs() []int64 {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	ids := make([]int64, 0, len(oc.chains))
	for chainID := range oc.chains {
		ids = append(ids, chainID)
	}

	return ids
}

func (oc *Orchestrator) addChain(observerSigner ObserverSigner) {
	chain := observerSigner.Chain()

	oc.mu.Lock()
	defer oc.mu.Unlock()

	// noop
	if _, ok := oc.chains[chain.ChainId]; ok {
		return
	}

	oc.chains[chain.ChainId] = observerSigner
	oc.logger.Info().Fields(chain.LogFields()).Msg("Added observer-signer")
}

func (oc *Orchestrator) removeChain(chainID int64) {
	// noop, should not happen
	if !oc.hasChain(chainID) {
		return
	}

	// blocking call
	oc.chains[chainID].Stop()

	oc.mu.Lock()
	delete(oc.chains, chainID)
	oc.mu.Unlock()

	oc.logger.Info().Int64(logs.FieldChain, chainID).Msg("Removed observer-signer")
}

// removeMissingChains stops and deletes chains
// that are not present in the list of chainIDs (e.g. after governance proposal)
func (oc *Orchestrator) removeMissingChains(presentChainIDs []int64) int {
	presentChainsSet := make(map[int64]struct{})
	for _, chainID := range presentChainIDs {
		presentChainsSet[chainID] = struct{}{}
	}

	existingIDs := oc.chainIDs()
	removed := 0

	for _, chainID := range existingIDs {
		if _, ok := presentChainsSet[chainID]; ok {
			// all good, chain is present
			continue
		}

		oc.removeChain(chainID)
		removed++
	}

	return removed
}

func validateConstructor(s *scheduler.Scheduler, dep *Dependencies) error {
	switch {
	case s == nil:
		return errors.New("scheduler is nil")
	case dep == nil:
		return errors.New("dependencies are nil")
	case dep.Musecore == nil:
		return errors.New("musecore is nil")
	case dep.TSS == nil:
		return errors.New("tss is nil")
	case dep.Telemetry == nil:
		return errors.New("telemetry is nil")
	case dep.DBPath == "":
		return errors.New("db path is empty")
	}

	return nil
}

func newLoggers(baseLogger base.Logger) loggers {
	std := baseLogger.Std.With().Str(logs.FieldModule, "orchestrator").Logger()

	return loggers{
		Logger:  std,
		sampled: std.Sample(&zerolog.BasicSampler{N: 10}),
		base:    baseLogger,
	}
}
