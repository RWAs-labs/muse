// Package signer implements the ChainSigner interface for EVM chains
package signer

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"runtime/debug"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/RWAs-labs/muse/museclient/chains/base"
	"github.com/RWAs-labs/muse/museclient/chains/evm/client"
	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	zctx "github.com/RWAs-labs/muse/museclient/context"
	"github.com/RWAs-labs/muse/museclient/logs"
	"github.com/RWAs-labs/muse/museclient/musecore"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/retry"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

const (
	// broadcastBackoff is the initial backoff duration for retrying broadcast
	broadcastBackoff = time.Second * 6

	// broadcastRetries is the maximum number of retries for broadcasting a transaction
	broadcastRetries = 5

	// broadcastTimeout is the timeout for broadcasting a transaction
	// we should allow enough time for the tx submission and avoid fast timeout
	broadcastTimeout = time.Second * 15
)

var (
	// zeroValue is for outbounds that carry no ETH (gas token) value
	zeroValue = big.NewInt(0)
)

// Signer deals with the signing EVM transactions and implements the ChainSigner interface
type Signer struct {
	*base.Signer

	// client is the EVM RPC client to interact with the EVM chain
	client *client.Client

	// museConnectorAddress is the address of the MuseConnector contract
	museConnectorAddress ethcommon.Address

	// er20CustodyAddress is the address of the ERC20Custody contract
	er20CustodyAddress ethcommon.Address

	// gatewayAddress is the address of the Gateway contract
	gatewayAddress ethcommon.Address
}

// New Signer constructor
func New(
	baseSigner *base.Signer,
	client *client.Client,
	museConnectorAddress ethcommon.Address,
	erc20CustodyAddress ethcommon.Address,
	gatewayAddress ethcommon.Address,
) (*Signer, error) {
	return &Signer{
		Signer:               baseSigner,
		client:               client,
		museConnectorAddress: museConnectorAddress,
		er20CustodyAddress:   erc20CustodyAddress,
		gatewayAddress:       gatewayAddress,
	}, nil
}

// SetMuseConnectorAddress sets the muse connector address
func (signer *Signer) SetMuseConnectorAddress(addr ethcommon.Address) {
	// noop
	if (addr == ethcommon.Address{}) || signer.museConnectorAddress == addr {
		return
	}

	signer.Logger().Std.Info().
		Str("signer.old_muse_connector_address", signer.museConnectorAddress.String()).
		Str("signer.new_muse_connector_address", addr.String()).
		Msg("Updated muse connector address")

	signer.Lock()
	signer.museConnectorAddress = addr
	signer.Unlock()
}

// SetERC20CustodyAddress sets the erc20 custody address
func (signer *Signer) SetERC20CustodyAddress(addr ethcommon.Address) {
	// noop
	if (addr == ethcommon.Address{}) || signer.er20CustodyAddress == addr {
		return
	}

	signer.Logger().Std.Info().
		Str("signer.old_erc20_custody_address", signer.er20CustodyAddress.String()).
		Str("signer.new_erc20_custody_address", addr.String()).
		Msg("Updated erc20 custody address")

	signer.Lock()
	signer.er20CustodyAddress = addr
	signer.Unlock()
}

// SetGatewayAddress sets the gateway address
func (signer *Signer) SetGatewayAddress(addrRaw string) {
	addr := ethcommon.HexToAddress(addrRaw)

	// noop
	if (addr == ethcommon.Address{}) || signer.gatewayAddress == addr {
		return
	}

	signer.Logger().Std.Info().
		Str("signer.old_gateway_address", signer.gatewayAddress.String()).
		Str("signer.new_gateway_address", addr.String()).
		Msg("Updated gateway address")

	signer.Lock()
	signer.gatewayAddress = addr
	signer.Unlock()
}

// GetMuseConnectorAddress returns the muse connector address
func (signer *Signer) GetMuseConnectorAddress() ethcommon.Address {
	return signer.museConnectorAddress
}

// GetERC20CustodyAddress returns the erc20 custody address
func (signer *Signer) GetERC20CustodyAddress() ethcommon.Address {
	return signer.er20CustodyAddress
}

// GetGatewayAddress returns the gateway address
func (signer *Signer) GetGatewayAddress() string {
	return signer.gatewayAddress.String()
}

// Sign given data, and metadata (gas, nonce, etc)
// returns a signed transaction, sig bytes, hash bytes, and error
func (signer *Signer) Sign(
	ctx context.Context,
	data []byte,
	to ethcommon.Address,
	amount *big.Int,
	gas Gas,
	nonce uint64,
	height uint64,
) (*ethtypes.Transaction, []byte, []byte, error) {
	signer.Logger().Std.Debug().
		Str("tss_pub_key", signer.TSS().PubKey().AddressEVM().String()).
		Msg("Signing evm transaction")

	chainID := big.NewInt(signer.Chain().ChainId)
	tx, err := newTx(chainID, data, to, amount, gas, nonce)
	if err != nil {
		return nil, nil, nil, err
	}

	hashBytes := signer.client.Hash(tx).Bytes()

	sig, err := signer.TSS().Sign(ctx, hashBytes, height, nonce, signer.Chain().ChainId)
	if err != nil {
		return nil, nil, nil, err
	}

	log.Debug().Msgf("Sign: Signature: %s", hex.EncodeToString(sig[:]))
	pubk, err := crypto.SigToPub(hashBytes, sig[:])
	if err != nil {
		signer.Logger().Std.Error().Err(err).Msgf("SigToPub error")
	}

	addr := crypto.PubkeyToAddress(*pubk)
	signer.Logger().Std.Info().Msgf("Sign: Ecrecovery of signature: %s", addr.Hex())
	signedTX, err := tx.WithSignature(signer.client.Signer, sig[:])
	if err != nil {
		return nil, nil, nil, err
	}

	return signedTX, sig[:], hashBytes[:], nil
}

func newTx(
	_ *big.Int,
	data []byte,
	to ethcommon.Address,
	amount *big.Int,
	gas Gas,
	nonce uint64,
) (*ethtypes.Transaction, error) {
	if err := gas.validate(); err != nil {
		return nil, errors.Wrap(err, "invalid gas parameters")
	}

	// https://github.com/RWAs-labs/muse/issues/3221
	//if gas.isLegacy() {
	return ethtypes.NewTx(&ethtypes.LegacyTx{
		To:       &to,
		Value:    amount,
		Data:     data,
		GasPrice: gas.Price,
		Gas:      gas.Limit,
		Nonce:    nonce,
	}), nil
	//}
	//
	//return ethtypes.NewTx(&ethtypes.DynamicFeeTx{
	//	ChainID:   chainID,
	//	To:        &to,
	//	Value:     amount,
	//	Data:      data,
	//	GasFeeCap: gas.Price,
	//	GasTipCap: gas.PriorityFee,
	//	Gas:       gas.Limit,
	//	Nonce:     nonce,
	//}), nil
}

func (signer *Signer) broadcast(ctx context.Context, tx *ethtypes.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, broadcastTimeout)
	defer cancel()

	return signer.client.SendTransaction(ctx, tx)
}

// TryProcessOutbound - signer interface implementation
// This function will attempt to build and sign an evm transaction using the TSS signer.
// It will then broadcast the signed transaction to the outbound chain.
// TODO(revamp): simplify function
func (signer *Signer) TryProcessOutbound(
	ctx context.Context,
	cctx *crosschaintypes.CrossChainTx,
	musecoreClient interfaces.MusecoreClient,
	height uint64,
) {
	outboundID := base.OutboundIDFromCCTX(cctx)
	signer.MarkOutbound(outboundID, true)

	// end outbound process on panic
	defer func() {
		signer.MarkOutbound(outboundID, false)
		if r := recover(); r != nil {
			signer.Logger().
				Std.Error().
				Str(logs.FieldMethod, "TryProcessOutbound").
				Str(logs.FieldCctx, cctx.Index).
				Interface("panic", r).
				Str("stack_trace", string(debug.Stack())).
				Msg("caught panic error")
		}
	}()

	// prepare logger and a few local variables
	var (
		params = cctx.GetCurrentOutboundParam()
		myID   = musecoreClient.GetKeys().GetOperatorAddress()
		logger = signer.Logger().Std.With().
			Str(logs.FieldMethod, "TryProcessOutbound").
			Int64(logs.FieldChain, signer.Chain().ChainId).
			Uint64(logs.FieldNonce, params.TssNonce).
			Str(logs.FieldCctx, cctx.Index).
			Str("cctx.receiver", params.Receiver).
			Str("cctx.amount", params.Amount.String()).
			Str("signer", myID.String()).
			Logger()
	)
	logger.Info().Msgf("TryProcessOutbound")

	// retrieve app context
	app, err := zctx.FromContext(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("error getting app context")
		return
	}

	// Setup Transaction input
	txData, skipTx, err := NewOutboundData(ctx, cctx, height, logger)
	if err != nil {
		logger.Err(err).Msg("error setting up transaction input fields")
		return
	}

	if skipTx {
		return
	}

	toChain, err := app.GetChain(txData.toChainID.Int64())
	switch {
	case err != nil:
		logger.Error().Err(err).Msgf("error getting toChain %d", txData.toChainID.Int64())
		return
	case toChain.IsMuse():
		// should not happen
		logger.Error().Msgf("unable to TryProcessOutbound when toChain is museChain (%d)", toChain.ID())
		return
	}

	// sign outbound
	tx, err := signer.SignOutboundFromCCTX(
		ctx,
		logger,
		cctx,
		txData,
		musecoreClient,
		toChain,
	)
	if err != nil {
		logger.Err(err).Msg("error signing outbound")
		return
	}

	// attach tx hash to logger and print log
	logger = logger.With().Str(logs.FieldTx, tx.Hash().Hex()).Logger()
	logger.Info().Msg("Successful keysign")

	// Broadcast Signed Tx
	signer.BroadcastOutbound(ctx, tx, cctx, logger, musecoreClient, txData)
}

// SignOutboundFromCCTX signs an outbound transaction from a given cctx
// TODO: simplify logic with all if else
// https://github.com/RWAs-labs/muse/issues/2050
func (signer *Signer) SignOutboundFromCCTX(
	ctx context.Context,
	logger zerolog.Logger,
	cctx *crosschaintypes.CrossChainTx,
	outboundData *OutboundData,
	musecoreClient interfaces.MusecoreClient,
	toChain zctx.Chain,
) (*ethtypes.Transaction, error) {
	if !signer.PassesCompliance(cctx) {
		// restricted cctx
		return signer.SignCancel(ctx, outboundData)
	} else if cctx.InboundParams.CoinType == coin.CoinType_Cmd {
		// admin command
		to := ethcommon.HexToAddress(cctx.GetCurrentOutboundParam().Receiver)
		if to == (ethcommon.Address{}) {
			return nil, fmt.Errorf("invalid receiver %s", cctx.GetCurrentOutboundParam().Receiver)
		}
		msg := strings.Split(cctx.RelayedMessage, ":")
		if len(msg) != 2 {
			return nil, fmt.Errorf("invalid message %s", msg)
		}
		// cmd field is used to determine whether to execute ERC20 whitelist or migrate TSS funds given that the coin type
		// from the cctx is coin.CoinType_Cmd
		cmd := msg[0]
		// params field is used to pass input parameters for command requests, currently it is used to pass the ERC20
		// contract address when a whitelist command is requested
		params := msg[1]
		return signer.SignAdminTx(ctx, outboundData, cmd, params)
	} else if cctx.ProtocolContractVersion == crosschaintypes.ProtocolContractVersion_V2 {
		// call sign outbound from cctx for v2 protocol contracts
		return signer.SignOutboundFromCCTXV2(ctx, cctx, outboundData)
	} else if IsPendingOutboundFromMuseChain(cctx, musecoreClient) {
		switch cctx.InboundParams.CoinType {
		case coin.CoinType_Gas:
			logger.Info().Msgf(
				"SignGasWithdraw: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				outboundData.gas.Price,
			)
			return signer.SignGasWithdraw(ctx, outboundData)
		case coin.CoinType_ERC20:
			logger.Info().Msgf(
				"SignERC20Withdraw: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				outboundData.gas.Price,
			)
			return signer.SignERC20Withdraw(ctx, outboundData)
		case coin.CoinType_Muse:
			logger.Info().Msgf(
				"SignConnectorOnReceive: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				outboundData.gas.Price,
			)
			return signer.SignConnectorOnReceive(ctx, outboundData)
		}
	} else if cctx.CctxStatus.Status == crosschaintypes.CctxStatus_PendingRevert && cctx.OutboundParams[0].ReceiverChainId == musecoreClient.Chain().ChainId {
		switch cctx.InboundParams.CoinType {
		case coin.CoinType_Muse:
			logger.Info().Msgf(
				"SignConnectorOnRevert: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(), cctx.GetCurrentOutboundParam().TssNonce,
				outboundData.gas.Price,
			)
			outboundData.srcChainID = big.NewInt(cctx.OutboundParams[0].ReceiverChainId)
			outboundData.toChainID = big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId)
			return signer.SignConnectorOnRevert(ctx, outboundData)
		case coin.CoinType_Gas:
			logger.Info().Msgf(
				"SignGasWithdraw: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				outboundData.gas.Price,
			)
			return signer.SignGasWithdraw(ctx, outboundData)
		case coin.CoinType_ERC20:
			logger.Info().Msgf("SignERC20Withdraw: %d => %d, nonce %d, gasPrice %d",
				cctx.InboundParams.SenderChainId,
				toChain.ID(),
				cctx.GetCurrentOutboundParam().TssNonce,
				outboundData.gas.Price,
			)
			return signer.SignERC20Withdraw(ctx, outboundData)
		}
	} else if cctx.CctxStatus.Status == crosschaintypes.CctxStatus_PendingRevert {
		logger.Info().Msgf(
			"SignConnectorOnRevert: %d => %d, nonce %d, gasPrice %d",
			cctx.InboundParams.SenderChainId,
			toChain.ID(),
			cctx.GetCurrentOutboundParam().TssNonce,
			outboundData.gas.Price,
		)
		outboundData.srcChainID = big.NewInt(cctx.OutboundParams[0].ReceiverChainId)
		outboundData.toChainID = big.NewInt(cctx.GetCurrentOutboundParam().ReceiverChainId)
		return signer.SignConnectorOnRevert(ctx, outboundData)
	} else if cctx.CctxStatus.Status == crosschaintypes.CctxStatus_PendingOutbound {
		logger.Info().Msgf(
			"SignConnectorOnReceive: %d => %d, nonce %d, gasPrice %d",
			cctx.InboundParams.SenderChainId,
			toChain.ID(),
			cctx.GetCurrentOutboundParam().TssNonce,
			outboundData.gas.Price,
		)
		return signer.SignConnectorOnReceive(ctx, outboundData)
	}

	return nil, fmt.Errorf("SignOutboundFromCCTX: can't determine how to sign outbound from cctx %s", cctx.String())
}

// BroadcastOutbound signed transaction through evm rpc client
func (signer *Signer) BroadcastOutbound(
	ctx context.Context,
	tx *ethtypes.Transaction,
	cctx *crosschaintypes.CrossChainTx,
	logger zerolog.Logger,
	musecoreClient interfaces.MusecoreClient,
	txData *OutboundData,
) {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		logger.Err(err).Msg("error getting app context")
		return
	}

	toChain, err := app.GetChain(txData.toChainID.Int64())
	switch {
	case err != nil:
		logger.Error().Err(err).Msgf("error getting toChain %d", txData.toChainID.Int64())
		return
	case toChain.IsMuse():
		// should not happen
		logger.Error().Msgf("unable to broadcast when toChain is museChain (%d)", toChain.ID())
		return
	case tx == nil:
		logger.Warn().Msgf("BroadcastOutbound: no tx to broadcast %s", cctx.Index)
		return
	}

	var (
		outboundHash = tx.Hash().Hex()
		nonce        = cctx.GetCurrentOutboundParam().TssNonce
	)

	// define broadcast function
	broadcast := func() error {
		// get latest TSS account nonce
		latestNonce, err := signer.client.NonceAt(ctx, signer.TSS().PubKey().AddressEVM(), nil)
		if err != nil {
			return errors.Wrap(err, "unable to get latest TSS account nonce")
		}

		// if TSS nonce is higher than CCTX nonce, there is no need to broadcast
		// this avoids foreseeable "nonce too low" error and unnecessary tracker report
		// Note: the latest finalized nonce is used here, not the pending nonce, making it possible to replacing pending txs
		if latestNonce > nonce {
			logger.Info().Uint64("latest_nonce", latestNonce).Msg("cctx nonce is too low, skip broadcasting tx")
			return nil
		}

		// broadcast success, report to tracker
		if err = signer.broadcast(ctx, tx); err == nil {
			signer.reportToOutboundTracker(ctx, musecoreClient, toChain.ID(), nonce, outboundHash, logger)
			return nil
		}

		// handle different broadcast errors
		retry, report := musecore.HandleBroadcastError(err, nonce, toChain.ID(), outboundHash)
		if report {
			signer.reportToOutboundTracker(ctx, musecoreClient, toChain.ID(), nonce, outboundHash, logger)
			return nil
		}
		if retry {
			return errors.Wrap(err, "unable to broadcast tx, retrying")
		}

		// no re-broadcast, no report, stop retry
		// e.g. "replacement transaction underpriced"
		return nil
	}

	// broadcast transaction with backoff to tolerate RPC error
	bo := backoff.NewConstantBackOff(broadcastBackoff)
	boWithMaxRetries := backoff.WithMaxRetries(bo, broadcastRetries)
	if err := retry.DoWithBackoff(broadcast, boWithMaxRetries); err != nil {
		logger.Error().Err(err).Msgf("unable to broadcast EVM outbound")
	}

	logger.Info().Msg("broadcasted EVM outbound")
}

// IsPendingOutboundFromMuseChain checks if the sender chain is MuseChain and if status is PendingOutbound
// TODO(revamp): move to another package more general for cctx functions
func IsPendingOutboundFromMuseChain(
	cctx *crosschaintypes.CrossChainTx,

	musecoreClient interfaces.MusecoreClient,
) bool {
	return cctx.InboundParams.SenderChainId == musecoreClient.Chain().ChainId &&
		cctx.CctxStatus.Status == crosschaintypes.CctxStatus_PendingOutbound
}

// ErrorMsg returns a error message for SignConnectorOnReceive failure with cctx data
func ErrorMsg(cctx *crosschaintypes.CrossChainTx) string {
	return fmt.Sprintf(
		"signer SignConnectorOnReceive error: nonce %d chain %d",
		cctx.GetCurrentOutboundParam().TssNonce,
		cctx.GetCurrentOutboundParam().ReceiverChainId,
	)
}
