package ton

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Faucet represents the faucet information.
//
//nolint:revive,stylecheck // comes from my-local-ton
type Faucet struct {
	InitialBalance   int64  `json:"initialBalance"`
	PrivateKey       string `json:"privateKey"`
	PublicKey        string `json:"publicKey"`
	WalletRawAddress string `json:"walletRawAddress"`
	Mnemonic         string `json:"mnemonic"`
	WalletVersion    string `json:"walletVersion"`
	WorkChain        int32  `json:"workChain"`
	SubWalletId      int    `json:"subWalletId"`
	Created          bool   `json:"created"`
}

// GetFaucet returns the faucet information.
func GetFaucet(ctx context.Context, url string) (Faucet, error) {
	resp, err := get(ctx, url)
	if err != nil {
		return Faucet{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Faucet{}, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	var faucet Faucet
	if err := json.NewDecoder(resp.Body).Decode(&faucet); err != nil {
		return Faucet{}, err
	}

	return faucet, nil
}

func get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}
