package blockchain

import (
	"log"
	"strconv"

	mc "github.com/codehakase/mobius-client-go/client"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	keypair "github.com/stellar/go/keypair"
	hp "github.com/stellar/go/protocols/horizon"
)

// Account represents an interface to interact with an account on the Stellar
// Network
type Account struct {
	Account        horizon.Account
	Keypair        keypair.KP
	AssetIssuers   []string
	ClientInstance *horizon.Client
}

// Build fetches information from the Stellar Network and returns an instance of
// Account
func Build(kp keypair.KP) (*Account, error) {
	accountID := kp.Address()
	account, err := mc.NewClient().HorizonClient.LoadAccount(accountID)
	if err != nil {
		return nil, err
	}
	return &Account{
		account,
		kp,
		[]string{},
		mc.NewClient().HorizonClient,
	}, nil
}

// GetKeypair returns the keypair for account
func (a *Account) GetKeypair() keypair.KP { return a.Keypair }

// Info returns the Account information
func (a *Account) Info() horizon.Account { return a.Account }

// Authorized confirms if the given keypair is added as a cosigner to account
func (a *Account) Authorized(kp keypair.KP) bool {
	signer := a.findSigner(kp.Address())
	if (hp.Signer{}) == signer {
		return false
	}
	return true
}

// Balance returns the balance for given asset
func (a *Account) Balance(asset build.Asset) float64 {
	balance, _ := strconv.ParseFloat(a.findBalance(asset).Balance, 64)
	return balance
}

// TrustLineExists returns true if a trustline exist for given asset, and limit is possitive
func (a *Account) TrustLineExists(asset build.Asset) bool {
	limit, _ := strconv.ParseFloat(a.findBalance(asset).Limit, 64)
	return (limit > 0)
}

// Reload invalidates current account information
func (a *Account) Reload() *Account {
	account, err := Build(a.Keypair)
	if err != nil {
		log.Fatalf("failed to reload account: %v", err)
	}
	return account
}

// findSigner returns a matched signer with public key
func (a *Account) findSigner(publicKey string) hp.Signer {
	var signer hp.Signer
	for _, s := range a.Account.Signers {
		if s.PublicKey == publicKey {
			signer = s
		}
	}
	return signer
}

// findBalance returns balance mathcing asset
func (a *Account) findBalance(asset build.Asset) hp.Balance {
	var balance hp.Balance
	for _, b := range a.Account.Balances {
		if a.balanceMatches(asset, b) {
			balance = b
		}
	}
	return balance
}

// balanceMatches compares asset, and confirms balance matches given asset
func (a *Account) balanceMatches(asset build.Asset, balance hp.Balance) bool {
	if asset.Native {
		return balance.Type == "native"
	} else {
		return (balance.Code == asset.Code && balance.Issuer == asset.Issuer)
	}
}
