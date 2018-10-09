package auth

import (
	"bytes"
	"encoding/base64"
	"log"

	"github.com/codehakase/mobius-client-go/utils"
	"github.com/codehakase/mobius-client-go/utils/custom/transaction"
	"github.com/stellar/go/xdr"
)

// Sign the user's challenge transaction
type Sign struct{}

// Sign adds a signature to the given transaction
func (s *Sign) Call(userSecret, xdrs, address string) string {
	tx, err := transaction.New(xdrs)
	if err != nil {
		log.Fatal(err)
	}
	kp := utils.KPFromSeed(userSecret)
	devKeypair := utils.KPFromAddress(address)
	if err != nil {
		log.Fatalf("failed to parse dev key pair, err: %v", err)
	}
	_ = utils.ValidateTx(devKeypair, tx)
	tx.Sign(kp)
	var txBytes bytes.Buffer
	_, err = xdr.Marshal(&txBytes, tx.Tx)
	if err != nil {
		log.Fatalf("marshal xdr failed, %v", err)
	}
	return base64.StdEncoding.EncodeToString(txBytes.Bytes())
}
