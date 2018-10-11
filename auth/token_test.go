package auth

import (
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	time "github.com/ssoroka/ttime"
)

var _ = Describe("Auth.Token", func() {
	BeforeEach(func() {
		genKeypairs()
	})
	AfterEach(func() {
		time.Unfreeze()
	})
	It("Token.Validate returns true if current time is within bounds", func() {
		time.Freeze(time.Now())
		tx := generateSignedChallengeTx(userKeypair, devKeypair)
		token, err := NewToken(devKeypair.Seed(), tx, userKeypair.Address())
		Expect(err).ToNot(HaveOccurred())
		ok, err := token.Validate(true)
		Expect(err).ToNot(HaveOccurred())
		Expect(ok).To(Equal(true))
	})
	It("Token.Validate throws and error if current time is outside bounds", func() {
		tx := generateSignedChallengeTx(userKeypair, devKeypair)
		token, err := NewToken(devKeypair.Seed(), tx, userKeypair.Address())
		Expect(err).ToNot(HaveOccurred())
		i := math.Floor(float64(tsm()/1000 + 3600*5))
		futureTime := time.Unix(int64(i), 0)
		time.Freeze(futureTime)
		_, err = token.Validate(true)
		Expect(err).To(HaveOccurred())
	})
	It("returns transaction hash", func() {
		tx := generateSignedChallengeTx(userKeypair, devKeypair)
		token, err := NewToken(devKeypair.Seed(), tx, userKeypair.Address())
		Expect(err).ToNot(HaveOccurred())
		Expect(token.Hash("")).ToNot(BeEmpty())
	})
})
