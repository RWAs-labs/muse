package cctx

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"

	musetoolchains "github.com/RWAs-labs/muse/cmd/musetool/chains"
	"github.com/RWAs-labs/muse/cmd/musetool/context"
	"github.com/RWAs-labs/muse/museclient/chains/bitcoin/client"
	museevmclient "github.com/RWAs-labs/muse/museclient/chains/evm/client"
	solanarpc "github.com/RWAs-labs/muse/museclient/chains/solana/rpc"
	museclientConfig "github.com/RWAs-labs/muse/museclient/config"
	"github.com/RWAs-labs/muse/pkg/chains"
)

func (c *TrackingDetails) CheckOutbound(ctx *context.Context) error {
	var (
		outboundChain = ctx.GetInboundChain()
	)

	// We do not need to handle the case for muse chain as the outbound is confirmed in the same block.
	switch {
	case outboundChain.IsEVMChain():
		return c.checkEvmOutboundTx(ctx)
	case outboundChain.IsBitcoinChain():
		return c.checkBitcoinOutboundTx(ctx)
	case outboundChain.IsSolanaChain():
		return c.checkSolanaOutboundTx(ctx)
	default:
		return fmt.Errorf("unsupported outbound chain")
	}
}

// checkEvmOutboundTx checks if the outbound transaction is confirmed on the outbound chain.
// If it's confirmed, we update the status to PendingOutboundVoting or PendingRevertVoting. Which means that the confirmation is done and we are not waiting for observers to vote
// Transition Status PendingConfirmation -> Status PendingVoting
func (c *TrackingDetails) checkEvmOutboundTx(ctx *context.Context) error {
	var (
		txHashList     = c.OutboundTrackerHashList
		outboundChain  = c.OutboundChain
		musecoreClient = ctx.GetMuseCoreClient()
		goCtx          = ctx.GetContext()
	)

	chainParams, err := musecoreClient.GetChainParamsForChainID(goCtx, outboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get chain params: %w", err)
	}

	// create evm client for the observation chain
	evmClient, err := musetoolchains.GetEvmClient(ctx, outboundChain)
	if err != nil {
		return fmt.Errorf("failed to create evm client: %w", err)
	}

	foundConfirmedTx := false

	// If one of the hash is confirmed, we update the status to pending voting
	// There might be a condition where we have multiple txs and the wrong tx is confirmed.
	// To verify that we need, check CCTX data
	for _, hash := range txHashList {
		tx, _, err := musetoolchains.GetEvmTx(ctx, evmClient, hash, outboundChain)
		if err != nil {
			continue
		}
		// Signer is unused
		c := museevmclient.New(evmClient, ethtypes.NewLondonSigner(tx.ChainId()))
		confirmed, err := c.IsTxConfirmed(goCtx, hash, chainParams.OutboundConfirmationSafe())
		if err != nil {
			continue
		}
		if confirmed {
			foundConfirmedTx = true
			break
		}
	}
	if foundConfirmedTx {
		c.updateOutboundVoting()
	}
	return nil
}

func (c *TrackingDetails) checkSolanaOutboundTx(ctx *context.Context) error {
	var (
		txHashList = c.OutboundTrackerHashList
		goCtx      = ctx.GetContext()
		cfg        = ctx.GetConfig()
	)

	foundConfirmedTx := false
	solClient := solrpc.New(cfg.SolanaRPC)
	if solClient == nil {
		return fmt.Errorf("error creating rpc client")
	}
	for _, hash := range txHashList {
		signature := solana.MustSignatureFromBase58(hash)
		_, err := solanarpc.GetTransaction(goCtx, solClient, signature)
		if err != nil {
			continue
		}
		foundConfirmedTx = true
	}

	if foundConfirmedTx {
		c.updateOutboundVoting()
	}
	return nil
}

func (c *TrackingDetails) checkBitcoinOutboundTx(ctx *context.Context) error {
	var (
		txHashList     = c.OutboundTrackerHashList
		outboundChain  = c.OutboundChain
		musecoreClient = ctx.GetMuseCoreClient()
		goCtx          = ctx.GetContext()
		cfg            = ctx.GetConfig()
		logger         = ctx.GetLogger()
	)

	chainParams, err := musecoreClient.GetChainParamsForChainID(goCtx, outboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get chain params: %w", err)
	}
	confirmationCount := chainParams.OutboundConfirmationSafe()

	params, err := chains.BitcoinNetParamsFromChainID(outboundChain.ChainId)
	if err != nil {
		return fmt.Errorf("unable to get bitcoin net params from chain id: %w", err)
	}

	connCfg := museclientConfig.BTCConfig{
		RPCUsername: cfg.BtcUser,
		RPCPassword: cfg.BtcPassword,
		RPCHost:     cfg.BtcHost,
		RPCParams:   params.Name,
	}

	btcClient, err := client.New(connCfg, outboundChain.ChainId, logger)
	if err != nil {
		return fmt.Errorf("unable to create rpc client: %w", err)
	}

	err = btcClient.Ping(goCtx)
	if err != nil {
		return fmt.Errorf("error ping the bitcoin server: %w", err)
	}

	foundConfirmedTx := false

	for _, hash := range txHashList {
		txHash, err := chainhash.NewHashFromStr(hash)
		if err != nil {
			continue
		}
		tx, err := btcClient.GetRawTransactionVerbose(goCtx, txHash)
		if err != nil {
			continue
		}

		if tx.Confirmations >= confirmationCount {
			foundConfirmedTx = true
		}
	}

	if foundConfirmedTx {
		c.updateOutboundVoting()
	}
	return nil
}
