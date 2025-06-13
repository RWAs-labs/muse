package chains

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/RWAs-labs/protocol-contracts/pkg/erc20custody.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/museconnector.non-eth.sol"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/RWAs-labs/muse/cmd/musetool/config"
	"github.com/RWAs-labs/muse/cmd/musetool/context"
	"github.com/RWAs-labs/muse/museclient/musecore"
	clienttypes "github.com/RWAs-labs/muse/museclient/types"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/constant"
	"github.com/RWAs-labs/muse/pkg/crypto"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func resolveRPC(chain chains.Chain, cfg *config.Config) string {
	return map[chains.Network]string{
		chains.Network_eth:     cfg.EthereumRPC,
		chains.Network_base:    cfg.BaseRPC,
		chains.Network_polygon: cfg.PolygonRPC,
		chains.Network_bsc:     cfg.BscRPC,
	}[chain.Network]
}

func GetEvmClient(ctx *context.Context, chain chains.Chain) (*ethclient.Client, error) {
	evmRRC := resolveRPC(chain, ctx.GetConfig())
	if evmRRC == "" {
		return nil, fmt.Errorf("rpc not found for chain %d network %s", chain.ChainId, chain.Network)
	}
	rpcClient, err := ethrpc.DialHTTP(evmRRC)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to eth rpc: %w", err)
	}
	return ethclient.NewClient(rpcClient), nil
}

func GetEvmTx(
	ctx *context.Context,
	evmClient *ethclient.Client,
	inboundHash string,
	chain chains.Chain,
) (*ethtypes.Transaction, *ethtypes.Receipt, error) {
	goCtx := ctx.GetContext()
	// Fetch transaction from the inbound
	hash := ethcommon.HexToHash(inboundHash)
	tx, isPending, err := evmClient.TransactionByHash(goCtx, hash)
	if err != nil {
		return nil, nil, fmt.Errorf("tx not found on chain: %w,chainID: %d", err, chain.ChainId)
	}
	if isPending {
		return nil, nil, fmt.Errorf("tx is still pending on chain: %d", chain.ChainId)
	}
	receipt, err := evmClient.TransactionReceipt(goCtx, hash)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get receipt: %w, tx hash: %s", err, inboundHash)
	}
	return tx, receipt, nil
}

func MuseTokenVoteV1(
	event *museconnector.MuseConnectorNonEthMuseSent,
	observationChain int64,
) *crosschaintypes.MsgVoteInbound {
	// note that this is most likely muse chain
	destChain, found := chains.GetChainFromChainID(event.DestinationChainId.Int64(), []chains.Chain{})
	if !found {
		return nil
	}

	destAddr := clienttypes.BytesToEthHex(event.DestinationAddress)
	sender := event.MuseTxSenderAddress.Hex()
	message := base64.StdEncoding.EncodeToString(event.Message)

	return musecore.GetInboundVoteMessage(
		sender,
		observationChain,
		event.SourceTxOriginAddress.Hex(),
		destAddr,
		destChain.ChainId,
		sdkmath.NewUintFromBigInt(event.MuseValueAndGas),
		message,
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		event.DestinationGasLimit.Uint64(),
		coin.CoinType_Muse,
		"",
		"",
		uint64(event.Raw.Index),
		crosschaintypes.InboundStatus_SUCCESS,
	)
}

func Erc20VoteV1(
	event *erc20custody.ERC20CustodyDeposited,
	sender ethcommon.Address,
	observationChain int64,
	musecoreChainID int64,
) *crosschaintypes.MsgVoteInbound {
	// donation check
	if bytes.Equal(event.Message, []byte(constant.DonationMessage)) {
		return nil
	}

	return musecore.GetInboundVoteMessage(
		sender.Hex(),
		observationChain,
		"",
		clienttypes.BytesToEthHex(event.Recipient),
		musecoreChainID,
		sdkmath.NewUintFromBigInt(event.Amount),
		hex.EncodeToString(event.Message),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		1_500_000,
		coin.CoinType_ERC20,
		event.Asset.String(),
		"",
		uint64(event.Raw.Index),
		crosschaintypes.InboundStatus_SUCCESS,
	)
}

func GasVoteV1(
	tx *ethtypes.Transaction,
	sender ethcommon.Address,
	blockNumber uint64,
	senderChainID int64,
	musecoreChainID int64,
) *crosschaintypes.MsgVoteInbound {
	message := string(tx.Data())
	data, _ := hex.DecodeString(message)
	if bytes.Equal(data, []byte(constant.DonationMessage)) {
		return nil
	}

	return musecore.GetInboundVoteMessage(
		sender.Hex(),
		senderChainID,
		sender.Hex(),
		sender.Hex(),
		musecoreChainID,
		sdkmath.NewUintFromString(tx.Value().String()),
		message,
		tx.Hash().Hex(),
		blockNumber,
		90_000,
		coin.CoinType_Gas,
		"",
		"",
		0, // not a smart contract call
		crosschaintypes.InboundStatus_SUCCESS,
	)
}

func DepositInboundVoteV2(event *gatewayevm.GatewayEVMDeposited,
	senderChainID int64,
	musecoreChainID int64) *crosschaintypes.MsgVoteInbound {
	// if event.Asset is zero, it's a native token
	coinType := coin.CoinType_ERC20
	if crypto.IsEmptyAddress(event.Asset) {
		coinType = coin.CoinType_Gas
	}

	// to maintain compatibility with previous gateway version, deposit event with a non-empty payload is considered as a call
	isCrossChainCall := false
	if len(event.Payload) > 0 {
		isCrossChainCall = true
	}

	return crosschaintypes.NewMsgVoteInbound(
		"",
		event.Sender.Hex(),
		senderChainID,
		"",
		event.Receiver.Hex(),
		musecoreChainID,
		sdkmath.NewUintFromBigInt(event.Amount),
		hex.EncodeToString(event.Payload),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		musecore.PostVoteInboundCallOptionsGasLimit,
		coinType,
		event.Asset.Hex(),
		uint64(event.Raw.Index),
		crosschaintypes.ProtocolContractVersion_V2,
		false, // currently not relevant since calls are not arbitrary
		crosschaintypes.InboundStatus_SUCCESS,
		crosschaintypes.ConfirmationMode_SAFE,
		crosschaintypes.WithEVMRevertOptions(event.RevertOptions),
		crosschaintypes.WithCrossChainCall(isCrossChainCall),
	)
}

func DepositAndCallInboundVoteV2(event *gatewayevm.GatewayEVMDepositedAndCalled,
	senderChainID int64,
	musecoreChainID int64) *crosschaintypes.MsgVoteInbound {
	// if event.Asset is zero, it's a native token
	coinType := coin.CoinType_ERC20
	if crypto.IsEmptyAddress(event.Asset) {
		coinType = coin.CoinType_Gas
	}

	return crosschaintypes.NewMsgVoteInbound(
		"",
		event.Sender.Hex(),
		senderChainID,
		"",
		event.Receiver.Hex(),
		musecoreChainID,
		sdkmath.NewUintFromBigInt(event.Amount),
		hex.EncodeToString(event.Payload),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		1_500_000,
		coinType,
		event.Asset.Hex(),
		uint64(event.Raw.Index),
		crosschaintypes.ProtocolContractVersion_V2,
		false, // currently not relevant since calls are not arbitrary
		crosschaintypes.InboundStatus_SUCCESS,
		crosschaintypes.ConfirmationMode_SAFE,
		crosschaintypes.WithEVMRevertOptions(event.RevertOptions),
		crosschaintypes.WithCrossChainCall(true),
	)
}

func CallInboundVoteV2(event *gatewayevm.GatewayEVMCalled,
	senderChainID int64,
	musecoreChainID int64) *crosschaintypes.MsgVoteInbound {
	return crosschaintypes.NewMsgVoteInbound(
		"",
		event.Sender.Hex(),
		senderChainID,
		"",
		event.Receiver.Hex(),
		musecoreChainID,
		sdkmath.ZeroUint(),
		hex.EncodeToString(event.Payload),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		musecore.PostVoteInboundCallOptionsGasLimit,
		coin.CoinType_NoAssetCall,
		"",
		uint64(event.Raw.Index),
		crosschaintypes.ProtocolContractVersion_V2,
		false, // currently not relevant since calls are not arbitrary
		crosschaintypes.InboundStatus_SUCCESS,
		crosschaintypes.ConfirmationMode_SAFE,
		crosschaintypes.WithEVMRevertOptions(event.RevertOptions),
	)
}
