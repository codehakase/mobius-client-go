# mobius-client-go

# Mobius DApp Store Go SDK

The Mobius DApp Store JS SDK makes it easy to integrate Mobius DApp Store MOBI payments into any Go application.

A big advantage of the Mobius DApp Store over centralized competitors such as the Apple App Store or Google Play Store is significantly lower fees - currently 0% compared to 30% - for in-app purchases.

## DApp Store Overview

The Mobius DApp Store will be an open-source, non-custodial "wallet" interface for easily sending crypto payments to apps. You can think of the DApp Store like https://stellarterm.com/ or https://www.myetherwallet.com/ but instead of a wallet interface it is an App Store interface.

The DApp Store is non-custodial meaning Mobius never holds the secret key of either the user or developer.

An overview of the DApp Store architecture is:

- Every application holds the private key for the account where it receives MOBI.
- An application specific unique account where a user deposits MOBI for use with the application is generated for each app based on the user's seed phrase.
- When a user opens an app through the DApp Store:
  1) Adds the application's public key as a signer so the application can access the MOBI and
  2) Signs a challenge transaction from the app with its secret key to authenticate that this user owns the account. This prevents a different person from pretending they own the account and spending the MOBI (more below under Authentication).

### Installation
Install with `go get`
```shell
$ go get github.com/codehakase/mobius-client-go
```

## Production Server Setup
Your production server must use HTTPS and set the below header on the `/auth` endpoint:

`Access-Control-Allow-Origin: *`

### Explanation
When a user opens an app through the DApp Store it tells the app what Mobius account it should use for payment.

The application needs to ensure that the user actually owns the secret key to the Mobius account and that this isn't a replay attack from a user who captured a previous request and is replaying it.

This authentication is accomplished through the following process:

When the user opens an app in the DApp Store it requests a challenge from the application.
The challenge is a payment transaction of 1 XLM from and to the application account. It is never sent to the network - it is just used for authentication.
The application generates the challenge transaction on request, signs it with its own private key, and sends it to user.
The user receives the challenge transaction and verifies it is signed by the application's secret key by checking it against the application's published public key (that it receives through the DApp Store). Then the user signs the transaction with its own private key and sends it back to application along with its public key.
Application checks that challenge transaction is now signed by itself and the public key that was passed in. Time bounds are also checked to make sure this isn't a replay attack. If everything passes the server replies with a token the application can pass in to "login" with the specified public key and use it for payment (it would have previously given the app access to the public key by adding the app's public key as a signer).
Note: the challenge transaction also has time bounds to restrict the time window when it can be used.

See demo at:

```shell
$ git clone https://github.com/codehakase/mobius-client-go.git $GOPATH/src/github.com/codehakase/mobius-client-go

$ cd $GOPATH/src/github.com/codehakase/mobius-client-go/example

$ go run main.go 

# navigate to http://localhost:3000 for demo

```

### Sample Server Implementation

```go
package main

import (
	"net/http"
	"encoding/json"

	MobiusClient "github.com/codehakase/mobius-client-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)
func main() {
	r := mux.NewRouter()
	// GET /auth
	// Generates and returns challenge transaction XDR signed by application to user
	r.HandleFunc("/auth", func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		chl := &MobiusClient.Auth.Challenge{}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte((chl.Call(APPLICATION_SECRET_KEY, 0)))
	}).Methods("GET")

	// POST /auth
	// Validates challenge transaction. It must be:
	//  - Signed by application and requesting user.
	//  - Not older than 10 seconds from now (see MobiusClient.Client.strictInterval`)
	type challengeReq struct{
		xdr       string `json:"xdr"`
		publicKey string `json:"public_key"`
	}
	r.HandleFunc("/auth", func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		var chreq challengeReq
		_ := json.NewDecoder(r.Body).Decode(&chreq) // handle errors
		token, err := MobiusClient.Auth.NewToken(
			APPLICATION_SECRET,
			chreq.xdr,
			chreq.publicKey	
		)
		if err != nil {
			http.Error(w, err, 500)
		}
		// Important! Otherwise, token will be considered valid
		_ = token.Validate(true)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token.Hash("hex")))
	}).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowCredentials: true,
		AllowedHeaders: []string{"Origin", "X-Requested-With", "Content-Type",
		"Accept"},
	})
	handler := c.Handler(r)
	log.Fatal(http.ListenAndServe(":3000", handler))
}
```

## Payment

### Explanation
After the user completes the authentication process they have a token. They now
pass it to the application to "login" which tells the application which Mobius
account to withdraw MOBI from (the user public key) when a payment is needed.
For a web application the token is generally passed in via a token request
parameter. Upon opening the website/loading the application it checks that the
token is valid (within time bounds etc) and the account in the token has added
the app as a signer so it can withdraw MOBI from it.

## Development
```shell
# Clone this repo

$ git clone https://github.com/codehakase/mobius-client-go.git $GOPATH/src/github.com/codehakase/mobius-client-go && cd $_

# Install dependencies (using go dep) http://github.com/golang/dep

$ dep ensure -v

# Run authentication example

$ make example:auth

# Run Tests

$ go test ./...
```

## Contributing

Bug reports and pull requests are welcome on GitHub at
https://github.com/codehakase/mobius-client-go. This project is intended to
be a safe, welcoming space for collaboration, and contributors are expected to
adhere to the [Contributor Covenant](http://contributor-covenant.org) code of
conduct.

## License

The SDK is available as open source under the terms of the [MIT
License](https://opensource.org/licenses/MIT).
