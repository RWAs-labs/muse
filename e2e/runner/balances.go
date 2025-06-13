package runner

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
)

// AccountBalances is a struct that contains the balances of the accounts used in the E2E test
type AccountBalances struct {
	MuseETH      *big.Int
	MuseMUSE     *big.Int
	MuseWMUSE    *big.Int
	MuseERC20    *big.Int
	MuseBTC      *big.Int
	MuseSOL      *big.Int
	MuseSPL      *big.Int
	MuseSui      *big.Int
	MuseSuiToken *big.Int
	MuseTON      *big.Int
	EvmETH       *big.Int
	EvmMUSE      *big.Int
	EvmERC20     *big.Int
	BtcBTC       string
	SolSOL       *big.Int
	SolSPL       *big.Int
	SuiSUI       uint64
	SuiToken     uint64
	TONTON       uint64
}

// AccountBalancesDiff is a struct that contains the difference in the balances of the accounts used in the E2E test
type AccountBalancesDiff struct {
	ETH   *big.Int
	MUSE  *big.Int
	ERC20 *big.Int
}

type ERC20BalanceOf interface {
	BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error)
}

func (r *E2ERunner) getERC20BalanceSafe(z ERC20BalanceOf, name string) *big.Int {
	// have to use reflect to check nil interface because go'ism
	if z == nil || reflect.ValueOf(z).IsNil() {
		r.Logger.Print("❓ balance of %s: nil", name)
		return new(big.Int)
	}
	res, err := z.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	if err != nil {
		r.Logger.Print("❓ balance of %s: %v", name, err)
		return new(big.Int)
	}
	return res
}

// GetAccountBalances returns the account balances of the accounts used in the E2E test
func (r *E2ERunner) GetAccountBalances(skipBTC bool) (AccountBalances, error) {
	// mevm
	museMuse, err := r.MEVMClient.BalanceAt(r.Ctx, r.EVMAddress(), nil)
	if err != nil {
		return AccountBalances{}, fmt.Errorf("get muse balance: %w", err)
	}
	museWMuse, err := r.WMuse.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	if err != nil {
		return AccountBalances{}, fmt.Errorf("get wmuse balance: %w", err)
	}
	museEth := r.getERC20BalanceSafe(r.ETHMRC20, "eth mrc20")
	museErc20 := r.getERC20BalanceSafe(r.ERC20MRC20, "erc20 mrc20")
	museBtc := r.getERC20BalanceSafe(r.BTCMRC20, "btc mrc20")
	museSol := r.getERC20BalanceSafe(r.SOLMRC20, "sol mrc20")
	museSPL := r.getERC20BalanceSafe(r.SPLMRC20, "spl mrc20")
	museSui := r.getERC20BalanceSafe(r.SUIMRC20, "sui mrc20")
	museSuiToken := r.getERC20BalanceSafe(r.SuiTokenMRC20, "sui token mrc20")
	museTon := r.getERC20BalanceSafe(r.TONMRC20, "ton mrc20")

	// evm
	evmEth, err := r.EVMClient.BalanceAt(r.Ctx, r.EVMAddress(), nil)
	if err != nil {
		return AccountBalances{}, fmt.Errorf("get eth balance: %w", err)
	}
	evmMuse := r.getERC20BalanceSafe(r.MuseEth, "muse eth")
	evmErc20 := r.getERC20BalanceSafe(r.ERC20, "eth erc20")

	// bitcoin
	var BtcBTC string
	if !skipBTC {
		if BtcBTC, err = r.GetBitcoinBalance(); err != nil {
			return AccountBalances{}, err
		}
	}

	// solana
	var solSOL *big.Int
	var solSPL *big.Int
	if r.Account.SolanaAddress != "" && r.Account.SolanaPrivateKey != "" && r.SolanaClient != nil {
		solanaAddr := solana.MustPublicKeyFromBase58(r.Account.SolanaAddress.String())
		privateKey := solana.MustPrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
		solSOLBalance, err := r.SolanaClient.GetBalance(
			r.Ctx,
			solanaAddr,
			rpc.CommitmentConfirmed,
		)
		if err != nil {
			return AccountBalances{}, fmt.Errorf("get sol balance: %w", err)
		}

		// #nosec G115 always in range
		solSOL = big.NewInt(int64(solSOLBalance.Value))

		if r.SPLAddr != (solana.PublicKey{}) {
			ata := r.ResolveSolanaATA(
				privateKey,
				solanaAddr,
				r.SPLAddr,
			)
			splBalance, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, ata, rpc.CommitmentConfirmed)
			if err != nil {
				return AccountBalances{}, fmt.Errorf("get spl balance: %w", err)
			}

			solSPLParsed, ok := new(big.Int).SetString(splBalance.Value.Amount, 10)
			if !ok {
				return AccountBalances{}, errors.New("can't parse spl balance")
			}

			solSPL = solSPLParsed
		}
	}

	// sui
	var suiSUI uint64
	var suiToken uint64
	if r.Clients.Sui != nil {
		signer, err := r.Account.SuiSigner()
		if err != nil {
			return AccountBalances{}, err
		}
		suiSUI = r.SuiGetSUIBalance(signer.Address())
		suiToken = r.SuiGetFungibleTokenBalance(signer.Address())
	}

	// TON
	var tonTON uint64
	if r.Clients.TON != nil {
		_, tonWallet, err := r.Account.AsTONWallet(r.Clients.TON)
		if err == nil {
			tonBalance, err := tonWallet.GetBalance(r.Ctx)
			if err == nil {
				tonTON = tonBalance
			}
		}
	}

	return AccountBalances{
		MuseETH:      museEth,
		MuseMUSE:     museMuse,
		MuseWMUSE:    museWMuse,
		MuseERC20:    museErc20,
		MuseBTC:      museBtc,
		MuseSOL:      museSol,
		MuseSPL:      museSPL,
		MuseSui:      museSui,
		MuseSuiToken: museSuiToken,
		MuseTON:      museTon,
		EvmETH:       evmEth,
		EvmMUSE:      evmMuse,
		EvmERC20:     evmErc20,
		BtcBTC:       BtcBTC,
		SolSOL:       solSOL,
		SolSPL:       solSPL,
		SuiSUI:       suiSUI,
		SuiToken:     suiToken,
		TONTON:       tonTON,
	}, nil
}

// GetBitcoinBalance returns the spendable BTC balance of the BTC address
func (r *E2ERunner) GetBitcoinBalance() (string, error) {
	address, _ := r.GetBtcKeypair()
	total, err := r.GetBitcoinBalanceByAddress(address)
	if err != nil {
		return "", err
	}

	return total.String(), nil
}

// GetBitcoinBalanceByAddress get btc balance by address.
func (r *E2ERunner) GetBitcoinBalanceByAddress(address btcutil.Address) (btcutil.Amount, error) {
	unspentList, err := r.BtcRPCClient.ListUnspentMinMaxAddresses(r.Ctx, 1, 9999999, []btcutil.Address{address})
	if err != nil {
		return 0, errors.Wrap(err, "failed to list unspent")
	}

	var total btcutil.Amount
	for _, unspent := range unspentList {
		total += btcutil.Amount(unspent.Amount * 1e8)
	}

	return total, nil
}

// PrintAccountBalances shows the account balances of the accounts used in the E2E test
// Note: USDT is mentioned as erc20 here because we want to show the balance of any erc20 contract
func (r *E2ERunner) PrintAccountBalances(balances AccountBalances) {
	r.Logger.Print(" ---💰 Account info ---")

	r.Logger.Print("Addresses:")
	r.Logger.Print("* EVM: %s", r.EVMAddress().Hex())
	r.Logger.Print("* Solana: %s", r.SolanaDeployerAddress.String())
	signer, err := r.Account.SuiSigner()
	if err != nil {
		r.Logger.Print("Error getting Sui address: %s", err.Error())
	} else {
		r.Logger.Print("* SUI: %s", signer.Address())
	}

	// mevm
	r.Logger.Print("MuseChain:")
	r.Logger.Print("* MUSE balance:  %s", balances.MuseMUSE.String())
	r.Logger.Print("* WMUSE balance: %s", balances.MuseWMUSE.String())
	r.Logger.Print("* ETH balance:   %s", balances.MuseETH.String())
	r.Logger.Print("* ERC20 balance: %s", balances.MuseERC20.String())
	r.Logger.Print("* BTC balance:   %s", balances.MuseBTC.String())
	r.Logger.Print("* SOL balance: %s", balances.MuseSOL.String())
	r.Logger.Print("* SPL balance: %s", balances.MuseSPL.String())
	r.Logger.Print("* SUI balance: %s", balances.MuseSui.String())
	r.Logger.Print("* SUI Token balance: %s", balances.MuseSuiToken.String())
	r.Logger.Print("* TON balance: %s", balances.MuseTON.String())

	// evm
	r.Logger.Print("EVM:")
	r.Logger.Print("* MUSE balance:  %s", balances.EvmMUSE.String())
	r.Logger.Print("* ETH balance:   %s", balances.EvmETH.String())
	r.Logger.Print("* ERC20 balance: %s", balances.EvmERC20.String())

	// bitcoin
	r.Logger.Print("Bitcoin:")
	r.Logger.Print("* BTC balance: %s", balances.BtcBTC)

	// solana
	r.Logger.Print("Solana:")
	if balances.SolSOL != nil {
		r.Logger.Print("* SOL balance: %s", balances.SolSOL.String())
	}
	if balances.SolSPL != nil {
		r.Logger.Print("* SPL balance: %s", balances.SolSPL.String())
	}

	// sui
	r.Logger.Print("Sui:")
	r.Logger.Print("* SUI balance: %d", balances.SuiSUI)
	r.Logger.Print("* SUI Token balance: %d", balances.SuiToken)

	// TON
	r.Logger.Print("TON:")
	if balances.TONTON != 0 {
		r.Logger.Print("* TON balance: %d", balances.TONTON)
	}
}

// PrintTotalDiff shows the difference in the account balances of the accounts used in the e2e test from two balances structs
func (r *E2ERunner) PrintTotalDiff(diffs AccountBalancesDiff) {
	r.Logger.Print(" ---💰 Total gas spent ---")

	// show the value only if it is not zero
	if diffs.MUSE.Cmp(big.NewInt(0)) != 0 {
		r.Logger.Print("* MUSE spent:  %s", diffs.MUSE.String())
	}
	if diffs.ETH.Cmp(big.NewInt(0)) != 0 {
		r.Logger.Print("* ETH spent:   %s", diffs.ETH.String())
	}
	if diffs.ERC20.Cmp(big.NewInt(0)) != 0 {
		r.Logger.Print("* ERC20 spent: %s", diffs.ERC20.String())
	}
}

// GetAccountBalancesDiff returns the difference in the account balances of the accounts used in the E2E test
func GetAccountBalancesDiff(balancesBefore, balancesAfter AccountBalances) AccountBalancesDiff {
	balancesBeforeMuse := big.NewInt(0).Add(balancesBefore.MuseMUSE, balancesBefore.EvmMUSE)
	balancesBeforeEth := big.NewInt(0).Add(balancesBefore.MuseETH, balancesBefore.EvmETH)
	balancesBeforeErc20 := big.NewInt(0).Add(balancesBefore.MuseERC20, balancesBefore.EvmERC20)

	balancesAfterMuse := big.NewInt(0).Add(balancesAfter.MuseMUSE, balancesAfter.EvmMUSE)
	balancesAfterEth := big.NewInt(0).Add(balancesAfter.MuseETH, balancesAfter.EvmETH)
	balancesAfterErc20 := big.NewInt(0).Add(balancesAfter.MuseERC20, balancesAfter.EvmERC20)

	diffMuse := big.NewInt(0).Sub(balancesBeforeMuse, balancesAfterMuse)
	diffEth := big.NewInt(0).Sub(balancesBeforeEth, balancesAfterEth)
	diffErc20 := big.NewInt(0).Sub(balancesBeforeErc20, balancesAfterErc20)

	return AccountBalancesDiff{
		ETH:   diffEth,
		MUSE:  diffMuse,
		ERC20: diffErc20,
	}
}

// formatBalances formats the AccountBalancesDiff into a one-liner string
func formatBalances(balances AccountBalancesDiff) string {
	parts := []string{}
	if balances.ETH != nil && balances.ETH.Cmp(big.NewInt(0)) > 0 {
		parts = append(parts, fmt.Sprintf("ETH:%s", balances.ETH.String()))
	}
	if balances.MUSE != nil && balances.MUSE.Cmp(big.NewInt(0)) > 0 {
		parts = append(parts, fmt.Sprintf("MUSE:%s", balances.MUSE.String()))
	}
	if balances.ERC20 != nil && balances.ERC20.Cmp(big.NewInt(0)) > 0 {
		parts = append(parts, fmt.Sprintf("ERC20:%s", balances.ERC20.String()))
	}
	return strings.Join(parts, ",")
}
