package runner

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/config"
	suicontract "github.com/RWAs-labs/muse/e2e/contracts/sui"
	"github.com/RWAs-labs/muse/e2e/txserver"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/constant"
	musesui "github.com/RWAs-labs/muse/pkg/contracts/sui"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

const changeTypeCreated = "created"

// RequestSuiFromFaucet requests SUI tokens from the faucet for the runner account
func (r *E2ERunner) RequestSuiFromFaucet(faucetURL, recipient string) {
	header := map[string]string{}
	err := sui.RequestSuiFromFaucet(faucetURL, recipient, header)
	require.NoError(r, err, "sui faucet request to %s", faucetURL)
}

// SetupSui initializes the gateway package on Sui and initialize the chain params on MuseChain
func (r *E2ERunner) SetupSui(faucetURL string) {
	r.Logger.Print("⚙️ initializing gateway package on Sui")

	deployerSigner, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	deployerAddress := deployerSigner.Address()

	// fund deployer
	r.RequestSuiFromFaucet(faucetURL, deployerAddress)

	// import deployer private key and select it as active address
	r.suiSetupDeployerAccount()

	// fund the TSS
	// request twice from the faucet to ensure TSS has enough funds for the first withdraw
	// TODO: this step might no longer necessary if a custom solution is implemented for the TSS funding
	r.RequestSuiFromFaucet(faucetURL, r.SuiTSSAddress)
	r.RequestSuiFromFaucet(faucetURL, r.SuiTSSAddress)

	// deploy gateway package
	whitelistCapID, withdrawCapID := r.suiDeployGateway()

	// update gateway package ID in Move.toml
	r.suiPatchMoveConfig()

	// deploy SUI mrc20
	r.deploySUIMRC20()

	// deploy fake USDC and whitelist it
	fakeUSDCCoinType := r.suiDeployFakeUSDC()
	r.whitelistSuiFakeUSDC(deployerSigner, fakeUSDCCoinType, whitelistCapID)

	// deploy example contract with on_call function
	r.suiDeployExample()

	// send withdraw cap to TSS
	r.suiSendWithdrawCapToTSS(deployerSigner, withdrawCapID)

	// set the chain params
	err = r.setSuiChainParams()
	require.NoError(r, err)
}

// suiSetupDeployerAccount imports a Sui deployer private key using the sui keytool import command
// and sets the deployer address as the active address.
func (r *E2ERunner) suiSetupDeployerAccount() {
	deployerSigner, err := r.Account.SuiSigner()
	require.NoError(r, err, "unable to get deployer signer")

	var (
		deployerAddress    = deployerSigner.Address()
		deployerPrivKeyHex = r.Account.RawPrivateKey.String()
	)

	// convert private key to bech32
	deployerPrivKeySecp256k1, err := musesui.PrivateKeyBech32Secp256k1FromHex(deployerPrivKeyHex)
	require.NoError(r, err)

	// import deployer private key using sui keytool import
	// #nosec G204, inputs are controlled in E2E test
	cmdImport := exec.Command("sui", "keytool", "import", deployerPrivKeySecp256k1, "secp256k1")
	require.NoError(r, cmdImport.Run(), "unable to import sui deployer private key")

	// switch to deployer address using sui client switch
	// #nosec G204, inputs are controlled in E2E test
	cmdSwitch := exec.Command("sui", "client", "switch", "--address", deployerAddress)
	require.NoError(r, cmdSwitch.Run(), "unable to switch to deployer address")

	// ensure the deployer address is active
	// #nosec G204, inputs are controlled in E2E test
	cmdList := exec.Command("sui", "client", "active-address")
	output, err := cmdList.Output()
	require.NoError(r, err)
	require.Equal(r, deployerAddress, strings.TrimSpace(string(output)))
}

// suiDeployGateway deploys the SUI gateway package on Sui
func (r *E2ERunner) suiDeployGateway() (whitelistCapID, withdrawCapID string) {
	const (
		filterGatewayType      = "gateway::Gateway"
		filterWithdrawCapType  = "gateway::WithdrawCap"
		filterWhitelistCapType = "gateway::WhitelistCap"
		filterUpgradeCapType   = "0x2::package::UpgradeCap"
	)

	objectTypeFilters := []string{
		filterGatewayType,
		filterWhitelistCapType,
		filterWithdrawCapType,
		filterUpgradeCapType,
	}
	packageID, objectIDs := r.suiDeployPackage(
		[]string{suicontract.GatewayBytecodeBase64(), suicontract.EVMBytecodeBase64()},
		objectTypeFilters,
	)

	gatewayID, ok := objectIDs[filterGatewayType]
	require.True(r, ok, "gateway object not found")

	whitelistCapID, ok = objectIDs[filterWhitelistCapType]
	require.True(r, ok, "whitelistCap object not found")

	withdrawCapID, ok = objectIDs[filterWithdrawCapType]
	require.True(r, ok, "withdrawCap object not found")

	r.SuiGatewayUpgradeCap, ok = objectIDs[filterUpgradeCapType]
	require.True(r, ok, "upgradeCap object not found")

	// set sui gateway
	r.SuiGateway = musesui.NewGateway(packageID, gatewayID)

	return whitelistCapID, withdrawCapID
}

// deploySUIMRC20 deploys the SUI mrc20 on MuseChain
func (r *E2ERunner) deploySUIMRC20() {
	// send message to deploy SUI mrc20
	liqCap := math.NewUint(10e18)
	adminAddr := r.MuseTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName)
	_, err := r.MuseTxServer.BroadcastTx(utils.AdminPolicyName, fungibletypes.NewMsgDeployFungibleCoinMRC20(
		adminAddr,
		"",
		chains.SuiLocalnet.ChainId,
		9,
		"SUI",
		"SUI",
		coin.CoinType_Gas,
		10000,
		&liqCap,
	))
	require.NoError(r, err)

	// set the address in the store
	r.SetupSUIMRC20()
}

// suiDeployFakeUSDC deploys the FakeUSDC contract on Sui
// it returns the treasuryCap object ID that allows to mint tokens
func (r *E2ERunner) suiDeployFakeUSDC() string {
	packageID, objectIDs := r.suiDeployPackage([]string{suicontract.FakeUSDCBytecodeBase64()}, []string{"TreasuryCap"})

	treasuryCap, ok := objectIDs["TreasuryCap"]
	require.True(r, ok, "treasuryCap not found")

	coinType := packageID + "::fake_usdc::FAKE_USDC"

	// strip 0x from packageID
	coinType = coinType[2:]

	// set asset value for mrc20 and treasuryCap object ID
	r.SuiTokenCoinType = coinType
	r.SuiTokenTreasuryCap = treasuryCap

	return coinType
}

// suiDeployExample deploys the example package on Sui
func (r *E2ERunner) suiDeployExample() {
	const (
		filterGlobalConfigType = "connected::GlobalConfig"
		filterPartnerType      = "connected::Partner"
		filterClockType        = "connected::Clock"
	)

	objectTypeFilters := []string{filterGlobalConfigType, filterPartnerType, filterClockType}
	packageID, objectIDs := r.suiDeployPackage(
		[]string{suicontract.ExampleFungibleTokenBytecodeBase64(), suicontract.ExampleConnectedBytecodeBase64()},
		objectTypeFilters,
	)
	r.Logger.Info("deployed example package with packageID: %s", packageID)

	globalConfigID, ok := objectIDs[filterGlobalConfigType]
	require.True(r, ok, "globalConfig object not found")

	partnerID, ok := objectIDs[filterPartnerType]
	require.True(r, ok, "partner object not found")

	clockID, ok := objectIDs[filterClockType]
	require.True(r, ok, "clock object not found")

	r.SuiExample = config.SuiExample{
		PackageID:      config.DoubleQuotedString(packageID),
		TokenType:      config.DoubleQuotedString(packageID + "::token::TOKEN"),
		GlobalConfigID: config.DoubleQuotedString(globalConfigID),
		PartnerID:      config.DoubleQuotedString(partnerID),
		ClockID:        config.DoubleQuotedString(clockID),
	}
}

// suiDeployPackage is a helper function that deploys a package on Sui
// It returns the packageID and a map of object types to their IDs
func (r *E2ERunner) suiDeployPackage(bytecodeBase64s []string, objectTypeFilters []string) (string, map[string]string) {
	client := r.Clients.Sui

	deployerSigner, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	deployerAddress := deployerSigner.Address()

	publishTx, err := client.Publish(r.Ctx, models.PublishRequest{
		Sender:          deployerAddress,
		CompiledModules: bytecodeBase64s,
		Dependencies: []string{
			"0x1", // Sui Framework
			"0x2", // Move Standard Library
		},
		GasBudget: "5000000000",
	})
	require.NoError(r, err, "create publish tx")

	signature, err := deployerSigner.SignTxBlock(publishTx)
	require.NoError(r, err, "sign transaction")

	resp, err := client.SuiExecuteTransactionBlock(r.Ctx, models.SuiExecuteTransactionBlockRequest{
		TxBytes:   publishTx.TxBytes,
		Signature: []string{signature},
		Options: models.SuiTransactionBlockOptions{
			ShowEffects:        true,
			ShowBalanceChanges: true,
			ShowEvents:         true,
			ShowObjectChanges:  true,
		},
		RequestType: "WaitForLocalExecution",
	})
	require.NoError(r, err)

	// find packageID
	var packageID string
	for _, change := range resp.ObjectChanges {
		if change.Type == "published" {
			packageID = change.PackageId
		}
	}
	require.NotEmpty(r, packageID, "packageID not found")

	// find objects by type filters
	objectIDs := make(map[string]string)
	for _, filter := range objectTypeFilters {
		for _, change := range resp.ObjectChanges {
			if change.Type == changeTypeCreated && strings.Contains(change.ObjectType, filter) {
				objectIDs[filter] = change.ObjectId
			}
		}
	}

	return packageID, objectIDs
}

// whitelistSuiFakeUSDC deploys the FakeUSDC mrc20 on MuseChain and whitelist it
func (r *E2ERunner) whitelistSuiFakeUSDC(signer *musesui.SignerSecp256k1, fakeUSDCCoinType, whitelistCap string) {
	// we use DeployFungibleCoinMRC20 and whitelist manually because whitelist cctx are currently not supported for Sui
	// TODO: change this logic and use MsgWhitelistERC20 once it's supported
	// https://github.com/RWAs-labs/muse/issues/3569

	// deploy mrc20
	liqCap := math.NewUint(10e18)
	res, err := r.MuseTxServer.BroadcastTx(utils.AdminPolicyName, fungibletypes.NewMsgDeployFungibleCoinMRC20(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		fakeUSDCCoinType,
		chains.SuiLocalnet.ChainId,
		6,
		"Sui's FakeUSDC",
		"USDC.SUI",
		coin.CoinType_ERC20,
		10000,
		&liqCap,
	))
	require.NoError(r, err)

	// extract the mrc20 address from event and set the erc20 address in the runner
	deployedEvent, ok := txserver.EventOfType[*fungibletypes.EventMRC20Deployed](res.Events)
	require.True(r, ok, "unable to find deployed mrc20 event")

	r.SuiTokenMRC20Addr = ethcommon.HexToAddress(deployedEvent.Contract)
	require.NotEqualValues(r, ethcommon.Address{}, r.SuiTokenMRC20Addr)
	r.SuiTokenMRC20, err = mrc20.NewMRC20(r.SuiTokenMRC20Addr, r.MEVMClient)
	require.NoError(r, err)

	// whitelist mrc20
	tx, err := r.Clients.Sui.MoveCall(r.Ctx, models.MoveCallRequest{
		Signer:          signer.Address(),
		PackageObjectId: r.SuiGateway.PackageID(),
		Module:          "gateway",
		Function:        "whitelist",
		TypeArguments:   []any{"0x" + fakeUSDCCoinType},
		Arguments:       []any{r.SuiGateway.ObjectID(), whitelistCap},
		GasBudget:       "5000000000",
	})
	require.NoError(r, err)

	r.suiExecuteTx(signer, tx)
}

// set the chain params for Sui
func (r *E2ERunner) setSuiChainParams() error {
	if r.MuseTxServer == nil {
		return errors.New("MuseTxServer is not initialized")
	}

	creator := r.MuseTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName)

	chainID := chains.SuiLocalnet.ChainId

	chainParams := &observertypes.ChainParams{
		ChainId:                     chainID,
		MuseTokenContractAddress:    constant.EVMZeroAddress,
		ConnectorContractAddress:    constant.EVMZeroAddress,
		Erc20CustodyContractAddress: constant.EVMZeroAddress,
		GasPriceTicker:              5,
		WatchUtxoTicker:             0,
		InboundTicker:               2,
		OutboundTicker:              2,
		OutboundScheduleInterval:    2,
		OutboundScheduleLookahead:   5,
		BallotThreshold:             observertypes.DefaultBallotThreshold,
		MinObserverDelegation:       observertypes.DefaultMinObserverDelegation,
		IsSupported:                 true,
		GatewayAddress:              fmt.Sprintf("%s,%s", r.SuiGateway.PackageID(), r.SuiGateway.ObjectID()),
		ConfirmationParams: &observertypes.ConfirmationParams{
			SafeInboundCount:  1,
			SafeOutboundCount: 1,
			FastInboundCount:  1,
			FastOutboundCount: 1,
		},
		ConfirmationCount: 1, // still need to be provided for now
	}
	if err := r.MuseTxServer.UpdateChainParams(chainParams); err != nil {
		return errors.Wrap(err, "unable to broadcast solana chain params tx")
	}

	resetMsg := observertypes.NewMsgResetChainNonces(creator, chainID, 0, 0)
	if _, err := r.MuseTxServer.BroadcastTx(utils.OperationalPolicyName, resetMsg); err != nil {
		return errors.Wrap(err, "unable to broadcast solana chain nonce reset tx")
	}

	query := &observertypes.QueryGetChainParamsForChainRequest{ChainId: chainID}

	const duration = 2 * time.Second

	for i := 0; i < 10; i++ {
		_, err := r.ObserverClient.GetChainParamsForChain(r.Ctx, query)
		if err == nil {
			r.Logger.Print("⚙️ Sui chain params are set")
			return nil
		}

		time.Sleep(duration)
	}

	return errors.New("unable to set Sui chain params")
}

func (r *E2ERunner) suiSendWithdrawCapToTSS(signer *musesui.SignerSecp256k1, withdrawCapID string) {
	tx, err := r.Clients.Sui.TransferObject(r.Ctx, models.TransferObjectRequest{
		Signer:    signer.Address(),
		ObjectId:  withdrawCapID,
		Recipient: r.SuiTSSAddress,
		GasBudget: "5000000000",
	})
	require.NoError(r, err)

	r.suiExecuteTx(signer, tx)
}
