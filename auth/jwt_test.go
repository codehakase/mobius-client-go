package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	time "github.com/ssoroka/ttime"
	"github.com/stellar/go/keypair"
)

var (
	userKeypair *keypair.Full
	devKeypair  *keypair.Full
)
var _ = Describe("Auth.Token", func() {
	var (
		token *Token
		jwts  *JWT
		tx    string
		err   error
	)
	BeforeEach(func() {
		time.Freeze(time.Now())
		genKeypairs()
		tx = generateSignedChallengeTx(userKeypair, devKeypair)
		token, err = NewToken(devKeypair.Seed(), tx, userKeypair.Address())
		Expect(err).NotTo(HaveOccurred())
		jwts = &JWT{Secret: "somekey"}
	})
	AfterEach(func() {
		time.Unfreeze()
	})
	It("jwt.Encode returns string jwt", func() {
		Expect(jwts.Encode(token, nil)).NotTo(BeEmpty())
	})
	It("jwt.Decode returns payload", func() {
		payload := jwts.Decode(jwts.Encode(token, nil))
		Expect(payload.Claims.(jwt.MapClaims)["sub"]).To(Equal(devKeypair.Address()))
	})
})

func generateSignedChallengeTx(userkeypair, devkeypair *keypair.Full) string {
	tx := new(Challenge).Call(devkeypair.Seed(), 0)
	signedTx := new(Sign).Call(userkeypair.Seed(), tx, devkeypair.Address())
	return signedTx
}

func genKeypairs() {
	if userKeypair == nil {
		userKeypair, _ = keypair.Random()
	}
	if devKeypair == nil {
		devKeypair, _ = keypair.Random()
	}
}
