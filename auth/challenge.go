package auth

import (
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/ssoroka/ttime"

	mc "github.com/codehakase/mobius-client-go/client"
	"github.com/codehakase/mobius-client-go/utils"
	"github.com/stellar/go/build"
	"github.com/stellar/go/network"
)

var Network build.Network

func init() {
	rand.Seed(ttime.Now().UnixNano())
	if os.Getenv("MOBIUS_NETWORK") == "test" {
		Network = build.Network{Passphrase: network.TestNetworkPassphrase}
	} else if os.Getenv("MOBIUS_NETWORK") == "public" {
		Network = build.Network{Passphrase: network.PublicNetworkPassphrase}
	}
}

// Challenge represents a model which generates challenge transactions on
// developer's account
type Challenge struct{}

// Call generates a challenge transaction signed by the developer's private key
func (c *Challenge) Call(devSecret string, expiresIn int64) string {
	if expiresIn < 1 {
		expiresIn = mc.ChallengeExpiresIn
	}
	kp := utils.KPFromSeed(devSecret)
	// randomKp, _ := keypair.Random()
	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: kp.Seed()},
		build.Sequence{Sequence: c.randomSequence()},
		Network,
		c.memo(),
		c.buildTimeBounds(expiresIn),
		build.Payment(
			build.Destination{AddressOrSeed: kp.Seed()},
			build.NativeAmount{Amount: "0.000001"},
		),
	)
	if err != nil {
		log.Fatalf("failed to build challenge transaction, err: %v", err)
	}
	txe, err := tx.Sign(kp.Seed())
	if err != nil {
		log.Fatalf("failed to sign challenge transaction, err: %v", err)
	}
	txtEnvelopeStr, err := txe.Base64()
	if err != nil {
		log.Fatalf("failed to generate txt envelope str for challenge txt, err: %v", err)
	}
	return txtEnvelopeStr
}

func (c *Challenge) randomSequence() uint64 {
	return uint64(99999999 - math.Floor(rand.Float64()*65536))
}

func (c *Challenge) memo() build.MemoText {
	return build.MemoText{Value: "Mobius authentication"}
}

func (c *Challenge) buildTimeBounds(exp int64) build.Timebounds {
	if exp < 1 {
		exp = mc.ChallengeExpiresIn
	}
	return build.Timebounds{
		MinTime: uint64(math.Floor(float64(tsm() / 1000))),
		MaxTime: uint64(math.Floor(float64(tsm()/1000 + exp))),
	}
}

func tsm() int64 {
	return ttime.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
