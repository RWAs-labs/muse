package main

import (
	"context"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/museclient/config"
	zctx "github.com/RWAs-labs/muse/museclient/context"
	"github.com/RWAs-labs/muse/museclient/keys"
	"github.com/RWAs-labs/muse/museclient/musecore"
)

// isObserverNode checks whether THIS node is an observer node.
func isObserverNode(ctx context.Context, zc *musecore.Client) (bool, error) {
	observers, err := zc.GetObserverList(ctx)
	if err != nil {
		return false, errors.Wrap(err, "unable to get observers list")
	}

	operatorAddress := zc.GetKeys().GetOperatorAddress().String()

	for _, observer := range observers {
		if observer == operatorAddress {
			return true, nil
		}
	}

	return false, nil
}

func isEnvFlagEnabled(flag string) bool {
	v, _ := strconv.ParseBool(os.Getenv(flag))
	return v
}

func btcChainIDsFromContext(app *zctx.AppContext) []int64 {
	var (
		btcChains   = app.FilterChains(zctx.Chain.IsBitcoin)
		btcChainIDs = make([]int64, len(btcChains))
	)

	for i, chain := range btcChains {
		btcChainIDs[i] = chain.ID()
	}

	return btcChainIDs
}

func resolveObserverPubKeyBech32(cfg config.Config, hotKeyPassword string) (string, error) {
	// Get observer's public key ("grantee pub key")
	_, granteePubKeyBech32, err := keys.GetKeyringKeybase(cfg, hotKeyPassword)
	if err != nil {
		return "", errors.Wrap(err, "unable to get keyring key base")
	}

	return granteePubKeyBech32, nil
}
