package blockchain

import (
	mc "github.com/codehakase/mobius-client-go/client"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	keypair "github.com/stellar/go/keypair"
)

// Cosigner represents a cosigner to add to an account
type Cosigner struct{}

func (c *Cosigner) Call(kp, cosignerKeypair *keypair.Full, weight int) (horizon.TransactionSuccess, error) {
	var ts horizon.TransactionSuccess
	client := mc.NewClient().HorizonClient
	account, err := Build(kp)
	// create transaction
	tx, err := c.tx(account, cosignerKeypair, weight)
	if err != nil {
		return ts, err
	}
	// sign transaciton
	akp := account.Keypair.(*keypair.Full)
	txe, err := tx.Sign(akp.Seed())
	if err != nil {
		return ts, err
	}
	txtEnvelopeStr, err := txe.Base64()
	if err != nil {
		return ts, err
	}
	return client.SubmitTransaction(txtEnvelopeStr)
}

func (c *Cosigner) tx(account *Account, cosignerKeypair *keypair.Full, weight int) (*build.TransactionBuilder, error) {
	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: account.Keypair.Address()},
		build.SetThresholds(1, 1, 10),
		build.MasterWeight(10),
		build.AddSigner(cosignerKeypair.Address(), uint32(weight)),
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
