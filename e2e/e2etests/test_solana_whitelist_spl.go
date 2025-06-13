package e2etests

import (
	sdkmath "cosmossdk.io/math"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/txserver"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/chains"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestSolanaWhitelistSPL(r *runner.E2ERunner, _ []string) {
	// Deploy a new SPL
	r.Logger.Info("Deploying new SPL")

	// load deployer private key
	privkey := r.GetSolanaPrivKey()

	// deploy SPL token, but don't whitelist in gateway
	spl := r.DeploySPL(&privkey, false)

	// check that whitelist entry doesn't exist for this spl
	seed := [][]byte{[]byte("whitelist"), spl.PublicKey().Bytes()}
	whitelistEntryPDA, _, err := solana.FindProgramAddress(seed, r.GatewayProgram)
	require.NoError(r, err)

	whitelistEntryInfo, err := r.SolanaClient.GetAccountInfoWithOpts(
		r.Ctx,
		whitelistEntryPDA,
		&rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed},
	)
	require.Error(r, err)
	require.Nil(r, whitelistEntryInfo)

	// whitelist sol mrc20
	r.Logger.Info("whitelisting spl on new network")
	res, err := r.MuseTxServer.BroadcastTx(utils.AdminPolicyName, crosschaintypes.NewMsgWhitelistERC20(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		spl.PublicKey().String(),
		chains.SolanaLocalnet.ChainId,
		"TESTSPL",
		"TESTSPL",
		6,
		100000,
		sdkmath.NewUintFromString("100000000000000000000000000"),
	))
	require.NoError(r, err)

	event, ok := txserver.EventOfType[*crosschaintypes.EventERC20Whitelist](res.Events)
	require.True(r, ok, "no EventERC20Whitelist in %s", res.TxHash)
	erc20mrc20Addr := event.Mrc20Address
	whitelistCCTXIndex := event.WhitelistCctxIndex

	err = r.MuseTxServer.InitializeLiquidityCaps(erc20mrc20Addr)
	require.NoError(r, err)

	// ensure CCTX created
	resCCTX, err := r.CctxClient.Cctx(r.Ctx, &crosschaintypes.QueryGetCctxRequest{Index: whitelistCCTXIndex})
	require.NoError(r, err)

	cctx := resCCTX.CrossChainTx
	r.Logger.CCTX(*cctx, "whitelist_cctx")

	// wait for the whitelist cctx to be mined
	r.WaitForMinedCCTXFromIndex(whitelistCCTXIndex)

	// check that whitelist entry exists for this spl
	whitelistEntryInfo, err = r.SolanaClient.GetAccountInfoWithOpts(
		r.Ctx,
		whitelistEntryPDA,
		&rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed},
	)
	require.NoError(r, err)
	require.NotNil(r, whitelistEntryInfo)
}
