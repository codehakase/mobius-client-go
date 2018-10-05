package utils

import (
	"crypto/sha256"
	"log"
	"os"

	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

func KPFromSeed(seed string) *keypair.Full {
	keypairHash := sha256.Sum256([]byte(seed))
	kp, err := keypair.FromRawSeed(keypairHash)
	if err != nil {
		log.Fatalf("failed to create keypair from developer secret key, err: %v", err)
		os.Exit(1)
	}
	return kp
}

func KPFromAddress(address string) *keypair.Full {
	if address == "" {
		log.Fatalf("failed to create keypair from user's public address")
		os.Exit(1)
	}
	kp, err := keypair.Parse(address)
	if err != nil {
		log.Fatalf("failed to create keypair from user's public address, err: %v", err)
		os.Exit(1)
	}
	return kp.(*keypair.Full)
}

func ValidateTx(devKeypair *keypair.Full, tx *xdr.TransactionEnvelope, t *build.TransactionBuilder) bool {
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
