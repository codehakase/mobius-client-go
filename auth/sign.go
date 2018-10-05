package auth

import (
	"log"

	"github.com/codehakase/mobius-client-go/utils"
	"github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
)

// Sign the user's challenge transaction
type Sign struct{}

// Sign adds a signature to the given transaction
func (s *Sign) Call(userSecret, xdrs, address string) string {
	tx := &xdr.TransactionEnvelope{}
	err := tx.Scan(xdrs)
	if err != nil {
		log.Fatalf("failed to scan challenge transaction xdr, err: %v", err)
	}
	kp := utils.KPFromSeed(userSecret)
	devKeypair := utils.KPFromAddress(address)
	if err != nil {
		log.Fatalf("failed to parse dev key pair, err: %v", err)
	}
	txt := &build.TransactionBuilder{TX: &tx.Tx}
	_ = utils.ValidateTx(devKeypair, tx, txt)
	txe, err := txt.Sign(kp.Seed())
	if err != nil {
		log.Fatalf("failed to sign challenge transaction, err: %v", err)
	}
	txtEnvelopeStr, err := txe.Base64()
	if err != nil {
		log.Fatalf("failed to generate txt envelope str, err: %v", err)
	}
	return txtEnvelopeStr
}
