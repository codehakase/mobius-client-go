package client_test

import (
	mc "github.com/codehakase/mobius-client-go/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/network"
)

var _ = Describe("Client", func() {
	var (
		client *mc.Client
	)
	Context("should respect global Network config..", func() {
		BeforeEach(func() {
			client = mc.NewClient()
		})
		It("should test on test network...", func() {
			client.UseTestNetwork()
			Expect(client.Network.Passphrase).To(Equal(network.TestNetworkPassphrase))
		})

		It("should test on public network", func() {
			client.UsePublicNetwork()
			Expect(client.Network.Passphrase).To(Equal(network.PublicNetworkPassphrase))
		})

		It("should return the correct testnet issuers", func() {
			client.UseTestNetwork()
			Expect(client.GetAssetIssuer()).To(Equal(mc.Issuers["TESTNET"]))
		})

		It("should return the correct publicnet issuers", func() {
			client.UsePublicNetwork()
			Expect(client.GetAssetIssuer()).To(Equal(mc.Issuers["PUBLIC"]))
		})
		It("should return the challenge expiration not less than 24 hours", func() {
			Expect(client.GetChallengeExpiration()).To(Equal(int64(60 * 60 * 24)))
		})
	})
	Context("client should reflect current Network settings", func() {
		BeforeEach(func() {
			client = mc.NewClient()
		})
		It("should return DefaultPublicNetClient on PublicNetwork", func() {
			client.UsePublicNetwork()
			Expect(client.GetHorizonClient()).To(Equal(horizon.DefaultPublicNetClient))
		})

		It("should return DefaultTestNetClient  on TestNetwork", func() {
			client.UseTestNetwork()
			Expect(client.GetHorizonClient()).To(Equal(horizon.DefaultTestNetClient))
		})

	})
})
