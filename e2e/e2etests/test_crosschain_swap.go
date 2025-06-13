package e2etests

import (
	"math/big"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestCrosschainSwap(r *runner.E2ERunner, _ []string) {
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()
	r.MEVMAuth.GasLimit = 10000000

	// TODO: move into setup and skip it if already initialized
	// https://github.com/RWAs-labs/muse-private/issues/88
	// it is kept as is for now to be consistent with the old implementation
	// if the tx fails due to already initialized, it will be ignored
	_, err := r.UniswapV2Factory.CreatePair(r.MEVMAuth, r.ERC20MRC20Addr, r.BTCMRC20Addr)
	if err != nil {
		r.Logger.Print("ℹ️ create pair error %s", err.Error())
	}

	txERC20MRC20Approve, err := r.ERC20MRC20.Approve(r.MEVMAuth, r.UniswapV2RouterAddr, big.NewInt(1e18))
	require.NoError(r, err)

	txBTCApprove, err := r.BTCMRC20.Approve(r.MEVMAuth, r.UniswapV2RouterAddr, big.NewInt(1e18))
	require.NoError(r, err)

	// Fund MEVMSwapApp with gas MRC20s
	txTransferETH, err := r.ETHMRC20.Transfer(r.MEVMAuth, r.MEVMSwapAppAddr, big.NewInt(1e7))
	require.NoError(r, err)

	txTransferBTC, err := r.BTCMRC20.Transfer(r.MEVMAuth, r.MEVMSwapAppAddr, big.NewInt(1e6))
	require.NoError(r, err)

	ensureTxReceipt := func(tx *ethtypes.Transaction, failMessage string) {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt, failMessage)
	}

	ensureTxReceipt(txERC20MRC20Approve, "MRC20 ERC20 approve failed")
	ensureTxReceipt(txBTCApprove, "BTC approve failed")
	ensureTxReceipt(txTransferETH, "ETH MRC20 transfer failed")
	ensureTxReceipt(txTransferBTC, "BTC MRC20 transfer failed")

	// Add 100 erc20 mrc20 liq and 0.001 BTC
	txAddLiquidity, err := r.UniswapV2Router.AddLiquidity(
		r.MEVMAuth,
		r.ERC20MRC20Addr,
		r.BTCMRC20Addr,
		big.NewInt(1e8),
		big.NewInt(1e8),
		big.NewInt(1e8),
		big.NewInt(1e5),
		r.EVMAddress(),
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	require.NoError(r, err)
	ensureTxReceipt(txAddLiquidity, "add liq failed")

	// msg would be [MEVMSwapAppAddr, memobytes]
	// memobytes is dApp specific; see the contracts/MEVMSwapApp.sol for details
	msg := []byte{}
	msg = append(msg, r.MEVMSwapAppAddr.Bytes()...)
	memobytes, err := r.MEVMSwapApp.EncodeMemo(
		&bind.CallOpts{},
		r.BTCMRC20Addr,
		[]byte(r.GetBtcAddress().EncodeAddress()),
	)
	require.NoError(r, err)

	r.Logger.Info("memobytes(%d) %x", len(memobytes), memobytes)
	msg = append(msg, memobytes...)

	r.Logger.Info("***** First test: ERC20 -> BTC")
	// Should deposit ERC20 for swap, swap for BTC and withdraw BTC
	txHash := r.LegacyDepositERC20WithAmountAndMessage(r.EVMAddress(), big.NewInt(8e7), msg)
	cctx1 := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

	// check the cctx status
	utils.RequireCCTXStatus(r, cctx1, types.CctxStatus_OutboundMined)

	// cctx1 index acts like the inboundHash for the second cctx (the one that withdraws BTC)
	cctx2 := utils.WaitCctxMinedByInboundHash(r.Ctx, cctx1.Index, r.CctxClient, r.Logger, r.CctxTimeout)

	// check the cctx status
	utils.RequireCCTXStatus(r, cctx2, types.CctxStatus_OutboundMined)

	r.Logger.Info("cctx2 outbound tx hash %s", cctx2.GetCurrentOutboundParam().Hash)

	r.Logger.Info("******* Second test: BTC -> ERC20MRC20")
	// list deployer utxos
	utxos := r.ListUTXOs()

	r.Logger.Info("#utxos %d", len(utxos))
	r.Logger.Info("memo address %s", r.ERC20MRC20Addr)

	memo, err := r.MEVMSwapApp.EncodeMemo(&bind.CallOpts{}, r.ERC20MRC20Addr, r.EVMAddress().Bytes())
	require.NoError(r, err)

	memo = append(r.MEVMSwapAppAddr.Bytes(), memo...)
	r.Logger.Info("memo length %d", len(memo))

	txID, err := r.SendToTSSWithMemo(0.01, memo)
	require.NoError(r, err)

	cctx3 := utils.WaitCctxMinedByInboundHash(r.Ctx, txID.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx3, types.CctxStatus_OutboundMined)

	r.Logger.Info("cctx3 index %s", cctx3.Index)
	r.Logger.Info("  inbound tx hash %s", cctx3.InboundParams.ObservedHash)
	r.Logger.Info("  status %s", cctx3.CctxStatus.Status.String())
	r.Logger.Info("  status msg: %s", cctx3.CctxStatus.StatusMessage)

	cctx4 := utils.WaitCctxMinedByInboundHash(r.Ctx, cctx3.Index, r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx4, types.CctxStatus_OutboundMined)

	r.Logger.Info("cctx4 index %s", cctx4.Index)
	r.Logger.Info("  outbound tx hash %s", cctx4.GetCurrentOutboundParam().Hash)
	r.Logger.Info("  status %s", cctx4.CctxStatus.Status.String())

	{
		r.Logger.Info("******* Third test: BTC -> ETH with contract call reverted; should refund BTC")
		// the following memo will result in a revert in the contract call as targetMRC20 is set to DeployerAddress
		// which is apparently not a MRC20 contract; the UNISWAP call will revert
		memo, err := r.MEVMSwapApp.EncodeMemo(&bind.CallOpts{}, r.EVMAddress(), r.EVMAddress().Bytes())
		require.NoError(r, err)

		memo = append(r.MEVMSwapAppAddr.Bytes(), memo...)
		r.Logger.Info("memo length %d", len(memo))

		amount := 0.1
		txid, err := r.SendToTSSWithMemo(amount, memo)
		require.NoError(r, err)

		cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txid.String(), r.CctxClient, r.Logger, r.CctxTimeout)
		r.Logger.Info("cctx3 index %s", cctx.Index)
		r.Logger.Info("  inbound tx hash %s", cctx.InboundParams.ObservedHash)
		r.Logger.Info("  status %s", cctx.CctxStatus.Status.String())
		r.Logger.Info("  status msg: %s", cctx.CctxStatus.StatusMessage)

		utils.RequireCCTXStatus(r, cctx, types.CctxStatus_Reverted)

		outboundHash, err := chainhash.NewHashFromStr(cctx.GetCurrentOutboundParam().Hash)
		require.NoError(r, err)

		txraw, err := r.BtcRPCClient.GetRawTransactionVerbose(r.Ctx, outboundHash)
		require.NoError(r, err)

		r.Logger.Info("out txid %s", txraw.Txid)
		for _, vout := range txraw.Vout {
			r.Logger.Info("  vout %d", vout.N)
			r.Logger.Info("  value %f", vout.Value)
			r.Logger.Info("  scriptPubKey %s", vout.ScriptPubKey.Hex)
			r.Logger.Info("  p2wpkh address: %s", utils.ScriptPKToAddress(vout.ScriptPubKey.Hex, r.BitcoinParams))
		}
	}
}
