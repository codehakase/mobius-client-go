package auth

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	mc "github.com/codehakase/mobius-client-go/client"
	"github.com/codehakase/mobius-client-go/utils"
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

// Token checks the challenge transaction signed by user on developer's side
type Token struct {
	DeveloperSecret string
	TX              *xdr.TransactionEnvelope
	TxProp          *build.TransactionBuilder
	Address         string // user public key
	devKeypair      *keypair.Full
	userKeypair     *keypair.Full
}

// NewToken creates a new token handler with provided developer seed, xdr and user
// address
func NewToken(devSecret, xdrs, address string) (*Token, error) {
	tx := &xdr.TransactionEnvelope{}
	err := tx.Scan(xdrs)
	if err != nil {
		return nil, fmt.Errorf("failed to scan challenge transaction xdr, err: %v", err)
	}
	txp := &build.TransactionBuilder{TX: &tx.Tx}
	return &Token{
		devSecret,
		tx,
		txp,
		address,
		nil,
		nil,
	}, nil
}

// TimeBounds returns the time bounds for a given transaction
func (t *Token) TimeBounds() *xdr.TimeBounds {
	if t.TX.Tx.TimeBounds == nil {
		log.Fatalf("wrong challenge transaction structure")
	}
	return t.TX.Tx.TimeBounds
}

// Validate transaction signed by developer and user
func (t *Token) Validate(strict bool) bool {
	if !t.signedCorrectly() {
		log.Fatalf("wrong challenge transaction signature")
		return false
	}
	if !t.timeNowCovers(t.TX.Tx.TimeBounds) {
		log.Fatalf("challenge transaction expired")
		return false
	}
	if strict && t.tooOld(t.TX.Tx.TimeBounds) {
		log.Fatalf("challenge transaction expired")
		return false
	}
	return true
}

// Hash validates and returns the transaction hash
func (t *Token) Hash(format string) interface{} {
	_ = t.Validate(true)
	if format == "binary" {
		txHash, _ := t.TxProp.Hash()
		return txHash[:]
	}
	tHexHash, _ := t.TxProp.HashHex()
	return tHexHash
}

// Address returns the address a current token is issued for
func (t *Token) GetAddress() string {
	return t.GetKeypair().Address()
}

// GetKeypair returns the keypair for developer
func (t *Token) GetKeypair() *keypair.Full {
	if t.devKeypair == nil {
		t.devKeypair = utils.KPFromSeed(t.DeveloperSecret)
	}
	return t.devKeypair
}

func (t *Token) GetUserKeypair() *keypair.Full {
	if t.userKeypair == nil {
		t.userKeypair = utils.KPFromAddress(t.Address)
	}
	return t.devKeypair
}

// signedCorrectly confirms is the transaction is correctly signed by user and
// developer
func (t *Token) signedCorrectly() bool {
	txt := build.TransactionBuilder{TX: &t.TX.Tx}
	isSignedByDev := t.validate(t.GetKeypair(), t.TX, txt)
	isSignedByUser := t.validate(t.GetUserKeypair(), t.TX, txt)
	return isSignedByDev && isSignedByUser
}

func (t *Token) validate(kp *keypair.Full, tx *xdr.TransactionEnvelope, txt build.TransactionBuilder) bool {
	if tx.Signatures == nil || len(tx.Signatures) < 1 {
		return false
	}
	hash, err := txt.Hash()
	if err != nil {
		log.Fatalf("failed to retrieve transaction hash, err: %v", err)
	}
	for _, signature := range tx.Signatures {
		if err := kp.Verify(hash[:], signature.Signature); err != nil {
			return true
		}
	}
	return false
}

// timeNowCovers returns true if current tie is within transaction time bounds
func (t *Token) timeNowCovers(timeBounds *xdr.TimeBounds) bool {
	now := math.Floor(float64(time.Now().UnixNano() / 1000))
	nMinTime, _ := strconv.ParseInt(string(timeBounds.MinTime), 10, 32)
	nMaxTime, _ := strconv.ParseInt(string(timeBounds.MaxTime), 10, 32)
	return (now >= float64(nMinTime) && now <= float64(nMaxTime))
}

// tooOld returns true if transaction is created more than 10 seconds from now
func (t *Token) tooOld(timeBounds *xdr.TimeBounds) bool {
	now := math.Floor(float64(time.Now().UnixNano() / 1000))
	nMinTime, _ := strconv.ParseInt((string(timeBounds.MinTime)), 10, 32)
	return (now > float64(nMinTime+mc.StrictInterval))
}
