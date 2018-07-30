package mobius

import (
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

var (
	// Issuers set for the mobius sdk
	Issuers = map[string]string{
		"PUBLIC":  "GA6HCMBLTZS5VYYBCATRBRZ3BZJMAFUDKYYF6AH6MVCMGWMRDNSWJPIH",
		"TESTNET": "GDRWBLJURXUKM4RWDZDTPJNX6XBYFO3PSE4H4GPUL6H6RCUQVKTSD4AT",
	}
	// Urls set to access the horizon server
	Urls = map[string]string{
		"PUBLIC":  "https://horizon.stellar.org",
		"TESTNET": "https://horizon-testnet.stellar.org",
	}
)

// Client is the entrypoint for the sdk.
type Client struct {
	MobiusHost         string
	StrictInterval     int
	ChallengeExpiresIn int
	AssetCode          string
	HorizonClient      *horizon.Client
	Network            *build.Network
	StellarAsset       build.Asset
}

// NewClient creates a new client with the defaults settings from sdk's
// design
func NewClient() *Client {
	client := &Client{
		MobiusHost:         "https://mobius.network",
		StrictInterval:     10,
		ChallengeExpiresIn: 60 * 60 * 24,
		AssetCode:          "MOBI",
		HorizonClient:      nil,
	}
	client.Network = &build.TestNetwork
	return client
}

// GetAssetIssuer retunrs the account ID of the Asset Issuer
func (c *Client) GetAssetIssuer() string {
	if c.Network.Passphrase == build.PublicNetwork.Passphrase {
		return Issuers["PUBLIC"]
	}
	return Issuers["TESTNET"]
}

// GetChallengeExpiration returns the challenge expiration value in seconds
// (defaults to 1 day)
func (c *Client) GetChallengeExpiration() int {
	return c.ChallengeExpiresIn
}

// GetMobiusHost returns the Mobius API host
func (c *Client) GetMobiusHost() string {
	return c.MobiusHost
}

// GetStellarAsset returns an instance of the asset used for payments
func (c *Client) GetStellarAsset() build.Asset {
	if (build.Asset{}) == c.StellarAsset {
		return c.StellarAsset
	}
	stellarAsset := build.CreditAsset(c.AssetCode, c.GetAssetIssuer())
	c.StellarAsset = stellarAsset
	return stellarAsset
}

// GetHorizonClient returns a StellarHorizon Server instance
func (c *Client) GetHorizonClient() *horizon.Client {
	if c.HorizonClient != nil {
		return c.HorizonClient
	}
	if c.Network.Passphrase == build.PublicNetwork.Passphrase {
		c.HorizonClient = horizon.DefaultPublicNetClient
		return c.HorizonClient
	}
	c.HorizonClient = horizon.DefaultTestNetClient
	return c.HorizonClient
}

// UsePublicNetwork sets network to mainnet
func (c *Client) UsePublicNetwork() {
	c.Network = &build.PublicNetwork
}

// UseTestNetwork sets network to testnet
func (c *Client) UseTestNetwork() {
	c.Network = &build.TestNetwork
}
