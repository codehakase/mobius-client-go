package utils

import (
	"log"
	"os"

	"github.com/codehakase/mobius-client-go/utils/custom/transaction"
	"github.com/stellar/go/keypair"
)

func KPFromSeed(seed string) *keypair.Full {
	kp, err := keypair.Parse(seed)
	if err != nil {
		log.Fatalf("failed to create keypair from developer secret key, err: %v", err)
		os.Exit(1)
	}
	switch kp.(type) {
	case *keypair.Full:
		return kp.(*keypair.Full)
	case *keypair.FromAddress:
		return new(keypair.Full)
	}
	log.Fatalf("failed to create keypair from developer secret key, err: %v", err)
	os.Exit(1)
	return nil
}

func KPFromAddress(address string) keypair.KP {
	if address == "" {
		log.Fatalf("failed to create keypair from user's public address")
		os.Exit(1)
	}
	kp, err := keypair.Parse(address)
	if err != nil {
		log.Fatalf("failed to create keypair from user's public address, err: %v", err)
		os.Exit(1)
	}
	return kp
}

func ValidateTx(devKeypair keypair.KP, tx *transaction.Transaction) bool {
	if tx.Tx.Signatures == nil || len(tx.Tx.Signatures) < 1 {
		return false
	}
	hash := tx.Hash()
	for _, signature := range tx.Tx.Signatures {
		if err := devKeypair.Verify(hash[:], signature.Signature); err != nil {
			return true
		}
	}
	log.Fatalf("wrong challenge transaction signature")
	return false
}
