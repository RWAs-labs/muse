package local

import (
	"context"
	"fmt"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	cometbfttypes "github.com/cometbft/cometbft/types"

	"github.com/RWAs-labs/muse/e2e/config"
)

// monitorBlockProductionCancel calls monitorBlockProduction and exits upon any error
func monitorBlockProductionCancel(ctx context.Context, cancel context.CancelCauseFunc, conf config.Config) {
	err := monitorBlockProduction(ctx, conf)
	if err != nil {
		cancel(err)
	}
}

// monitorBlockProduction subscribes to new block events to monitor if blocks are being produced
// at least every four seconds
func monitorBlockProduction(ctx context.Context, conf config.Config) error {
	rpcClient, err := rpchttp.New(conf.RPCs.MuseCoreRPC, "/websocket")
	if err != nil {
		return fmt.Errorf("new musecore rpc: %w", err)
	}

	err = rpcClient.WSEvents.Start()
	if err != nil {
		return fmt.Errorf("start ws events: %w", err)
	}
	blockEventChan, err := rpcClient.WSEvents.Subscribe(ctx, "", "tm.event='NewBlock'")
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}
	latestNewBlockEvent := cometbfttypes.EventDataNewBlock{}
	for {
		select {
		case event := <-blockEventChan:
			newBlockEvent, ok := event.Data.(cometbfttypes.EventDataNewBlock)
			if !ok {
				return fmt.Errorf("expecting new block event, got %T", event.Data)
			}
			latestNewBlockEvent = newBlockEvent
		case <-time.After(10 * time.Second):
			return fmt.Errorf("timed out waiting for new block (last block %d)", latestNewBlockEvent.Block.Height)
		case <-ctx.Done():
			return nil
		}
	}
}
