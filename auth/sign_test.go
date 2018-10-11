package auth

import (
	"github.com/codehakase/mobius-client-go/utils"
	"github.com/codehakase/mobius-client-go/utils/custom/transaction"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth.Sign", func() {
	BeforeEach(func() {
		genKeypairs()
	})
	It("signs challenge correctly by user", func() {
		tx := generateSignedChallengeTx(userKeypair, devKeypair)
		txt, err := transaction.New(tx)
		Expect(err).ToNot(HaveOccurred())
		Expect(utils.ValidateTx(userKeypair, txt)).To(Equal(true))
	})
})
