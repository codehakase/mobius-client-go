package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/codehakase/mobius-client-go/app"
	"github.com/codehakase/mobius-client-go/auth"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var APP_KEY string

func init() {
	APP_KEY = os.Getenv("APPLICATION_SECRET_KEY")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		chl := &auth.Challenge{}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte((chl.Call(APP_KEY, 0))))
	}).Methods("GET")

	r.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		token, err := auth.NewToken(
			APP_KEY,
			r.FormValue("xdr"),
			r.FormValue("public_key"),
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		// Important! Otherwise, token will be considered valid
		if _, err := token.Validate(true); err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		j := &auth.JWT{Secret: APP_KEY}
		t := j.Encode(token, nil)
		w.Write([]byte(t))
	}).Methods("POST")

	r.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		publicKey := getPublicKeyFromToken(r.URL.Query().Get("token"))
		dapp, err := app.Build(APP_KEY, publicKey)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%v", dapp.UserBalance())))
	}).Methods("GET")

	r.HandleFunc("/pay", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		publicKey := getPublicKeyFromToken(r.URL.Query().Get("token"))
		err := r.ParseForm()
		if err != nil {
			log.Fatalf("failed to parse request, %v", err)
		}
		var amount float64
		if r.FormValue("amount") != "" {
			amount, _ = strconv.ParseFloat(r.FormValue("amount"), 64)
		} else {
			amount = 1
		}
		dapp, err := app.Build(APP_KEY, publicKey)
		if err != nil {
			log.Fatal(err)
		}
		res, err := dapp.Charge(amount)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(
			"{'status': 'ok','tx_hash': '%s','balance': %v}",
			res.Hash,
			dapp.UserBalance(),
		)))
	}).Methods("POST")

	r.HandleFunc("/payout", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
		targetAddress, publicKey := r.FormValue("target_address"), getPublicKeyFromToken(r.URL.Query().Get("token"))
		if amount < 1 {
			log.Fatal("invalid amount")
		}
		if targetAddress == "" {
			log.Fatal("invalid target address")
		}
		dapp, err := app.Build(APP_KEY, publicKey)
		if err != nil {
			log.Fatal(err)
		}
		res, err := dapp.Payout(float64(amount), targetAddress)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("{'status': 'ok','tx_hash': '%s', 'balance': %v}", res.Hash, dapp.UserBalance())))
	}).Methods("POST")

	r.HandleFunc("/transfer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
		targetAddress, publicKey := r.FormValue("target_address"), getPublicKeyFromToken(r.URL.Query().Get("token"))
		if amount < 1 {
			log.Fatal("invalid amount")
		}
		if targetAddress == "" {
			log.Fatal("invalid target address")
		}
		dapp, err := app.Build(APP_KEY, publicKey)
		if err != nil {
			log.Fatal(err)
		}
		res, err := dapp.Transfer(amount, targetAddress)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("{'status': 'ok','tx_hash': '%s', 'balance': %v}", res.Hash, dapp.UserBalance())))

	}).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders: []string{"Origin", "X-Requested-With", "Content-Type",
			"Accept"},
	})
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./"))))
	handler := c.Handler(r)
	log.Println("Starting app server...")
	log.Fatal(http.ListenAndServe(":3000", handler))

}

func getPublicKeyFromToken(token string) string {
	if token == "" {
		log.Fatal("token is sent empty")
		os.Exit(1)
	}
	j := &auth.JWT{Secret: APP_KEY}
	jwtToken := j.Decode(token)
	return jwtToken.Claims.(jwt.MapClaims)["sub"].(string)
}
