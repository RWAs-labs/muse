package keeper_test

import (
	"fmt"
	"github.com/RWAs-labs/muse/pkg/contracts/sui"
	"math/big"
	"strconv"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/constant"
	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	crosschainkeeper "github.com/RWAs-labs/muse/x/crosschain/keeper"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// SetupStateForProcessLogsMuseSent sets up additional state required for processing logs for MuseSent events
// This sets up the gas coin, mrc20 contract, gas price, mrc20 pool.
// This should be used in conjunction with SetupStateForProcessLogs for processing MuseSent events
func SetupStateForProcessLogsMuseSent(
	t *testing.T,
	ctx sdk.Context,
	k *crosschainkeeper.Keeper,
	zk keepertest.MuseKeepers,
	sdkk keepertest.SDKKeepers,
	chain chains.Chain,
	admin string,
) {

	assetAddress := sample.EthAddress().String()
	gasMRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chain.ChainId, "ethereum", "ETH")
	mrc20Addr := deployMRC20(
		t,
		ctx,
		zk.FungibleKeeper,
		sdkk.EvmKeeper,
		chain.ChainId,
		"ethereum",
		assetAddress,
		"ETH",
	)

	_, err := zk.FungibleKeeper.UpdateMRC20ProtocolFlatFee(ctx, gasMRC20, big.NewInt(withdrawFee))
	require.NoError(t, err)
	_, err = zk.FungibleKeeper.UpdateMRC20ProtocolFlatFee(ctx, mrc20Addr, big.NewInt(withdrawFee))
	require.NoError(t, err)

	k.SetGasPrice(ctx, crosschaintypes.GasPrice{
		ChainId:     chain.ChainId,
		MedianIndex: 0,
		Prices:      []uint64{gasPrice},
	})
	setupMRC20Pool(
		t,
		ctx,
		zk.FungibleKeeper,
		sdkk.BankKeeper,
		mrc20Addr,
	)
}

// SetupStateForProcessLogs sets up observer state for required for processing logs
// It deploys system contracts, sets up TSS, gas price, chain nonce's, pending nonce's.These are all required to create a cctx from a log
func SetupStateForProcessLogs(
	t *testing.T,
	ctx sdk.Context,
	k *crosschainkeeper.Keeper,
	zk keepertest.MuseKeepers,
	sdkk keepertest.SDKKeepers,
	chain chains.Chain,
) {

	deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
	tss := sample.Tss()
	zk.ObserverKeeper.SetTSS(ctx, tss)
	k.SetGasPrice(ctx, crosschaintypes.GasPrice{
		ChainId: chain.ChainId,
		Prices:  []uint64{100},
	})

	zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{
		ChainId: chain.ChainId,
		Nonce:   0,
	})
	zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{
		NonceLow:  0,
		NonceHigh: 0,
		ChainId:   chain.ChainId,
		Tss:       tss.TssPubkey,
	})
}

func TestParseMRC20WithdrawalEvent(t *testing.T) {
	t.Run("unable to parse an event with an invalid address in event log", func(t *testing.T) {
		for i, log := range sample.InvalidMRC20WithdrawToExternalReceipt(t).Logs {
			event, err := crosschainkeeper.ParseMRC20WithdrawalEvent(*log)
			if i < 3 {
				require.ErrorContains(t, err, "event signature mismatch")
				require.Nil(t, event)
				continue
			}
			require.NoError(t, err)
			require.NotNil(t, event)
			require.Equal(t, "1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3", string(event.To))
		}
	})

	t.Run("successfully parse event for a valid BTC withdrawal", func(t *testing.T) {
		for i, log := range sample.ValidMRC20WithdrawToBTCReceipt(t).Logs {
			event, err := crosschainkeeper.ParseMRC20WithdrawalEvent(*log)
			if i < 3 {
				require.ErrorContains(t, err, "event signature mismatch")
				require.Nil(t, event)
				continue
			}
			require.NoError(t, err)
			require.NotNil(t, event)
			require.Equal(t, "0x33EaD83db0D0c682B05ead61E8d8f481Bb1B4933", event.From.Hex())
			require.Equal(t, "bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu", string(event.To))
		}
	})

	t.Run("successfully parse valid event for ETH withdrawal", func(t *testing.T) {
		for i, log := range sample.ValidMrc20WithdrawToETHReceipt(t).Logs {
			event, err := crosschainkeeper.ParseMRC20WithdrawalEvent(*log)
			if i != 11 {
				require.ErrorContains(t, err, "event signature mismatch")
				require.Nil(t, event)
				continue
			}
			require.NoError(t, err)
			require.NotNil(t, event)
			require.Equal(t, "0x5daBFdd153Aaab4a970fD953DcFEEE8BF6Bb946E", ethcommon.BytesToAddress(event.To).String())
			require.Equal(t, "0x8E0f8E7E9E121403e72151d00F4937eACB2D9Ef3", event.From.Hex())
		}
	})

	t.Run("failed to parse event with a valid address but no topic present", func(t *testing.T) {
		for _, log := range sample.ValidMRC20WithdrawToBTCReceipt(t).Logs {
			log.Topics = nil
			event, err := crosschainkeeper.ParseMRC20WithdrawalEvent(*log)
			require.ErrorContains(t, err, "invalid log - no topics")
			require.Nil(t, event)
		}
	})
}

func TestValidateMrc20WithdrawEvent(t *testing.T) {
	t.Run("successfully validate a valid BTC withdrawal event", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		btcMainNetWithdrawalEvent, err := crosschainkeeper.ParseMRC20WithdrawalEvent(
			*sample.ValidMRC20WithdrawToBTCReceipt(t).Logs[3],
		)
		require.NoError(t, err)
		err = k.ValidateMRC20WithdrawEvent(
			ctx,
			btcMainNetWithdrawalEvent,
			chains.BitcoinMainnet.ChainId,
			coin.CoinType_Gas,
		)
		require.NoError(t, err)
	})

	t.Run("successfully validate a valid SOL withdrawal event", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		// 1000000 lamports is the minimum amount (rent exempt) that can be withdrawn
		chainID := chains.SolanaMainnet.ChainId
		to := []byte(sample.SolanaAddress(t))
		value := big.NewInt(constant.SolanaWalletRentExempt)
		solWithdrawalEvent := sample.MRC20Withdrawal(to, value)

		// 1000000 lamports can be withdrawn
		err := k.ValidateMRC20WithdrawEvent(ctx, solWithdrawalEvent, chainID, coin.CoinType_Gas)
		require.NoError(t, err)
	})

	t.Run("successfully validate a small amount of SPL withdrawal event", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		// set SPL token amount to 1
		chainID := chains.SolanaMainnet.ChainId
		to := []byte(sample.SolanaAddress(t))
		solWithdrawalEvent := sample.MRC20Withdrawal(to, big.NewInt(1))

		// should withdraw successfully
		err := k.ValidateMRC20WithdrawEvent(ctx, solWithdrawalEvent, chainID, coin.CoinType_ERC20)
		require.NoError(t, err)
	})

	t.Run("unable to validate a btc withdrawal event with an invalid amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		btcMainNetWithdrawalEvent, err := crosschainkeeper.ParseMRC20WithdrawalEvent(
			*sample.ValidMRC20WithdrawToBTCReceipt(t).Logs[3],
		)
		require.NoError(t, err)

		// 1000 satoshis is the minimum amount that can be withdrawn
		btcMainNetWithdrawalEvent.Value = big.NewInt(constant.BTCWithdrawalDustAmount)
		err = k.ValidateMRC20WithdrawEvent(
			ctx,
			btcMainNetWithdrawalEvent,
			chains.BitcoinMainnet.ChainId,
			coin.CoinType_Gas,
		)
		require.NoError(t, err)

		// 999 satoshis cannot be withdrawn
		btcMainNetWithdrawalEvent.Value = big.NewInt(constant.BTCWithdrawalDustAmount - 1)
		err = k.ValidateMRC20WithdrawEvent(
			ctx,
			btcMainNetWithdrawalEvent,
			chains.BitcoinMainnet.ChainId,
			coin.CoinType_Gas,
		)
		require.ErrorContains(t, err, "less than dust amount")
	})

	t.Run("unable to validate a event with an invalid chain ID", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		btcMainNetWithdrawalEvent, err := crosschainkeeper.ParseMRC20WithdrawalEvent(
			*sample.ValidMRC20WithdrawToBTCReceipt(t).Logs[3],
		)
		require.NoError(t, err)
		err = k.ValidateMRC20WithdrawEvent(
			ctx,
			btcMainNetWithdrawalEvent,
			chains.BitcoinTestnet.ChainId,
			coin.CoinType_Gas,
		)
		require.ErrorContains(t, err, "invalid address")
	})

	t.Run("unable to validate an unsupported address type", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		btcMainNetWithdrawalEvent, err := crosschainkeeper.ParseMRC20WithdrawalEvent(
			*sample.ValidMRC20WithdrawToBTCReceipt(t).Logs[3],
		)
		require.NoError(t, err)
		btcMainNetWithdrawalEvent.To = []byte("04b2891ba8cb491828db3ebc8a780d43b169e7b3974114e6e50f9bab6ec" +
			"63c2f20f6d31b2025377d05c2a704d3bd799d0d56f3a8543d79a01ab6084a1cb204f260")
		err = k.ValidateMRC20WithdrawEvent(
			ctx,
			btcMainNetWithdrawalEvent,
			chains.BitcoinMainnet.ChainId,
			coin.CoinType_Gas,
		)
		require.ErrorContains(t, err, "unsupported Bitcoin address")
	})

	t.Run("unable to validate an event with an invalid solana address", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		// create a withdrawal event with an invalid address (contains invalid character 'l')
		to := []byte("DCAK36VfExkPdAkYUQg6ewgxyinvcEyPLyHjRbmveKFl")
		value := big.NewInt(constant.SolanaWalletRentExempt)
		solWithdrawalEvent := sample.MRC20Withdrawal(to, value)

		err := k.ValidateMRC20WithdrawEvent(ctx, solWithdrawalEvent, chains.SolanaMainnet.ChainId, coin.CoinType_Gas)
		require.ErrorContains(t, err, "invalid address")
	})

	t.Run("unable to validate a SOL withdrawal event with an invalid amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		// 1000000 lamports is the minimum amount (rent exempt) that can be withdrawn
		chainID := chains.SolanaMainnet.ChainId
		to := []byte(sample.SolanaAddress(t))
		value := big.NewInt(constant.SolanaWalletRentExempt - 1)
		solWithdrawalEvent := sample.MRC20Withdrawal(to, value)

		// 999999 lamports cannot be withdrawn
		err := k.ValidateMRC20WithdrawEvent(ctx, solWithdrawalEvent, chainID, coin.CoinType_Gas)
		require.ErrorContains(t, err, "less than rent exempt")
	})

	t.Run("unable to validate an event with an invalid sui address", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		// create a withdrawal event with an invalid address (contains additional character 'aa')
		value := big.NewInt(1000000)
		suiWithdrawalEvent := sample.MRC20Withdrawal(
			[]byte("0x25db16c3ca555f6702c07860503107bb73cce9f6c1d6df00464529db15d5a5abaa"),
			value,
		)

		err := k.ValidateMRC20WithdrawEvent(ctx, suiWithdrawalEvent, chains.SuiMainnet.ChainId, coin.CoinType_Gas)
		require.ErrorContains(t, err, "invalid Sui address")
	})

	t.Run("validate valid Sui event", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		addr, err := sui.EncodeAddress("0x25db16c3ca555f6702c07860503107bb73cce9f6c1d6df00464529db15d5a5ab")
		require.NoError(t, err)

		value := big.NewInt(1000000)
		suiWithdrawalEvent := sample.MRC20Withdrawal(
			addr,
			value,
		)

		err = k.ValidateMRC20WithdrawEvent(ctx, suiWithdrawalEvent, chains.SuiMainnet.ChainId, coin.CoinType_Gas)
		require.NoError(t, err)
	})
}

func TestKeeper_ProcessMRC20WithdrawalEvent(t *testing.T) {
	t.Run("successfully process MRC20Withdrawal to BTC chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.BitcoinMainnet
		chainID := chain.ChainId
		senderChain := chains.MuseChainMainnet
		setSupportedChain(ctx, zk, []int64{chainID, senderChain.ChainId}...)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		event, err := crosschainkeeper.ParseMRC20WithdrawalEvent(*sample.ValidMRC20WithdrawToBTCReceipt(t).Logs[3])
		require.NoError(t, err)
		mrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "bitcoin", "BTC")
		event.Raw.Address = mrc20
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err = k.ProcessMRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex())
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 1)
		require.Equal(t, "bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu", cctxList[0].GetCurrentOutboundParam().Receiver)
		require.Equal(t, emittingContract.Hex(), cctxList[0].InboundParams.Sender)
		require.Equal(t, txOrigin.Hex(), cctxList[0].InboundParams.TxOrigin)
	})

	t.Run("successfully process MRC20Withdrawal to ETH chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.Ethereum
		chainID := chain.ChainId
		senderChain := chains.MuseChainMainnet
		setSupportedChain(ctx, zk, []int64{chainID, senderChain.ChainId}...)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		event, err := crosschainkeeper.ParseMRC20WithdrawalEvent(*sample.ValidMrc20WithdrawToETHReceipt(t).Logs[11])
		require.NoError(t, err)
		mrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		event.Raw.Address = mrc20
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err = k.ProcessMRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex())
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 1)
		require.Equal(t, "0x5daBFdd153Aaab4a970fD953DcFEEE8BF6Bb946E", cctxList[0].GetCurrentOutboundParam().Receiver)
		require.Equal(t, emittingContract.Hex(), cctxList[0].InboundParams.Sender)
		require.Equal(t, txOrigin.Hex(), cctxList[0].InboundParams.TxOrigin)
	})

	t.Run("unable to process MRC20Withdrawal if foreign coin is not found", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.Ethereum
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		event, err := crosschainkeeper.ParseMRC20WithdrawalEvent(*sample.ValidMrc20WithdrawToETHReceipt(t).Logs[11])
		require.NoError(t, err)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err = k.ProcessMRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex())
		require.ErrorContains(t, err, "cannot find foreign coin with emittingContract address")
		require.Empty(t, k.GetAllCrossChainTx(ctx))
	})

	t.Run("unable to process MRC20Withdrawal if receiver chain is not supported", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.Ethereum
		chainID := chain.ChainId
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		event, err := crosschainkeeper.ParseMRC20WithdrawalEvent(*sample.ValidMrc20WithdrawToETHReceipt(t).Logs[11])
		require.NoError(t, err)
		mrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		event.Raw.Address = mrc20
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err = k.ProcessMRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex())
		require.ErrorContains(t, err, "chain not supported")
		require.Empty(t, k.GetAllCrossChainTx(ctx))
	})

	t.Run("unable to process MRC20Withdrawal if muse chainID is not correctly set", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.Ethereum
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		event, err := crosschainkeeper.ParseMRC20WithdrawalEvent(*sample.ValidMrc20WithdrawToETHReceipt(t).Logs[11])
		require.NoError(t, err)
		mrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		event.Raw.Address = mrc20
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		ctx = ctx.WithChainID("test_21-1")

		err = k.ProcessMRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex())
		require.ErrorContains(t, err, "failed to convert chainID")
		require.Empty(t, k.GetAllCrossChainTx(ctx))
	})

	t.Run("unable to process MRC20Withdrawal if to address is not in correct format", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.Ethereum
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		event, err := crosschainkeeper.ParseMRC20WithdrawalEvent(*sample.ValidMrc20WithdrawToETHReceipt(t).Logs[11])
		require.NoError(t, err)
		mrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		event.Raw.Address = mrc20
		event.To = ethcommon.Address{}.Bytes()
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err = k.ProcessMRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex())
		require.ErrorContains(t, err, "cannot encode address")
		require.Empty(t, k.GetAllCrossChainTx(ctx))
	})

	t.Run("unable to process MRC20Withdrawal if gaslimit not set on mrc20 contract", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.Ethereum
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		event, err := crosschainkeeper.ParseMRC20WithdrawalEvent(*sample.ValidMrc20WithdrawToETHReceipt(t).Logs[11])
		require.NoError(t, err)
		mrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		event.Raw.Address = mrc20
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		fc, _ := zk.FungibleKeeper.GetForeignCoins(ctx, mrc20.Hex())

		fungibleMock.On("GetForeignCoins", mock.Anything, mock.Anything).Return(fc, true)
		fungibleMock.On("QueryGasLimit", mock.Anything, mock.Anything).
			Return(big.NewInt(0), fmt.Errorf("error querying gas limit"))
		err = k.ProcessMRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex())
		require.ErrorContains(t, err, "error querying gas limit")
		require.Empty(t, k.GetAllCrossChainTx(ctx))
	})

	t.Run("unable to process MRC20Withdrawal if gasprice is not set in crosschain keeper", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.Ethereum
		chainID := chain.ChainId
		senderChain := chains.MuseChainMainnet
		setSupportedChain(ctx, zk, []int64{chainID, senderChain.ChainId}...)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		k.RemoveGasPrice(ctx, strconv.FormatInt(chainID, 10))

		event, err := crosschainkeeper.ParseMRC20WithdrawalEvent(*sample.ValidMrc20WithdrawToETHReceipt(t).Logs[11])
		require.NoError(t, err)
		mrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		event.Raw.Address = mrc20
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err = k.ProcessMRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex())
		require.ErrorContains(t, err, "gasprice not found")
		require.Empty(t, k.GetAllCrossChainTx(ctx))
	})

	t.Run("unable to process MRC20Withdrawal if process cctx fails", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.Ethereum
		chainID := chain.ChainId
		senderChain := chains.MuseChainMainnet
		setSupportedChain(ctx, zk, []int64{chainID, senderChain.ChainId}...)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{
			ChainId: chain.ChainId,
			Nonce:   1,
		})

		event, err := crosschainkeeper.ParseMRC20WithdrawalEvent(*sample.ValidMrc20WithdrawToETHReceipt(t).Logs[11])
		require.NoError(t, err)
		mrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		event.Raw.Address = mrc20
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err = k.ProcessMRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex())
		require.ErrorContains(t, err, "nonce mismatch")
		require.Empty(t, k.GetAllCrossChainTx(ctx))
	})
}

func TestKeeper_ParseMuseSentEvent(t *testing.T) {
	t.Run("successfully parse a valid event", func(t *testing.T) {
		logs := sample.ValidMuseSentDestinationExternalReceipt(t).Logs
		for i, log := range logs {
			connector := log.Address
			event, err := crosschainkeeper.ParseMuseSentEvent(*log, connector)
			if i < 4 {
				require.ErrorContains(t, err, "event signature mismatch")
				require.Nil(t, event)
				continue
			}
			require.Equal(t, chains.Ethereum.ChainId, event.DestinationChainId.Int64())
			require.Equal(t, "70000000000000000000", event.MuseValueAndGas.String())
			require.Equal(t, "0x60983881bdf302dcfa96603A58274D15D5966209", event.SourceTxOriginAddress.String())
			require.Equal(t, "0xF0a3F93Ed1B126142E61423F9546bf1323Ff82DF", event.MuseTxSenderAddress.String())
		}
	})

	t.Run("unable to parse if topics field is empty", func(t *testing.T) {
		logs := sample.ValidMuseSentDestinationExternalReceipt(t).Logs
		for _, log := range logs {
			connector := log.Address
			log.Topics = nil
			event, err := crosschainkeeper.ParseMuseSentEvent(*log, connector)
			require.ErrorContains(t, err, "ParseMuseSentEvent: invalid log - no topics")
			require.Nil(t, event)
		}
	})

	t.Run("unable to parse if connector address does not match", func(t *testing.T) {
		logs := sample.ValidMuseSentDestinationExternalReceipt(t).Logs
		for i, log := range logs {
			event, err := crosschainkeeper.ParseMuseSentEvent(*log, sample.EthAddress())
			if i < 4 {
				require.ErrorContains(t, err, "event signature mismatch")
				require.Nil(t, event)
				continue
			}
			require.ErrorContains(t, err, "does not match connectorMEVM")
			require.Nil(t, event)
		}
	})
}
func TestKeeper_ProcessMuseSentEvent(t *testing.T) {
	t.Run("successfully process MuseSentEvent", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)

		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.Ethereum
		chainID := chain.ChainId
		senderChain := chains.MuseChainMainnet
		setSupportedChain(ctx, zk, []int64{chainID, senderChain.ChainId}...)

		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		admin := keepertest.SetAdminPolicies(ctx, zk.AuthorityKeeper)
		SetupStateForProcessLogsMuseSent(t, ctx, k, zk, sdkk, chain, admin)

		amount, ok := sdkmath.NewIntFromString("20000000000000000000000")
		require.True(t, ok)
		err := sdkk.BankKeeper.MintCoins(
			ctx,
			fungibletypes.ModuleName,
			sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)),
		)
		require.NoError(t, err)

		event, err := crosschainkeeper.ParseMuseSentEvent(
			*sample.ValidMuseSentDestinationExternalReceipt(t).Logs[4],
			sample.ValidMuseSentDestinationExternalReceipt(t).Logs[4].Address,
		)
		require.NoError(t, err)
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err = k.ProcessMuseSentEvent(ctx, event, emittingContract, txOrigin.Hex())
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 1)
		require.Equal(
			t,
			strings.Compare(
				"0x60983881bdf302dcfa96603a58274d15d5966209",
				cctxList[0].GetCurrentOutboundParam().Receiver,
			),
			0,
		)
		require.Equal(t, chains.Ethereum.ChainId, cctxList[0].GetCurrentOutboundParam().ReceiverChainId)
		require.Equal(t, emittingContract.Hex(), cctxList[0].InboundParams.Sender)
		require.Equal(t, txOrigin.Hex(), cctxList[0].InboundParams.TxOrigin)
	})

	t.Run("unable to process MuseSentEvent if fungible module does not have enough balance", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.Ethereum
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		admin := keepertest.SetAdminPolicies(ctx, zk.AuthorityKeeper)
		SetupStateForProcessLogsMuseSent(t, ctx, k, zk, sdkk, chain, admin)

		event, err := crosschainkeeper.ParseMuseSentEvent(
			*sample.ValidMuseSentDestinationExternalReceipt(t).Logs[4],
			sample.ValidMuseSentDestinationExternalReceipt(t).Logs[4].Address,
		)
		require.NoError(t, err)
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err = k.ProcessMuseSentEvent(ctx, event, emittingContract, txOrigin.Hex())
		require.ErrorContains(t, err, "ProcessMuseSentEvent: failed to burn coins from fungible")
	})

	t.Run("unable to process MuseSentEvent if receiver chain is not supported", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.Ethereum
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		admin := keepertest.SetAdminPolicies(ctx, zk.AuthorityKeeper)
		SetupStateForProcessLogsMuseSent(t, ctx, k, zk, sdkk, chain, admin)

		amount, ok := sdkmath.NewIntFromString("20000000000000000000000")
		require.True(t, ok)
		err := sdkk.BankKeeper.MintCoins(
			ctx,
			fungibletypes.ModuleName,
			sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)),
		)
		require.NoError(t, err)

		event, err := crosschainkeeper.ParseMuseSentEvent(
			*sample.ValidMuseSentDestinationExternalReceipt(t).Logs[4],
			sample.ValidMuseSentDestinationExternalReceipt(t).Logs[4].Address,
		)
		require.NoError(t, err)
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err = k.ProcessMuseSentEvent(ctx, event, emittingContract, txOrigin.Hex())
		require.ErrorContains(t, err, "chain not supported")
	})

	t.Run("unable to process MuseSentEvent if musechain chain id not correctly set in context", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.Ethereum
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		admin := keepertest.SetAdminPolicies(ctx, zk.AuthorityKeeper)
		SetupStateForProcessLogsMuseSent(t, ctx, k, zk, sdkk, chain, admin)

		amount, ok := sdkmath.NewIntFromString("20000000000000000000000")
		require.True(t, ok)
		err := sdkk.BankKeeper.MintCoins(
			ctx,
			fungibletypes.ModuleName,
			sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)),
		)
		require.NoError(t, err)

		event, err := crosschainkeeper.ParseMuseSentEvent(
			*sample.ValidMuseSentDestinationExternalReceipt(t).Logs[4],
			sample.ValidMuseSentDestinationExternalReceipt(t).Logs[4].Address,
		)
		require.NoError(t, err)
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		ctx = ctx.WithChainID("test-21-1")
		err = k.ProcessMuseSentEvent(ctx, event, emittingContract, txOrigin.Hex())
		require.ErrorContains(t, err, "ProcessMuseSentEvent: failed to convert chainID")
	})

	t.Run("unable to process MuseSentEvent if gas pay fails", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.Ethereum
		chainID := chain.ChainId
		senderChain := chains.MuseChainMainnet
		setSupportedChain(ctx, zk, []int64{chainID, senderChain.ChainId}...)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		amount, ok := sdkmath.NewIntFromString("20000000000000000000000")
		require.True(t, ok)
		err := sdkk.BankKeeper.MintCoins(
			ctx,
			fungibletypes.ModuleName,
			sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)),
		)
		require.NoError(t, err)
		event, err := crosschainkeeper.ParseMuseSentEvent(
			*sample.ValidMuseSentDestinationExternalReceipt(t).Logs[4],
			sample.ValidMuseSentDestinationExternalReceipt(t).Logs[4].Address,
		)
		require.NoError(t, err)
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err = k.ProcessMuseSentEvent(ctx, event, emittingContract, txOrigin.Hex())
		require.ErrorContains(t, err, "gas coin contract invalid address")
	})

	t.Run("unable to process MuseSentEvent if process cctx fails", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.Ethereum
		chainID := chain.ChainId
		senderChain := chains.MuseChainMainnet
		setSupportedChain(ctx, zk, []int64{chainID, senderChain.ChainId}...)

		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		admin := keepertest.SetAdminPolicies(ctx, zk.AuthorityKeeper)
		SetupStateForProcessLogsMuseSent(t, ctx, k, zk, sdkk, chain, admin)

		zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{
			ChainId: chain.ChainId,
			Nonce:   1,
		})
		amount, ok := sdkmath.NewIntFromString("20000000000000000000000")
		require.True(t, ok)
		err := sdkk.BankKeeper.MintCoins(
			ctx,
			fungibletypes.ModuleName,
			sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)),
		)
		require.NoError(t, err)

		event, err := crosschainkeeper.ParseMuseSentEvent(
			*sample.ValidMuseSentDestinationExternalReceipt(t).Logs[4],
			sample.ValidMuseSentDestinationExternalReceipt(t).Logs[4].Address,
		)
		require.NoError(t, err)
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err = k.ProcessMuseSentEvent(ctx, event, emittingContract, txOrigin.Hex())
		require.ErrorContains(t, err, "nonce mismatch")
	})
}

func TestKeeper_ProcessLogs(t *testing.T) {
	t.Run("successfully parse and process MRC20Withdrawal to BTC chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.BitcoinMainnet
		chainID := chain.ChainId
		senderChain := chains.MuseChainMainnet
		setSupportedChain(ctx, zk, []int64{chainID, senderChain.ChainId}...)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.ValidMRC20WithdrawToBTCReceipt(t)
		gasMRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "bitcoin", "BTC")
		for _, log := range block.Logs {
			log.Address = gasMRC20
		}
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err := k.ProcessLogs(ctx, block.Logs, emittingContract, txOrigin.Hex())
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 1)
		require.Equal(t, "bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu", cctxList[0].GetCurrentOutboundParam().Receiver)
		require.Equal(t, emittingContract.Hex(), cctxList[0].InboundParams.Sender)
		require.Equal(t, txOrigin.Hex(), cctxList[0].InboundParams.TxOrigin)
	})

	t.Run("successfully parse and process gateway withdraw to SOL chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.SolanaDevnet
		chainID := chain.ChainId
		senderChain := chains.MuseChainMainnet
		setSupportedChain(ctx, zk, []int64{chainID, senderChain.ChainId}...)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.ValidGatewayWithdrawToSOLChainReceipt(t)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "solana", "SOL")
		txOrigin := sample.EthAddress()

		err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), txOrigin.Hex())
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 1)
		require.Equal(t, "9fA4vYZfCa9k9UHjnvYCk4YoipsooapGciKMgaTBw9UH", cctxList[0].GetCurrentOutboundParam().Receiver)
		require.Equal(t, txOrigin.Hex(), cctxList[0].InboundParams.TxOrigin)
	})

	t.Run("successfully parse and process gateway withdraw and call to SOL chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.SolanaDevnet
		chainID := chain.ChainId
		senderChain := chains.MuseChainMainnet
		setSupportedChain(ctx, zk, []int64{chainID, senderChain.ChainId}...)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.ValidGatewayWithdrawAndCallToSOLChainReceipt(t)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "solana", "SOL")
		txOrigin := sample.EthAddress()

		err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), txOrigin.Hex())
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 1)
		require.Equal(t, "4xEw862A2SEwMjofPkUyd4NEekmVJKJsdHkK3UkAtDrc", cctxList[0].GetCurrentOutboundParam().Receiver)
		require.Equal(t, txOrigin.Hex(), cctxList[0].InboundParams.TxOrigin)
	})

	t.Run("fails to parse and process invalid gateway withdraw to SOL chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.SolanaDevnet
		chainID := chain.ChainId
		senderChain := chains.MuseChainMainnet
		setSupportedChain(ctx, zk, []int64{chainID, senderChain.ChainId}...)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.InvalidGatewayWithdrawToSOLChainReceipt(t)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "solana", "SOL")
		txOrigin := sample.EthAddress()

		err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), txOrigin.Hex())
		require.Error(t, err)
	})

	t.Run("successfully parse and process gateway call to SOL chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.SolanaDevnet
		chainID := chain.ChainId
		senderChain := chains.MuseChainMainnet
		setSupportedChain(ctx, zk, []int64{chainID, senderChain.ChainId}...)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.ValidGatewayCallToSOLChainReceipt(t)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "solana", "SOL")
		txOrigin := sample.EthAddress()

		err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), txOrigin.Hex())
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 1)
		require.Equal(t, "4xEw862A2SEwMjofPkUyd4NEekmVJKJsdHkK3UkAtDrc", cctxList[0].GetCurrentOutboundParam().Receiver)
		require.Zero(t, cctxList[0].GetCurrentOutboundParam().Amount.BigInt().Int64())
		require.Equal(t, txOrigin.Hex(), cctxList[0].InboundParams.TxOrigin)
	})

	t.Run("fails to parse and process invalid gateway call to SOL chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.SolanaDevnet
		chainID := chain.ChainId
		senderChain := chains.MuseChainMainnet
		setSupportedChain(ctx, zk, []int64{chainID, senderChain.ChainId}...)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.InvalidGatewayCallToSOLChainReceipt(t)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "solana", "SOL")
		txOrigin := sample.EthAddress()

		err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), txOrigin.Hex())
		require.Error(t, err)
	})

	t.Run("successfully parse and process MuseSentEvent", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.Ethereum
		chainID := chain.ChainId
		senderChain := chains.MuseChainMainnet
		setSupportedChain(ctx, zk, []int64{chainID, senderChain.ChainId}...)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		admin := keepertest.SetAdminPolicies(ctx, zk.AuthorityKeeper)
		SetupStateForProcessLogsMuseSent(t, ctx, k, zk, sdkk, chain, admin)

		amount, ok := sdkmath.NewIntFromString("20000000000000000000000")
		require.True(t, ok)
		err := sdkk.BankKeeper.MintCoins(
			ctx,
			fungibletypes.ModuleName,
			sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)),
		)
		require.NoError(t, err)
		block := sample.ValidMuseSentDestinationExternalReceipt(t)
		system, found := zk.FungibleKeeper.GetSystemContract(ctx)
		require.True(t, found)
		for _, log := range block.Logs {
			log.Address = ethcommon.HexToAddress(system.ConnectorMevm)
		}
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err = k.ProcessLogs(ctx, block.Logs, emittingContract, txOrigin.Hex())
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 1)
		require.Equal(
			t,
			strings.Compare(
				"0x60983881bdf302dcfa96603a58274d15d5966209",
				cctxList[0].GetCurrentOutboundParam().Receiver,
			),
			0,
		)
		require.Equal(t, chains.Ethereum.ChainId, cctxList[0].GetCurrentOutboundParam().ReceiverChainId)
		require.Equal(t, emittingContract.Hex(), cctxList[0].InboundParams.Sender)
		require.Equal(t, txOrigin.Hex(), cctxList[0].InboundParams.TxOrigin)
	})

	t.Run("unable to process logs if system contract not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		err := k.ProcessLogs(ctx, sample.ValidMRC20WithdrawToBTCReceipt(t).Logs, sample.EthAddress(), "")
		require.ErrorContains(t, err, "cannot find system contract")
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 0)
	})

	t.Run("no cctx created for logs containing no events", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.BitcoinMainnet
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.ValidMRC20WithdrawToBTCReceipt(t)
		gasMRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "bitcoin", "BTC")
		for _, log := range block.Logs {
			log.Address = gasMRC20
		}
		block.Logs = block.Logs[:3]

		err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), "")
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 0)
	})

	t.Run(
		"no cctx created  for logs containing proper event but not emitted from a known MRC20 contract",
		func(t *testing.T) {
			k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
			k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
			chain := chains.BitcoinMainnet
			chainID := chain.ChainId
			setSupportedChain(ctx, zk, chainID)
			SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

			block := sample.ValidMRC20WithdrawToBTCReceipt(t)
			setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "bitcoin", "BTC")
			for _, log := range block.Logs {
				log.Address = sample.EthAddress()
			}

			err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), "")
			require.NoError(t, err)
			cctxList := k.GetAllCrossChainTx(ctx)
			require.Len(t, cctxList, 0)
		},
	)

	t.Run("no cctx created for valid logs if Inbound is disabled", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.BitcoinMainnet
		chainID := chain.ChainId
		senderChain := chains.MuseChainMainnet
		setSupportedChain(ctx, zk, []int64{chainID, senderChain.ChainId}...)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.ValidMRC20WithdrawToBTCReceipt(t)
		gasMRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "bitcoin", "BTC")
		for _, log := range block.Logs {
			log.Address = gasMRC20
		}
		zk.ObserverKeeper.SetCrosschainFlags(ctx, observertypes.CrosschainFlags{
			IsInboundEnabled: false,
		})

		err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), "")
		require.ErrorContains(t, err, observertypes.ErrInboundDisabled.Error())
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 0)
	})

	t.Run("error returned for invalid event data", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		// use the wrong (testnet) chain ID to make the btc address parsing fail
		chain := chains.BitcoinTestnet
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.InvalidMRC20WithdrawToExternalReceipt(t)
		gasMRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "bitcoin", "BTC")
		for _, log := range block.Logs {
			log.Address = gasMRC20
		}

		err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), "")
		require.ErrorContains(t, err, "invalid address")
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 0)
	})

	t.Run("error returned if unable to process an event", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.BitcoinMainnet
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.ValidMRC20WithdrawToBTCReceipt(t)
		gasMRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "bitcoin", "BTC")
		for _, log := range block.Logs {
			log.Address = gasMRC20
		}
		ctx = ctx.WithChainID("test-21-1")

		err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), "")
		require.ErrorContains(t, err, "ProcessMRC20WithdrawalEvent: failed to convert chainID")
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 0)
	})
}
