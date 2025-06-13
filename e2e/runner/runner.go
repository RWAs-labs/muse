package runner

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	erc20custodyv2 "github.com/RWAs-labs/protocol-contracts/pkg/erc20custody.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	museeth "github.com/RWAs-labs/protocol-contracts/pkg/muse.eth.sol"
	museconnectoreth "github.com/RWAs-labs/protocol-contracts/pkg/museconnector.eth.sol"
	connectormevm "github.com/RWAs-labs/protocol-contracts/pkg/museconnectormevm.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/systemcontract.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/wmuse.sol"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/ton"

	"github.com/RWAs-labs/muse/e2e/config"
	"github.com/RWAs-labs/muse/e2e/contracts/contextapp"
	"github.com/RWAs-labs/muse/e2e/contracts/erc20"
	"github.com/RWAs-labs/muse/e2e/contracts/mevmswap"
	"github.com/RWAs-labs/muse/e2e/contracts/testdappv2"
	"github.com/RWAs-labs/muse/e2e/txserver"
	"github.com/RWAs-labs/muse/e2e/utils"
	btcclient "github.com/RWAs-labs/muse/museclient/chains/bitcoin/client"
	"github.com/RWAs-labs/muse/pkg/constant"
	"github.com/RWAs-labs/muse/pkg/contracts/sui"
	"github.com/RWAs-labs/muse/pkg/contracts/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/RWAs-labs/muse/pkg/contracts/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	emissionstypes "github.com/RWAs-labs/muse/x/emissions/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	lightclienttypes "github.com/RWAs-labs/muse/x/lightclient/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

type E2ERunnerOption func(*E2ERunner)

// Important ENV
const (
	EnvKeyLocalnetMode = "LOCALNET_MODE"

	LocalnetModeUpgrade = "upgrade"
)

func WithMuseTxServer(txServer *txserver.MuseTxServer) E2ERunnerOption {
	return func(r *E2ERunner) {
		r.MuseTxServer = txServer
	}
}

func WithTestFilter(testFilter *regexp.Regexp) E2ERunnerOption {
	return func(r *E2ERunner) {
		r.TestFilter = testFilter
	}
}

// E2ERunner stores all the clients and addresses needed for E2E test
// Exposes a method to run E2E test
// It also provides some helper functions
type E2ERunner struct {
	// accounts
	Account               config.Account
	TSSAddress            ethcommon.Address
	BTCTSSAddress         btcutil.Address
	SuiTSSAddress         string
	SolanaDeployerAddress solana.PublicKey
	FeeCollectorAddress   types.AccAddress

	// all clients.
	// a reference to this type is required to enable creating a new E2ERunner.
	Clients Clients

	// rpc clients
	MEVMClient   *ethclient.Client
	EVMClient    *ethclient.Client
	BtcRPCClient *btcclient.Client
	SolanaClient *rpc.Client

	// musecored grpc clients
	AuthorityClient    authoritytypes.QueryClient
	CctxClient         crosschaintypes.QueryClient
	FungibleClient     fungibletypes.QueryClient
	AuthClient         authtypes.QueryClient
	BankClient         banktypes.QueryClient
	StakingClient      stakingtypes.QueryClient
	ObserverClient     observertypes.QueryClient
	LightclientClient  lightclienttypes.QueryClient
	DistributionClient distributiontypes.QueryClient
	EmissionsClient    emissionstypes.QueryClient

	// optional muse (cosmos) client
	// typically only in test runners that need it
	// (like admin tests)
	MuseTxServer *txserver.MuseTxServer

	// evm auth
	EVMAuth  *bind.TransactOpts
	MEVMAuth *bind.TransactOpts

	// programs on Solana
	GatewayProgram solana.PublicKey
	SPLAddr        solana.PublicKey

	// TON related
	TONGateway ton.AccountID

	// contract Sui
	SuiGateway *sui.Gateway

	// SuiGatewayUpgradeCap is the upgrade cap used for upgrading the Sui gateway package
	SuiGatewayUpgradeCap string

	// SuiTokenCoinType is the coin type identifying the fungible token for SUI
	SuiTokenCoinType string

	// SuiTokenTreasuryCap is the treasury cap for the SUI token that allows minting, only using in local tests
	SuiTokenTreasuryCap string

	// SuiExample contains the example package information for Sui
	SuiExample config.SuiExample

	// contracts evm
	MuseEthAddr       ethcommon.Address
	MuseEth           *museeth.MuseEth
	ConnectorEthAddr  ethcommon.Address
	ConnectorEth      *museconnectoreth.MuseConnectorEth
	ERC20CustodyAddr  ethcommon.Address
	ERC20Custody      *erc20custodyv2.ERC20Custody
	ERC20Addr         ethcommon.Address
	ERC20             *erc20.ERC20
	EvmTestDAppAddr   ethcommon.Address
	GatewayEVMAddr    ethcommon.Address
	GatewayEVM        *gatewayevm.GatewayEVM
	TestDAppV2EVMAddr ethcommon.Address
	TestDAppV2EVM     *testdappv2.TestDAppV2

	// contracts mevm
	// mrc20 contracts
	ERC20MRC20Addr    ethcommon.Address
	ERC20MRC20        *mrc20.MRC20
	SPLMRC20Addr      ethcommon.Address
	SPLMRC20          *mrc20.MRC20
	ETHMRC20Addr      ethcommon.Address
	ETHMRC20          *mrc20.MRC20
	BTCMRC20Addr      ethcommon.Address
	BTCMRC20          *mrc20.MRC20
	SOLMRC20Addr      ethcommon.Address
	SOLMRC20          *mrc20.MRC20
	TONMRC20Addr      ethcommon.Address
	TONMRC20          *mrc20.MRC20
	SUIMRC20Addr      ethcommon.Address
	SUIMRC20          *mrc20.MRC20
	SuiTokenMRC20Addr ethcommon.Address
	SuiTokenMRC20     *mrc20.MRC20

	// other contracts
	UniswapV2FactoryAddr ethcommon.Address
	UniswapV2Factory     *uniswapv2factory.UniswapV2Factory
	UniswapV2RouterAddr  ethcommon.Address
	UniswapV2Router      *uniswapv2router.UniswapV2Router02
	ConnectorMEVMAddr    ethcommon.Address
	ConnectorMEVM        *connectormevm.MuseConnectorMEVM
	WMuseAddr            ethcommon.Address
	WMuse                *wmuse.WETH9
	MEVMSwapAppAddr      ethcommon.Address
	MEVMSwapApp          *mevmswap.MEVMSwapApp
	ContextAppAddr       ethcommon.Address
	ContextApp           *contextapp.ContextApp
	SystemContractAddr   ethcommon.Address
	SystemContract       *systemcontract.SystemContract
	MevmTestDAppAddr     ethcommon.Address
	GatewayMEVMAddr      ethcommon.Address
	GatewayMEVM          *gatewaymevm.GatewayMEVM
	TestDAppV2MEVMAddr   ethcommon.Address
	TestDAppV2MEVM       *testdappv2.TestDAppV2

	// config
	CctxTimeout    time.Duration
	ReceiptTimeout time.Duration

	// other
	Name             string
	Ctx              context.Context
	CtxCancel        context.CancelCauseFunc
	Logger           *Logger
	BitcoinParams    *chaincfg.Params
	TestFilter       *regexp.Regexp
	mutex            sync.Mutex
	musecoredVersion string
}

func NewE2ERunner(
	ctx context.Context,
	name string,
	ctxCancel context.CancelCauseFunc,
	account config.Account,
	clients Clients,
	logger *Logger,
	opts ...E2ERunnerOption,
) *E2ERunner {
	r := &E2ERunner{
		Name:      name,
		CtxCancel: ctxCancel,

		Account: account,

		FeeCollectorAddress: authtypes.NewModuleAddress(authtypes.FeeCollectorName),

		Clients: clients,

		MEVMClient:         clients.Mevm,
		EVMClient:          clients.Evm,
		AuthorityClient:    clients.Musecore.Authority,
		CctxClient:         clients.Musecore.Crosschain,
		FungibleClient:     clients.Musecore.Fungible,
		AuthClient:         clients.Musecore.Auth,
		BankClient:         clients.Musecore.Bank,
		StakingClient:      clients.Musecore.Staking,
		ObserverClient:     clients.Musecore.Observer,
		LightclientClient:  clients.Musecore.Lightclient,
		DistributionClient: clients.Musecore.Distribution,
		EmissionsClient:    clients.Musecore.Emissions,

		EVMAuth:      clients.EvmAuth,
		MEVMAuth:     clients.MevmAuth,
		BtcRPCClient: clients.BtcRPC,
		SolanaClient: clients.Solana,

		Logger: logger,
	}

	r.Ctx = utils.WithTesting(ctx, r)

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// CopyAddressesFrom copies addresses from another E2ETestRunner that initialized the contracts
func (r *E2ERunner) CopyAddressesFrom(other *E2ERunner) (err error) {
	// copy TSS address
	r.TSSAddress = other.TSSAddress
	r.BTCTSSAddress = other.BTCTSSAddress
	r.SuiTSSAddress = other.SuiTSSAddress

	// copy addresses
	r.MuseEthAddr = other.MuseEthAddr
	r.ConnectorEthAddr = other.ConnectorEthAddr
	r.ERC20CustodyAddr = other.ERC20CustodyAddr
	r.ERC20Addr = other.ERC20Addr
	r.ERC20MRC20Addr = other.ERC20MRC20Addr
	r.ETHMRC20Addr = other.ETHMRC20Addr
	r.BTCMRC20Addr = other.BTCMRC20Addr
	r.SOLMRC20Addr = other.SOLMRC20Addr
	r.TONMRC20Addr = other.TONMRC20Addr
	r.SUIMRC20Addr = other.SUIMRC20Addr
	r.SuiTokenMRC20Addr = other.SuiTokenMRC20Addr
	r.UniswapV2FactoryAddr = other.UniswapV2FactoryAddr
	r.UniswapV2RouterAddr = other.UniswapV2RouterAddr
	r.ConnectorMEVMAddr = other.ConnectorMEVMAddr
	r.WMuseAddr = other.WMuseAddr
	r.EvmTestDAppAddr = other.EvmTestDAppAddr
	r.MEVMSwapAppAddr = other.MEVMSwapAppAddr
	r.ContextAppAddr = other.ContextAppAddr
	r.SystemContractAddr = other.SystemContractAddr
	r.MevmTestDAppAddr = other.MevmTestDAppAddr

	r.GatewayProgram = other.GatewayProgram

	r.TONGateway = other.TONGateway

	r.SuiGateway = other.SuiGateway
	r.SuiGatewayUpgradeCap = other.SuiGatewayUpgradeCap
	r.SuiTokenCoinType = other.SuiTokenCoinType
	r.SuiTokenTreasuryCap = other.SuiTokenTreasuryCap
	r.SuiExample = other.SuiExample

	// create instances of contracts
	r.MuseEth, err = museeth.NewMuseEth(r.MuseEthAddr, r.EVMClient)
	if err != nil {
		return err
	}
	r.ConnectorEth, err = museconnectoreth.NewMuseConnectorEth(r.ConnectorEthAddr, r.EVMClient)
	if err != nil {
		return err
	}
	r.ERC20Custody, err = erc20custodyv2.NewERC20Custody(r.ERC20CustodyAddr, r.EVMClient)
	if err != nil {
		return err
	}
	r.ERC20, err = erc20.NewERC20(r.ERC20Addr, r.EVMClient)
	if err != nil {
		return err
	}
	r.ERC20MRC20, err = mrc20.NewMRC20(r.ERC20MRC20Addr, r.MEVMClient)
	if err != nil {
		return err
	}
	r.ETHMRC20, err = mrc20.NewMRC20(r.ETHMRC20Addr, r.MEVMClient)
	if err != nil {
		return err
	}
	r.BTCMRC20, err = mrc20.NewMRC20(r.BTCMRC20Addr, r.MEVMClient)
	if err != nil {
		return err
	}
	r.SOLMRC20, err = mrc20.NewMRC20(r.SOLMRC20Addr, r.MEVMClient)
	if err != nil {
		return err
	}
	r.TONMRC20, err = mrc20.NewMRC20(r.TONMRC20Addr, r.MEVMClient)
	if err != nil {
		return err
	}
	r.SUIMRC20, err = mrc20.NewMRC20(r.SUIMRC20Addr, r.MEVMClient)
	if err != nil {
		return err
	}
	r.SuiTokenMRC20, err = mrc20.NewMRC20(r.SuiTokenMRC20Addr, r.MEVMClient)
	if err != nil {
		return err
	}

	r.UniswapV2Factory, err = uniswapv2factory.NewUniswapV2Factory(r.UniswapV2FactoryAddr, r.MEVMClient)
	if err != nil {
		return err
	}
	r.UniswapV2Router, err = uniswapv2router.NewUniswapV2Router02(r.UniswapV2RouterAddr, r.MEVMClient)
	if err != nil {
		return err
	}
	r.ConnectorMEVM, err = connectormevm.NewMuseConnectorMEVM(r.ConnectorMEVMAddr, r.MEVMClient)
	if err != nil {
		return err
	}
	r.WMuse, err = wmuse.NewWETH9(r.WMuseAddr, r.MEVMClient)
	if err != nil {
		return err
	}

	r.MEVMSwapApp, err = mevmswap.NewMEVMSwapApp(r.MEVMSwapAppAddr, r.MEVMClient)
	if err != nil {
		return err
	}
	r.ContextApp, err = contextapp.NewContextApp(r.ContextAppAddr, r.MEVMClient)
	if err != nil {
		return err
	}
	r.SystemContract, err = systemcontract.NewSystemContract(r.SystemContractAddr, r.MEVMClient)
	if err != nil {
		return err
	}

	// v2 contracts
	r.GatewayEVMAddr = other.GatewayEVMAddr
	r.GatewayEVM, err = gatewayevm.NewGatewayEVM(r.GatewayEVMAddr, r.EVMClient)
	if err != nil {
		return err
	}

	r.TestDAppV2EVMAddr = other.TestDAppV2EVMAddr
	r.TestDAppV2EVM, err = testdappv2.NewTestDAppV2(r.TestDAppV2EVMAddr, r.EVMClient)
	if err != nil {
		return err
	}

	r.GatewayMEVMAddr = other.GatewayMEVMAddr
	r.GatewayMEVM, err = gatewaymevm.NewGatewayMEVM(r.GatewayMEVMAddr, r.MEVMClient)
	if err != nil {
		return err
	}
	r.TestDAppV2MEVMAddr = other.TestDAppV2MEVMAddr
	r.TestDAppV2MEVM, err = testdappv2.NewTestDAppV2(r.TestDAppV2MEVMAddr, r.MEVMClient)
	if err != nil {
		return err
	}

	return nil
}

// Lock locks the mutex
func (r *E2ERunner) Lock() {
	r.mutex.Lock()
}

// Unlock unlocks the mutex
func (r *E2ERunner) Unlock() {
	r.mutex.Unlock()
}

// PrintContractAddresses prints the addresses of the contracts
// the printed contracts are grouped in a mevm and evm section
// there is a padding used to print the addresses at the same position
func (r *E2ERunner) PrintContractAddresses() {
	r.Logger.Print(" --- 📜Solana addresses ---")
	r.Logger.Print("GatewayProgram: %s", r.GatewayProgram.String())
	r.Logger.Print("SPL:            %s", r.SPLAddr.String())

	r.Logger.Print(" --- 📜TON addresses ---")
	if !r.TONGateway.IsZero() {
		r.Logger.Print("Gateway:        %s", r.TONGateway.ToRaw())
	} else {
		r.Logger.Print("Gateway:        not set! 💤")
	}

	r.Logger.Print(" --- 📜Sui addresses ---")
	if r.SuiGateway != nil {
		r.Logger.Print("GatewayPackageID: %s", r.SuiGateway.PackageID())
		r.Logger.Print("GatewayObjectID:  %s", r.SuiGateway.ObjectID())
		r.Logger.Print("GatewayUpgradeCap: %s", r.SuiGatewayUpgradeCap)
	} else {
		r.Logger.Print("💤 Sui tests disabled")
	}

	// mevm contracts
	r.Logger.Print(" --- 📜mEVM contracts ---")
	r.Logger.Print("SystemContract: %s", r.SystemContractAddr.Hex())
	r.Logger.Print("ETHMRC20:       %s", r.ETHMRC20Addr.Hex())
	r.Logger.Print("ERC20MRC20:     %s", r.ERC20MRC20Addr.Hex())
	r.Logger.Print("BTCMRC20:       %s", r.BTCMRC20Addr.Hex())
	r.Logger.Print("SOLMRC20:       %s", r.SOLMRC20Addr.Hex())
	r.Logger.Print("SPLMRC20:       %s", r.SPLMRC20Addr.Hex())
	r.Logger.Print("TONMRC20:       %s", r.TONMRC20Addr.Hex())
	r.Logger.Print("SUIMRC20:       %s", r.SUIMRC20Addr.Hex())
	r.Logger.Print("SuiTokenMRC20:  %s", r.SuiTokenMRC20Addr.Hex())
	r.Logger.Print("UniswapFactory: %s", r.UniswapV2FactoryAddr.Hex())
	r.Logger.Print("UniswapRouter:  %s", r.UniswapV2RouterAddr.Hex())
	r.Logger.Print("ConnectorMEVM:  %s", r.ConnectorMEVMAddr.Hex())
	r.Logger.Print("WMuse:          %s", r.WMuseAddr.Hex())
	r.Logger.Print("GatewayMEVM:    %s", r.GatewayMEVMAddr.Hex())
	r.Logger.Print("TestDAppV2MEVM: %s", r.TestDAppV2MEVMAddr.Hex())

	// evm contracts
	r.Logger.Print(" --- 📜EVM contracts ---")
	r.Logger.Print("MuseEth:        %s", r.MuseEthAddr.Hex())
	r.Logger.Print("ConnectorEth:   %s", r.ConnectorEthAddr.Hex())
	r.Logger.Print("ERC20Custody:   %s", r.ERC20CustodyAddr.Hex())
	r.Logger.Print("ERC20:          %s", r.ERC20Addr.Hex())
	r.Logger.Print("GatewayEVM:     %s", r.GatewayEVMAddr.Hex())
	r.Logger.Print("TestDAppV2EVM:  %s", r.TestDAppV2EVMAddr.Hex())

	r.Logger.Print(" --- 📜Legacy contracts ---")

	r.Logger.Print("MEVMSwapApp:    %s", r.MEVMSwapAppAddr.Hex())
	r.Logger.Print("ContextApp:     %s", r.ContextAppAddr.Hex())
	r.Logger.Print("TestDappMEVM:   %s", r.MevmTestDAppAddr.Hex())
	r.Logger.Print("TestDappEVM:    %s", r.EvmTestDAppAddr.Hex())
}

// IsRunningUpgrade returns true if the test is running an upgrade test suite.
func (r *E2ERunner) IsRunningUpgrade() bool {
	return os.Getenv(EnvKeyLocalnetMode) == LocalnetModeUpgrade
}

// Errorf logs an error message. Mimics the behavior of testing.T.Errorf
func (r *E2ERunner) Errorf(format string, args ...any) {
	r.Logger.Error(format, args...)
}

// FailNow implemented to mimic the behavior of testing.T.FailNow
func (r *E2ERunner) FailNow() {
	err := fmt.Errorf("(*E2ERunner).FailNow for runner %q. Exiting", r.Name)

	r.Logger.Error("Failure: %s", err.Error())
	r.CtxCancel(err)

	// this panic ensures that the test routine exits fast.
	// it should be caught and handled gracefully so long
	// as the test is being executed by RunE2ETest().
	panic(err)
}

func (r *E2ERunner) requireTxSuccessful(receipt *ethtypes.Receipt, msgAndArgs ...any) {
	utils.RequireTxSuccessful(r, receipt, msgAndArgs...)
}

// EVMAddress is shorthand to get the EVM address of the account
func (r *E2ERunner) EVMAddress() ethcommon.Address {
	return r.Account.EVMAddress()
}

func (r *E2ERunner) GetSolanaPrivKey() solana.PrivateKey {
	privkey, err := solana.PrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
	require.NoError(r, err)
	return privkey
}

func (r *E2ERunner) GetMusecoredVersion() string {
	if r.musecoredVersion != "" {
		return r.musecoredVersion
	}
	nodeInfo, err := r.Clients.Musecore.GetNodeInfo(r.Ctx)
	require.NoError(r, err, "get node info")
	r.musecoredVersion = constant.NormalizeVersion(nodeInfo.ApplicationVersion.Version)
	return r.musecoredVersion
}
