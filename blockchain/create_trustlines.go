package blockchain

import (
	mc "github.com/codehakase/mobius-client-go/client"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	keypair "github.com/stellar/go/keypair"
)

// CreateTrustline reprents a model for creating multiple trustline for a given asset
type CreateTrustline struct{}

func (ct *CreateTrustline) Call(kp keypair.KP, asset build.Asset) (horizon.TransactionSuccess, error) {
	var ts horizon.TransactionSuccess
	client := mc.NewClient().HorizonClient
	account, err := Build(kp)
	// create transaction
	tx, err := ct.tx(account, asset)
	if err != nil {
		return ts, err
	}
	// sign transaciton
	txe, err := tx.Sign(account.Keypair.Address())
	if err != nil {
		return ts, err
	}
	txtEnvelopeStr, err := txe.Base64()
	if err != nil {
		return ts, err
	}
	return client.SubmitTransaction(txtEnvelopeStr)
}

func (ct *CreateTrustline) tx(account *Account, asset build.Asset) (*build.TransactionBuilder, error) {
	tx, err := build.Transaction(
		build.ChangeTrust(asset),
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
