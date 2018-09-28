package auth

import (
	"log"

	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
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
	kp := kpFromSeed(userSecret)
	devKeypair, err := keypair.Parse(address)
	if err != nil {
		log.Fatalf("failed to parse dev key pair, err: %v", err)
	}
	txt := build.TransactionBuilder{TX: &tx.Tx}
	_ = s.validate(devKeypair, tx, txt)
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

func (s *Sign) validate(devKeypair keypair.KP, tx *xdr.TransactionEnvelope, t build.TransactionBuilder) bool {
	if tx.Signatures == nil || len(tx.Signatures) < 1 {
		return false
	}
	hash, err := t.Hash()
	if err != nil {
		log.Fatalf("failed to retrieve transaction hash, err: %v", err)
	}
	for _, signature := range tx.Signatures {
		if err := devKeypair.Verify(hash[:], signature.Signature); err != nil {
			return true
		}
	}
	log.Fatalf("wrong challenge transaction signature")
	return false
}
