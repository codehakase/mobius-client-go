package auth

import (
	"math"

	"github.com/codehakase/mobius-client-go/client"
	"github.com/codehakase/mobius-client-go/utils"
	"github.com/codehakase/mobius-client-go/utils/custom/transaction"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	time "github.com/ssoroka/ttime"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

var _ = Describe("Auth.Challenge", func() {
	var (
		tx  *transaction.Transaction
		kp  *keypair.Full
		err error
	)
	BeforeEach(func() {
		time.Freeze(time.Now())
		kp, err = keypair.Random()
		Expect(err).NotTo(HaveOccurred(), "random keypair generation shouldn't fail")
		ch := &Challenge{}
		tx, err = transaction.New(ch.Call(kp.Seed(), 0))
		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		time.Unfreeze()
	})
	Context("when a transaction is being made", func() {
		It("signs challenge transaction correctly by developer", func() {
			Expect(utils.ValidateTx(kp, tx)).To(Equal(true))
		})
		It("contains memo", func() {
			t, _ := tx.Tx.Tx.Memo.GetText()
			Expect(t).To(Equal("Mobius authentication"))
		})
		It("contains time bounds", func() {
			Expect(tx.Tx.Tx.TimeBounds).NotTo(BeNil())
		})
		It("contains correct minimum time bound", func() {
			now := xdr.Uint64(math.Floor(float64(tsm() / 1000)))
			Expect(tx.Tx.Tx.TimeBounds.MinTime).To(Equal(now))
		})
		It("contains correct maximum time bound", func() {
			now := xdr.Uint64(math.Floor(float64(tsm()/1000 + client.ChallengeExpiresIn)))
			Expect(tx.Tx.Tx.TimeBounds.MaxTime).To(Equal(now))
		})
		It("contains correct custom maximum time bound", func() {
			ch := &Challenge{}
			tx, err = transaction.New(ch.Call(kp.Seed(), 60))
			Expect(err).NotTo(HaveOccurred())
			now := xdr.Uint64(math.Floor(float64(tsm()/1000 + 60)))
			Expect(tx.Tx.Tx.TimeBounds.MaxTime).To(Equal(now))
		})
	})
})
