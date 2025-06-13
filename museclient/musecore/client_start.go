package musecore

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/RWAs-labs/muse/museclient/authz"
	"github.com/RWAs-labs/muse/museclient/config"
	"github.com/RWAs-labs/muse/museclient/keys"
	museauthz "github.com/RWAs-labs/muse/pkg/authz"
	"github.com/RWAs-labs/muse/pkg/ticker"
)

// This file contains some high level logic for creating a musecore client
// when starting museclientd in cmd/museclientd/start.go

// NewFromConfig creates a new client from the given config.
// It also makes sure that the musecore (i.e. the node) is ready to be used.
func NewFromConfig(
	ctx context.Context,
	cfg *config.Config,
	hotkeyPassword string,
	logger zerolog.Logger,
) (*Client, error) {
	hotKey := cfg.AuthzHotkey

	chainIP := cfg.MuseCoreURL

	kb, _, err := keys.GetKeyringKeybase(*cfg, hotkeyPassword)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get keyring base")
	}

	granterAddress, err := sdk.AccAddressFromBech32(cfg.AuthzGranter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get granter address")
	}

	k := keys.NewKeysWithKeybase(kb, granterAddress, cfg.AuthzHotkey, hotkeyPassword)

	// All votes broadcasts to musecore are wrapped in authz.
	// This is to ensure that the user does not need to keep their operator key online,
	// and can use a cold key to sign votes
	signerAddress, err := k.GetAddress()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get signer address")
	}

	authz.SetupAuthZSignerList(k.GetOperatorAddress().String(), signerAddress)

	// Create client
	client, err := NewClient(k, chainIP, hotKey, cfg.ChainID, logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the client")
	}

	// Make sure that the node produces blocks
	if err = ensureBlocksProduction(ctx, client); err != nil {
		return nil, errors.Wrap(err, "musecore unavailable")
	}

	// Prepare the client
	if err = prepareMusecoreClient(ctx, client, cfg); err != nil {
		return nil, errors.Wrap(err, "failed to prepare the client")
	}

	return client, nil
}

// ensureBlocksProduction waits for musecore to be ready (i.e. producing blocks)
func ensureBlocksProduction(ctx context.Context, zc *Client) error {
	const (
		interval = 5 * time.Second
		attempts = 15
	)

	var (
		retryCount = 0
		start      = time.Now()
	)

	task := func(ctx context.Context, t *ticker.Ticker) error {
		blockHeight, err := zc.GetBlockHeight(ctx)

		if err == nil && blockHeight > 1 {
			zc.logger.Info().Msgf("Musecore is ready, block height: %d", blockHeight)
			t.Stop()
			return nil
		}

		retryCount++
		if retryCount > attempts {
			return fmt.Errorf("musecore is not ready, timeout %s", time.Since(start).String())
		}

		zc.logger.Info().Msgf("Failed to get block number, retry: %d/%d", retryCount, attempts)
		return nil
	}

	return ticker.Run(ctx, interval, task)
}

// prepareMusecoreClient prepares the musecore client for use.
func prepareMusecoreClient(ctx context.Context, zc *Client, cfg *config.Config) error {
	// Set grantee account number and sequence number
	if err := zc.SetAccountNumber(museauthz.MuseClientGranteeKey); err != nil {
		return errors.Wrap(err, "failed to set account number")
	}

	res, err := zc.GetNodeInfo(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get node info")
	}

	network := res.GetDefaultNodeInfo().Network
	if network != cfg.ChainID {
		zc.logger.Warn().
			Str("got", cfg.ChainID).
			Str("want", network).
			Msg("Musecore chain id config mismatch. Forcing update from the network")

		cfg.ChainID = network
		if err = zc.UpdateChainID(cfg.ChainID); err != nil {
			return errors.Wrap(err, "failed to update chain id")
		}
	}

	return nil
}
