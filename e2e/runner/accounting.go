package runner

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	musecrypto "github.com/RWAs-labs/muse/pkg/crypto"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

const (
	// MRC20InitialSupply is the initial supply of the MRC20 token
	MRC20SOLInitialSupply = 100000000

	// SolanaPDAInitialBalance is the initial balance (in lamports) of the gateway PDA account
	SolanaPDAInitialBalance = 1447680
)

type Amount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Response struct {
	Amount Amount `json:"amount"`
}

func (r *E2ERunner) CheckMRC20BalanceAndSupply() {
	r.Logger.Info("Checking MRC20 Balance vs. Supply")

	err := r.checkETHTSSBalance()
	require.NoError(r, err, "ETH balance check failed")

	err = r.checkERC20TSSBalance()
	require.NoError(r, err, "ERC20 balance check failed")

	err = r.checkMuseTSSBalance()
	require.NoError(r, err, "MUSE balance check failed")

	err = r.CheckBTCTSSBalance()
	require.NoError(r, err, "BTC balance check failed")
}

func (r *E2ERunner) checkETHTSSBalance() error {
	allTssAddress, err := r.ObserverClient.TssHistory(r.Ctx, &observertypes.QueryTssHistoryRequest{})
	if err != nil {
		return err
	}

	tssTotalBalance := big.NewInt(0)

	for _, tssAddress := range allTssAddress.TssList {
		evmAddress, err := r.ObserverClient.GetTssAddressByFinalizedHeight(
			r.Ctx,
			&observertypes.QueryGetTssAddressByFinalizedHeightRequest{
				FinalizedMuseHeight: tssAddress.FinalizedMuseHeight,
			},
		)
		if err != nil {
			continue
		}

		tssBal, err := r.EVMClient.BalanceAt(r.Ctx, common.HexToAddress(evmAddress.Eth), nil)
		if err != nil {
			continue
		}
		tssTotalBalance.Add(tssTotalBalance, tssBal)
	}

	mrc20Supply, err := r.ETHMRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return err
	}
	if tssTotalBalance.Cmp(mrc20Supply) < 0 {
		return fmt.Errorf("ETH: TSS balance (%d) < MRC20 TotalSupply (%d) ", tssTotalBalance, mrc20Supply)
	}
	r.Logger.Info("ETH: TSS balance (%d) >= MRC20 TotalSupply (%d)", tssTotalBalance, mrc20Supply)
	return nil
}

func (r *E2ERunner) CheckBTCTSSBalance() error {
	allTssAddress, err := r.ObserverClient.TssHistory(r.Ctx, &observertypes.QueryTssHistoryRequest{})
	if err != nil {
		return err
	}

	tssTotalBalance := float64(0)

	for _, tssAddress := range allTssAddress.TssList {
		btcTssAddress, err := musecrypto.GetTSSAddrBTC(tssAddress.TssPubkey, r.BitcoinParams)
		if err != nil {
			continue
		}
		utxos, err := r.BtcRPCClient.ListUnspent(r.Ctx)
		if err != nil {
			continue
		}
		for _, utxo := range utxos {
			if utxo.Address == btcTssAddress {
				tssTotalBalance += utxo.Amount
			}
		}
	}

	mrc20Supply, err := r.BTCMRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return err
	}

	// check the balance in TSS is greater than the total supply on MuseChain
	// the amount minted to initialize the pool is subtracted from the total supply
	// #nosec G701 test - always in range
	if int64(tssTotalBalance*1e8) < (mrc20Supply.Int64() - 10000000) {
		// #nosec G701 test - always in range
		return fmt.Errorf(
			"BTC: TSS Balance (%d) < MRC20 TotalSupply (%d)",
			int64(tssTotalBalance*1e8),
			mrc20Supply.Int64()-10000000,
		)
	}
	// #nosec G115 test - always in range
	r.Logger.Info(
		"BTC: Balance (%d) >= MRC20 TotalSupply (%d)",
		int64(tssTotalBalance*1e8),
		mrc20Supply.Int64()-10000000,
	)

	return nil
}

// CheckSolanaTSSBalance compares the gateway PDA balance with the total supply of the SOL MRC20 on MuseChain
func (r *E2ERunner) CheckSolanaTSSBalance() error {
	mrc20Supply, err := r.SOLMRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return err
	}

	// get PDA received amount
	pda := r.ComputePdaAddress()
	balance, err := r.SolanaClient.GetBalance(r.Ctx, pda, rpc.CommitmentConfirmed)
	require.NoError(r, err)
	pdaReceivedAmount := balance.Value - SolanaPDAInitialBalance

	// the SOL balance in gateway PDA must not be less than the total supply on MuseChain
	// the amount minted to initialize the pool is subtracted from the total supply
	// #nosec G115 test - always in range
	if pdaReceivedAmount < (mrc20Supply.Uint64() - MRC20SOLInitialSupply) {
		// #nosec G115 test - always in range
		return fmt.Errorf(
			"SOL: Gateway PDA Received (%d) < MRC20 TotalSupply (%d)",
			pdaReceivedAmount,
			mrc20Supply.Uint64()-MRC20SOLInitialSupply,
		)
	}
	// #nosec G115 test - always in range
	r.Logger.Info(
		"SOL: Gateway PDA Received (%d) >= MRC20 TotalSupply (%d)",
		pdaReceivedAmount,
		mrc20Supply.Int64()-MRC20SOLInitialSupply,
	)

	return nil
}

func (r *E2ERunner) checkERC20TSSBalance() error {
	custodyBalance, err := r.ERC20.BalanceOf(&bind.CallOpts{}, r.ERC20CustodyAddr)
	if err != nil {
		return err
	}

	erc20mrc20Supply, err := r.ERC20MRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return err
	}
	if custodyBalance.Cmp(erc20mrc20Supply) < 0 {
		return fmt.Errorf("ERC20: TSS balance (%d) < MRC20 TotalSupply (%d) ", custodyBalance, erc20mrc20Supply)
	}
	r.Logger.Info("ERC20: TSS balance (%d) >= ERC20 MRC20 TotalSupply (%d)", custodyBalance, erc20mrc20Supply)
	return nil
}

func (r *E2ERunner) checkMuseTSSBalance() error {
	museLocked, err := r.ConnectorEth.GetLockedAmount(&bind.CallOpts{})
	if err != nil {
		return err
	}
	resp, err := http.Get("http://musecore0:1317/cosmos/bank/v1beta1/supply/by_denom?denom=amuse")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	museSupply, _ := big.NewInt(0).SetString(result.Amount.Amount, 10)
	if museLocked.Cmp(museSupply) < 0 {
		r.Logger.Info("MUSE: TSS balance (%d) < MRC20 TotalSupply (%d)", museLocked, museSupply)
	} else {
		r.Logger.Info("MUSE: TSS balance (%d) >= MRC20 TotalSupply (%d)", museLocked, museSupply)
	}
	return nil
}
