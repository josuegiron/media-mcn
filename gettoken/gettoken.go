package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/auth", SetToken).Methods("GET")
	log.Fatal(http.ListenAndServe(":3010", router))
}

func SetToken(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("holi si funciona")
	log.Println(r.URL.Query().Get("code"))

}

type Credentials struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

var (
	conf *oauth2.Config
	ctx  context.Context
)
var c Credentials

func Init() {

	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &c)

	ctx = context.Background()
	conf = &oauth2.Config{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		Scopes:       []string{""},
		Endpoint:     spotify.Endpoint,
		RedirectURL:  c.RedirectURL,
	}

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)
	// Handle the exchange code to initiate a transport.

}
