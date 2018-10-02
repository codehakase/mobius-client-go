package app

import (
	"fmt"
	"log"
	"strconv"

	"github.com/codehakase/mobius-client-go/blockchain"
	mc "github.com/codehakase/mobius-client-go/client"
	"github.com/codehakase/mobius-client-go/utils"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

// App defines an interface to the user in a DApp
type App struct {
	AppAccount     *blockchain.Account
	ClientInstance *horizon.Client
	UserAccount    *blockchain.Account
}

// Build fetches information from the Stellar Network, and returns an instance
// of a connected DApp
func Build(developerSecret, address string) (*App, error) {
	devKeypair := utils.KPFromSeed(developerSecret)
	devAccount, err := blockchain.Build(devKeypair)
	if err != nil {
		return nil, err
	}
	userKeypair := utils.KPFromSeed(address)
	userAccount, err := blockchain.Build(userKeypair)
	if err != nil {
		return nil, err
	}
	return &App{
		devAccount,
		mc.NewClient().HorizonClient,
		userAccount,
	}, nil
}

// Authorized confirms the connected developer is authorized to use app
func (a *App) Authorized() bool {
	return a.UserAccount.Authorized(a.AppKeypair())
}

// AppAccount returns the associated app account
func (a *App) Account() *blockchain.Account { return a.AppAccount }

// Balance returns the app balance
func (a *App) Balance() float64 { return a.AppAccount.Balance() }

// GetUserAccount returns the user account authenticated to app
func (a *App) GetUserAccount() *blockchain.Account { return a.UserAccount }

// UserBalance reutrns the app user balance
func (a *App) UserBalance() float64 {
	// validate user balance
	if !a.Authorized() {
		log.Fatal("Authorization is missing")
	}
	if !a.UserAccount.TrustLineExists() {
		log.Fatal("Trustline not found for user account")
	}
	return a.UserAccount.Balance()
}

// AppKeypair returns keypair associated with app account
func (a *App) AppKeypair() *keypair.Full { return a.AppAccount.GetKeypair() }

// UserKeypair returns keypair associated with connected user
func (a *App) UserKeypair() *keypair.Full { return a.UserAccount.GetKeypair() }

// Charge charges specified amount from user account and then optionally
// transfers it from app account to a thrid party in same transaction
func (a *App) Charge(amount float64, destination ...string) (horizon.TransactionSuccess, error) {
	if a.UserBalance() < amount {
		log.Fatal("Insufficient Funds in user wallet to cover transaction")
	}
	var payment build.PaymentBuilder
	payment = a.userPaymentOp(amount, a.AppKeypair().Address())
	if destination != nil {
		payment.Mutate(a.appPaymentOp(amount, destination[0]))
	}
	return a.submitTx(payment)
}

// Payout sends money from application account to a user or third party.
func (a *App) Payout(amount float64, destination ...string) (horizon.TransactionSuccess, error) {
	if a.Balance() < amount {
		log.Fatal("Insufficient Funds in app wallet to cover transaction")
	}
	var payto string
	// set destination to user wallet if destination is nil
	if len(destination) < 1 {
		payto = a.UserKeypair().Address()
	} else {
		payto = destination[0]
	}
	return a.submitTx(a.appPaymentOp(amount, payto))
}

// Transfer sends money from the user account to a thrid party directly
func (a *App) Transfer(amount float64, destination string) (horizon.TransactionSuccess, error) {
	if a.UserBalance() < amount {
		log.Fatal("Insufficient Funds in user wallet to cover transaction")
	}
	return a.submitTx(a.userPaymentOp(amount, destination))
}

//

// submitTx builds and submit a transaction to the Stellar Network
func (a *App) submitTx(paymentOps build.PaymentBuilder) (horizon.TransactionSuccess, error) {
	var hts horizon.TransactionSuccess
	tx, err := build.Transaction(
		paymentOps,
	)
	if err != nil {
		return hts, fmt.Errorf("failed to build transaction, err: %v", err)
	}
	// sign transaction
	txe, err := tx.Sign(a.AppKeypair().Seed())
	if err != nil {
		return hts, fmt.Errorf("failed to sign transaction, err: %v", err)
	}
	txtEnvelopeStr, err := txe.Base64()
	if err != nil {
		return hts, fmt.Errorf("failed to generate txt envelope str, err: %v", err)
	}
	return a.ClientInstance.SubmitTransaction(txtEnvelopeStr)
}

// userPaymentOp creates a new payment operation setting source to app user
func (a *App) userPaymentOp(amount float64, destination string) build.PaymentBuilder {
	return build.Payment(
		build.Destination{AddressOrSeed: destination},
		build.NativeAmount{Amount: strconv.FormatFloat(amount, 'g', -1, 64)},
		build.SourceAccount{AddressOrSeed: a.UserKeypair().Address()},
	)
}

// appPaymentOp creates a new payment operation setting source to app account
func (a *App) appPaymentOp(amount float64, destination string) build.PaymentBuilder {
	return build.Payment(
		build.Destination{AddressOrSeed: destination},
		build.NativeAmount{Amount: strconv.FormatFloat(amount, 'g', -1, 64)},
		build.SourceAccount{AddressOrSeed: a.AppKeypair().Address()},
	)
}
