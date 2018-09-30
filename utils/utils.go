package utils

import (
	"crypto/sha256"
	"log"
	"os"

	"github.com/stellar/go/keypair"
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
