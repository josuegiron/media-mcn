package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/mediocregopher/radix.v2/redis"
)

type Credentials struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

type Token struct {
	AccesToken string `json:"access_token"`
	TokenType  string `json:"token_type"`
	ExpiresIn  int    `json:"expires_in"`
	Scope      string `json:"scope"`
}

var c Credentials
var token Token

func main() {
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &c)

	for {
		GetTokenToSpotify()
		time.Sleep(3599 * time.Second)
	}
}

func GetTokenToSpotify() {

	body := url.Values{"grant_type": {"client_credentials"}}
	urlAuth := "https://accounts.spotify.com/api/token"
	credentials := []byte(c.ClientId + ":" + c.ClientSecret)
	credentialsEncode := base64.StdEncoding.EncodeToString(credentials)
	basicAuth := "Basic " + credentialsEncode

	reqToken, er := http.NewRequest("POST", urlAuth, strings.NewReader(body.Encode()))
	reqToken.Header.Add("Authorization", basicAuth)
	reqToken.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	clientToken := &http.Client{}
	resToken, er := clientToken.Do(reqToken)
	if er != nil {
		log.Println(er)
	}
	defer resToken.Body.Close()
	log.Println(resToken.StatusCode)

	json.NewDecoder(resToken.Body).Decode(&token)

	log.Println("Token: ", token.AccesToken)

	SetTokenToCache()
}

func SetTokenToCache() {

	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	err = conn.Cmd("HMSET", "Token", "access_token", token.AccesToken, "token_type", token.TokenType, "scope", token.Scope, "expire_in", token.ExpiresIn).Err
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Cmd("EXPIRE", "Token", token.ExpiresIn).Err
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Actualiza token en cache en redis...")
}
