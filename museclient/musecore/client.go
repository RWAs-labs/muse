// Package musecore provides the client to interact with musecore node via GRPC.
package musecore

import (
	"fmt"
	"strings"
	"sync"

	etherminttypes "github.com/RWAs-labs/ethermint/types"
	cometbftrpc "github.com/cometbft/cometbft/rpc/client"
	cometbfthttp "github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/types"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/RWAs-labs/muse/app"
	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	"github.com/RWAs-labs/muse/museclient/config"
	keyinterfaces "github.com/RWAs-labs/muse/museclient/keys/interfaces"
	"github.com/RWAs-labs/muse/museclient/logs"
	"github.com/RWAs-labs/muse/pkg/authz"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/fanout"
	musecorerpc "github.com/RWAs-labs/muse/pkg/rpc"
)

var _ interfaces.MusecoreClient = &Client{}

// Client is the client to send tx to musecore
type Client struct {
	musecorerpc.Clients

	logger zerolog.Logger
	config config.ClientConfiguration

	cosmosClientContext cosmosclient.Context
	cometBFTClient      cometbftrpc.Client

	blockHeight   int64
	accountNumber map[authz.KeyType]uint64
	seqNumber     map[authz.KeyType]uint64

	encodingCfg etherminttypes.EncodingConfig
	keys        keyinterfaces.ObserverKeys
	chainID     string
	chain       chains.Chain

	// blocksFanout that receives new block events from Musecore via websockets
	blocksFanout *fanout.FanOut[ctypes.EventDataNewBlock]

	mu sync.RWMutex
}

var unsecureGRPC = grpc.WithTransportCredentials(insecure.NewCredentials())

type constructOpts struct {
	customCometBFT bool
	cometBFTClient cometbftrpc.Client

	customAccountRetriever bool
	accountRetriever       cosmosclient.AccountRetriever
}

type Opt func(cfg *constructOpts)

// WithCometBFTClient sets custom CometBFT client
func WithCometBFTClient(client cometbftrpc.Client) Opt {
	return func(c *constructOpts) {
		c.customCometBFT = true
		c.cometBFTClient = client
	}
}

// WithCustomAccountRetriever sets custom CometBFT client
func WithCustomAccountRetriever(ac cosmosclient.AccountRetriever) Opt {
	return func(c *constructOpts) {
		c.customAccountRetriever = true
		c.accountRetriever = ac
	}
}

// NewClient create a new instance of Client
func NewClient(
	keys keyinterfaces.ObserverKeys,
	chainIP string,
	signerName string,
	chainID string,
	logger zerolog.Logger,
	opts ...Opt,
) (*Client, error) {
	var constructOptions constructOpts
	for _, opt := range opts {
		opt(&constructOptions)
	}

	museChain, err := chains.MuseChainFromCosmosChainID(chainID)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid chain id %q", chainID)
	}

	log := logger.With().Str(logs.FieldModule, "musecoreClient").Logger()

	cfg := config.ClientConfiguration{
		ChainHost:    cosmosREST(chainIP),
		SignerName:   signerName,
		SignerPasswd: "password",
		ChainRPC:     CometBFTRPC(chainIP),
	}

	encodingCfg := app.MakeEncodingConfig()

	musecoreClients, err := musecorerpc.NewGRPCClients(cosmosGRPC(chainIP), unsecureGRPC)
	if err != nil {
		return nil, errors.Wrap(err, "grpc dial fail")
	}

	accountsMap := make(map[authz.KeyType]uint64)
	seqMap := make(map[authz.KeyType]uint64)
	for _, keyType := range authz.GetAllKeyTypes() {
		accountsMap[keyType] = 0
		seqMap[keyType] = 0
	}

	cosmosContext, err := buildCosmosClientContext(chainID, keys, cfg, encodingCfg, constructOptions)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build cosmos client context")
	}

	cometBFTClientIface := constructOptions.cometBFTClient

	// create a cometbft client if one was not provided in the constructOptions
	if !constructOptions.customCometBFT {
		cometBFTURL := "http://" + CometBFTRPC(chainIP)
		cometBFTClient, err := cometbfthttp.New(cometBFTURL, "/websocket")
		if err != nil {
			return nil, errors.Wrapf(err, "new cometbft client (%s)", cometBFTURL)
		}
		// start websockets
		err = cometBFTClient.WSEvents.Start()
		if err != nil {
			return nil, errors.Wrap(err, "cometbft start")
		}
		cometBFTClientIface = cometBFTClient
	}

	return &Client{
		Clients: musecoreClients,
		logger:  log,
		config:  cfg,

		cosmosClientContext: cosmosContext,
		cometBFTClient:      cometBFTClientIface,

		accountNumber: accountsMap,
		seqNumber:     seqMap,

		encodingCfg: encodingCfg,
		keys:        keys,
		chainID:     chainID,
		chain:       museChain,
	}, nil
}

// buildCosmosClientContext constructs a valid context with all relevant values set
func buildCosmosClientContext(
	chainID string,
	keys keyinterfaces.ObserverKeys,
	config config.ClientConfiguration,
	encodingConfig etherminttypes.EncodingConfig,
	opts constructOpts,
) (cosmosclient.Context, error) {
	if keys == nil {
		return cosmosclient.Context{}, errors.New("client key are not set")
	}

	addr, err := keys.GetAddress()
	if err != nil {
		return cosmosclient.Context{}, errors.Wrap(err, "fail to get address from key")
	}

	var (
		input   = strings.NewReader("")
		client  cosmosclient.CometRPC
		nodeURI string
	)

	// if password is needed, set it as input
	password := keys.GetHotkeyPassword()
	if password != "" {
		input = strings.NewReader(fmt.Sprintf("%[1]s\n%[1]s\n", password))
	}

	// note that in rare cases, this might give FALSE positive
	// (google "golang nil interface comparison")
	client = opts.cometBFTClient
	if !opts.customCometBFT {
		remote := config.ChainRPC
		if !strings.HasPrefix(config.ChainHost, "http") {
			remote = fmt.Sprintf("tcp://%s", remote)
		}

		wsClient, err := cometbfthttp.New(remote, "/websocket")
		if err != nil {
			return cosmosclient.Context{}, err
		}

		client = wsClient
		nodeURI = remote
	}

	var accountRetriever cosmosclient.AccountRetriever
	if opts.customAccountRetriever {
		accountRetriever = opts.accountRetriever
	} else {
		accountRetriever = authtypes.AccountRetriever{}
	}

	return cosmosclient.Context{
		Client:        client,
		NodeURI:       nodeURI,
		FromAddress:   addr,
		ChainID:       chainID,
		Keyring:       keys.GetKeybase(),
		BroadcastMode: "sync",
		HomeDir:       config.ChainHomeFolder,
		FromName:      config.SignerName,

		AccountRetriever: accountRetriever,

		Codec:             encodingConfig.Codec,
		InterfaceRegistry: encodingConfig.InterfaceRegistry,
		TxConfig:          encodingConfig.TxConfig,
		LegacyAmino:       encodingConfig.Amino,

		Input: input,
	}, nil
}

func (c *Client) UpdateChainID(chainID string) error {
	if c.chainID != chainID {
		c.chainID = chainID

		museChain, err := chains.MuseChainFromCosmosChainID(chainID)
		if err != nil {
			return fmt.Errorf("invalid chain id %s, %w", chainID, err)
		}
		c.chain = museChain
	}

	return nil
}

// Chain returns the Chain chain object
func (c *Client) Chain() chains.Chain {
	return c.chain
}

func (c *Client) GetKeys() keyinterfaces.ObserverKeys {
	return c.keys
}

// GetAccountNumberAndSequenceNumber We do not use multiple KeyType for now , but this can be optionally used in the future to seprate TSS signer from Museclient GRantee
func (c *Client) GetAccountNumberAndSequenceNumber(_ authz.KeyType) (uint64, uint64, error) {
	address, err := c.keys.GetAddress()
	if err != nil {
		return 0, 0, err
	}
	return c.cosmosClientContext.AccountRetriever.GetAccountNumberSequence(c.cosmosClientContext, address)
}

// SetAccountNumber sets the account number and sequence number for the given keyType
// todo remove method and make it part of the client constructor.
func (c *Client) SetAccountNumber(keyType authz.KeyType) error {
	address, err := c.keys.GetAddress()
	if err != nil {
		return errors.Wrap(err, "fail to get address")
	}

	accN, seq, err := c.cosmosClientContext.AccountRetriever.GetAccountNumberSequence(c.cosmosClientContext, address)
	if err != nil {
		return errors.Wrap(err, "fail to get account number and sequence number")
	}

	c.accountNumber[keyType] = accN
	c.seqNumber[keyType] = seq

	return nil
}

func cosmosREST(host string) string {
	return fmt.Sprintf("%s:1317", host)
}

func cosmosGRPC(host string) string {
	return fmt.Sprintf("%s:9090", host)
}

func CometBFTRPC(host string) string {
	return fmt.Sprintf("%s:26657", host)
}
