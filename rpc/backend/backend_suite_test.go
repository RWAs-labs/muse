package backend

import (
	"bufio"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	protov2 "google.golang.org/protobuf/proto"

	"github.com/RWAs-labs/ethermint/app"
	"github.com/RWAs-labs/ethermint/crypto/ethsecp256k1"
	"github.com/RWAs-labs/ethermint/crypto/hd"
	"github.com/RWAs-labs/ethermint/indexer"
	"github.com/RWAs-labs/ethermint/tests"
	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/suite"

	"github.com/RWAs-labs/muse/rpc/backend/mocks"
	rpctypes "github.com/RWAs-labs/muse/rpc/types"
)

type BackendTestSuite struct {
	suite.Suite
	backend *Backend
	acc     sdk.AccAddress
	signer  keyring.Signer
}

// testTx is a dummy implementation of cosmos Tx used for testing.
type testTx struct {
}

func (tx testTx) GetMsgs() []sdk.Msg                    { return nil }
func (tx testTx) GetMsgsV2() ([]protov2.Message, error) { return nil, nil }
func (tx testTx) GetSigners() []sdk.AccAddress          { return nil }

func (tx testTx) ValidateBasic() error { return nil }
func (t testTx) ProtoMessage()         { panic("not implemented") }
func (t testTx) Reset()                { panic("not implemented") }

func (t testTx) String() string { panic("not implemented") }

func (t testTx) Bytes() []byte { panic("not implemented") }

func (t testTx) VerifySignature(msg []byte, sig []byte) bool { panic("not implemented") }

func (t testTx) Type() string { panic("not implemented") }

var (
	_ sdk.Tx  = (*testTx)(nil)
	_ sdk.Msg = (*testTx)(nil)
)

func TestBackendTestSuite(t *testing.T) {
	suite.Run(t, new(BackendTestSuite))
}

const ChainID = "musechain_7001-1"

// SetupTest is executed before every BackendTestSuite test
func (suite *BackendTestSuite) SetupTest() {
	ctx := server.NewDefaultContext()
	ctx.Viper.Set("telemetry.global-labels", []interface{}{})

	baseDir := suite.T().TempDir()
	nodeDirName := "node"
	clientDir := filepath.Join(baseDir, nodeDirName, "evmoscli")
	keyRing, err := suite.generateTestKeyring(clientDir)
	if err != nil {
		panic(err)
	}

	// Create Account with set sequence

	suite.acc = sdk.AccAddress(tests.GenerateAddress().Bytes())
	accounts := map[string]client.TestAccount{}
	accounts[suite.acc.String()] = client.TestAccount{
		Address: suite.acc,
		Num:     uint64(1),
		Seq:     uint64(1),
	}

	priv, err := ethsecp256k1.GenerateKey()
	suite.signer = tests.NewSigner(priv)
	suite.Require().NoError(err)

	encodingConfig := app.MakeConfigForTest()
	clientCtx := client.Context{}.WithChainID(ChainID).
		WithHeight(1).
		WithTxConfig(encodingConfig.TxConfig).
		WithKeyringDir(clientDir).
		WithKeyring(keyRing).
		WithAccountRetriever(client.TestAccountRetriever{Accounts: accounts})

	allowUnprotectedTxs := false
	idxer := indexer.NewKVIndexer(dbm.NewMemDB(), ctx.Logger, clientCtx)

	suite.backend = NewBackend(ctx, ctx.Logger, clientCtx, allowUnprotectedTxs, idxer)
	suite.backend.queryClient.QueryClient = mocks.NewEVMQueryClient(suite.T())
	suite.backend.clientCtx.Client = mocks.NewClient(suite.T())
	suite.backend.queryClient.FeeMarket = mocks.NewFeeMarketQueryClient(suite.T())
	suite.backend.ctx = rpctypes.ContextWithHeight(1)

	// Add codec
	encCfg := app.MakeConfigForTest()
	suite.backend.clientCtx.Codec = encCfg.Codec
}

// buildEthereumTx returns an example legacy Ethereum transaction
func (suite *BackendTestSuite) buildEthereumTx() (*evmtypes.MsgEthereumTx, []byte) {
	msgEthereumTx := evmtypes.NewTx(
		suite.backend.chainID,
		uint64(0),
		&common.Address{},
		big.NewInt(0),
		100000,
		big.NewInt(1),
		nil,
		nil,
		nil,
		nil,
	)
	suite.signAndEncodeEthTx(msgEthereumTx)

	txBuilder := suite.backend.clientCtx.TxConfig.NewTxBuilder()
	txBuilder.SetSignatures()
	err := txBuilder.SetMsgs(msgEthereumTx)
	suite.Require().NoError(err)

	bz, err := suite.backend.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	suite.Require().NoError(err)
	return msgEthereumTx, bz
}

func (suite *BackendTestSuite) buildSyntheticTxResult(txHash string) ([]byte, abci.ExecTxResult) {
	testTx := &testTx{}
	txBuilder := suite.backend.clientCtx.TxConfig.NewTxBuilder()
	txBuilder.SetSignatures()
	txBuilder.SetMsgs(testTx)
	bz, _ := suite.backend.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	return bz, abci.ExecTxResult{
		Code: 0,
		Events: []abci.Event{
			{Type: evmtypes.EventTypeEthereumTx, Attributes: []abci.EventAttribute{
				{Key: "ethereumTxHash", Value: txHash},
				{Key: "txIndex", Value: "8888"},
				{Key: "txData", Value: "0x1234"},
				{Key: "amount", Value: "1000"},
				{Key: "txGasUsed", Value: "21000"},
				{Key: "txGasLimit", Value: "21000"},
				{Key: "txHash", Value: ""},
				{Key: "recipient", Value: "0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7"},
			}},
			{
				Type: "message", Attributes: []abci.EventAttribute{
					{Key: "sender", Value: "0x735b14BB79463307AAcBED86DAf3322B1e6226aB"},
					{Key: "txType", Value: "88"},
					{Key: "txNonce", Value: "1"},
				},
			},
		},
	}
}

// buildFormattedBlock returns a formatted block for testing
func (suite *BackendTestSuite) buildFormattedBlock(
	blockRes *tmrpctypes.ResultBlockResults,
	resBlock *tmrpctypes.ResultBlock,
	fullTx bool,
	tx *evmtypes.MsgEthereumTx,
	validator sdk.AccAddress,
	baseFee *big.Int,
) map[string]interface{} {
	header := resBlock.Block.Header
	gasLimit := int64(^uint32(0)) // for `MaxGas = -1` (DefaultConsensusParams)
	gasUsed := new(big.Int).SetUint64(uint64(blockRes.TxsResults[0].GasUsed))

	root := common.Hash{}.Bytes()
	receipt := ethtypes.NewReceipt(root, false, gasUsed.Uint64())
	bloom := ethtypes.CreateBloom(ethtypes.Receipts{receipt})

	ethRPCTxs := []interface{}{}
	if tx != nil {
		if fullTx {
			rpcTx, err := rpctypes.NewRPCTransaction(
				tx.AsTransaction(),
				common.BytesToHash(header.Hash()),
				uint64(header.Height),
				uint64(0),
				baseFee,
				suite.backend.chainID,
			)
			suite.Require().NoError(err)
			ethRPCTxs = []interface{}{rpcTx}
		} else {
			ethRPCTxs = []interface{}{common.HexToHash(tx.Hash)}
		}
	}

	return rpctypes.FormatBlock(
		header,
		resBlock.Block.Size(),
		gasLimit,
		gasUsed,
		ethRPCTxs,
		bloom,
		common.BytesToAddress(validator.Bytes()),
		baseFee,
	)
}

func (suite *BackendTestSuite) generateTestKeyring(clientDir string) (keyring.Keyring, error) {
	buf := bufio.NewReader(os.Stdin)
	encCfg := app.MakeConfigForTest()
	return keyring.New(
		sdk.KeyringServiceName(),
		keyring.BackendTest,
		clientDir,
		buf,
		encCfg.Codec,
		[]keyring.Option{hd.EthSecp256k1Option()}...)
}

func (suite *BackendTestSuite) signAndEncodeEthTx(msgEthereumTx *evmtypes.MsgEthereumTx) []byte {
	from, priv := tests.NewAddrKey()
	signer := tests.NewSigner(priv)

	queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
	RegisterParamsWithoutHeader(queryClient, 1)

	ethSigner := ethtypes.LatestSigner(suite.backend.ChainConfig())
	msgEthereumTx.From = from.String()
	err := msgEthereumTx.Sign(ethSigner, signer)
	suite.Require().NoError(err)

	tx, err := msgEthereumTx.BuildTx(suite.backend.clientCtx.TxConfig.NewTxBuilder(), "amuse")
	suite.Require().NoError(err)

	txEncoder := suite.backend.clientCtx.TxConfig.TxEncoder()
	txBz, err := txEncoder(tx)
	suite.Require().NoError(err)

	return txBz
}
