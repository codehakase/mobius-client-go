package auth

import (
	"log"
	"math"
	"math/rand"
	"time"

	mc "github.com/codehakase/mobius-client-go/client"
	"github.com/codehakase/mobius-client-go/utils"
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Challenge represents a model which generates challenge transactions on
// developer's account
type Challenge struct{}

// Call generates a challenge transaction signed by the developer's private key
func (c *Challenge) Call(devSecret string, expersIn int) string {
	kp := utils.KPFromSeed(devSecret)
	randomKp, _ := keypair.Random()
	tx, err := build.Transaction(
		build.Payment(
			build.SourceAccount{AddressOrSeed: randomKp.Address()},
			build.Destination{AddressOrSeed: kp.Address()},
			build.Sequence{Sequence: c.randomSequence()},
			build.NativeAmount{Amount: "0.000001"},
		),
	)
	if err != nil {
		log.Fatalf("failed to build challenge transaction, err: %v", err)
	}
	txe, err := tx.Sign(kp.Address())
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
		MinTime: uint64(math.Floor(float64(time.Now().UnixNano() / 1000))),
		MaxTime: uint64(math.Floor(float64(time.Now().UnixNano()/1000 + exp))),
	}
}
