package main

import (
	"log"
	"net/http"
	"os"

	MobiusAuth "github.com/codehakase/mobius-client-go/auth"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var APPLICATION_SECRET_KEY string

func init() {
	APPLICATION_SECRET_KEY = os.Getenv("APPLICATION_SECRET_KEY")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		chl := &MobiusAuth.Challenge{}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte((chl.Call(APPLICATION_SECRET_KEY, 0))))
	}).Methods("GET")

	r.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		token, err := MobiusAuth.NewToken(
			APPLICATION_SECRET_KEY,
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
		h := token.Hash("hex").(string)
		w.Write([]byte(h))
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
