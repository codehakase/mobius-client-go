package auth

import (
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"strconv"

	mc "github.com/codehakase/mobius-client-go/client"
	"github.com/codehakase/mobius-client-go/utils"
	"github.com/codehakase/mobius-client-go/utils/custom/transaction"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

// Token checks the challenge transaction signed by user on developer's side
type Token struct {
	DeveloperSecret string
	TX              *transaction.Transaction
	Address         string // user public key
	devKeypair      *keypair.Full
	userKeypair     keypair.KP
}

// NewToken creates a new token handler with provided developer seed, xdr and user
// address
func NewToken(devSecret, xdrs, address string) (*Token, error) {
	tx, err := transaction.New(xdrs)
	if err != nil {
		return nil, err
	}
	return &Token{
		devSecret,
		tx,
		address,
		nil,
		nil,
	}, nil
}

// TimeBounds returns the time bounds for a given transaction
func (t *Token) TimeBounds() *xdr.TimeBounds {
	timebounds := t.TX.Tx.Tx.TimeBounds
	if timebounds == nil {
		log.Fatalf("wrong challenge transaction structure")
	}
	return timebounds
}

// Validate transaction signed by developer and user
func (t *Token) Validate(strict bool) (bool, error) {
	if !t.signedCorrectly() {
		return false, fmt.Errorf("wrong challenge transaction signature")
	}
	if !t.timeNowCovers(t.TX.Tx.Tx.TimeBounds) {
		return false, fmt.Errorf("challenge transaction expired")
	}
	if strict && t.tooOld(t.TX.Tx.Tx.TimeBounds) {
		return false, fmt.Errorf("challenge transaction expired; too old")
	}
	return true, nil
}

// Hash validates and returns the transaction hash
func (t *Token) Hash(format string) interface{} {
	if _, err := t.Validate(true); err != nil {
		log.Fatal(err)
	}
	hash := t.TX.Hash()
	if format == "binary" {
		return hash
	}
	return fmt.Sprintf("%s", hex.EncodeToString(hash[:]))
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
	isSignedByDev := utils.ValidateTx(t.GetKeypair(), t.TX)
	isSignedByUser := utils.ValidateTx(t.GetUserKeypair(), t.TX)
	return isSignedByDev && isSignedByUser
}

// timeNowCovers returns true if current tie is within transaction time bounds
func (t *Token) timeNowCovers(timeBounds *xdr.TimeBounds) bool {
	now := fmt.Sprintf("%.1f", math.Floor(float64(tsm()/1000)))
	nMinTime := strconv.Itoa(int(timeBounds.MinTime))
	nMaxTime := strconv.Itoa(int(timeBounds.MaxTime))
	return (now >= nMinTime && now <= nMaxTime)
}

// tooOld returns true if transaction is created more than 10 seconds from now
func (t *Token) tooOld(timeBounds *xdr.TimeBounds) bool {
	now := fmt.Sprintf("%.1f", math.Floor(float64(tsm()/1000)))
	nMinTime, _ := strconv.ParseInt(strconv.Itoa(int(timeBounds.MinTime)), 10, 64)
	nMinTime = nMinTime + mc.StrictInterval
	return (now > strconv.Itoa(int(nMinTime)))
}
