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

	"./structures"
	"github.com/gorilla/mux"
)

var c structures.Credentials

func main() {

	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &c)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/search/song", Media).Methods("GET")
	log.Fatal(http.ListenAndServe(":3006", router))
}

func Media(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := structures.Err{Code: 400, Messaje: "The required information could not be obtained..."}

	q := r.URL.Query().Get("q")
	provider := r.URL.Query().Get("provider")

	if q == "" || provider == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&err)
	} else {
		switch provider {
		case "spotify":
			SpotifyRequest(q, r, w)
		case "deezer":
			DeezerRequest(q, w)
		default:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&err)
		}
	}
}

func SpotifyRequest(q string, r *http.Request, w http.ResponseWriter) {

	body := url.Values{"grant_type": {"client_credentials"}}

	var urlAuth = "https://accounts.spotify.com/api/token"

	reqToken, er := http.NewRequest("POST", urlAuth, strings.NewReader(body.Encode()))

	cidcs := []byte(c.ClientId + ":" + c.ClientSecret)
	secret := base64.StdEncoding.EncodeToString(cidcs)
	basicAuth := "Basic " + secret
	reqToken.Header.Add("Authorization", basicAuth)
	reqToken.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	clientToken := &http.Client{}
	resToken, er := clientToken.Do(reqToken)
	if er != nil {
		panic(er)
	}
	defer resToken.Body.Close()

	var token structures.Token
	json.NewDecoder(resToken.Body).Decode(&token)

	url := "https://api.spotify.com/v1/search?q=" + url.QueryEscape(q) + "&type=track&market=ES&limit=10&offset=1"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.AccesToken)
	client2 := &http.Client{}
	resp, err := client2.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer resp.Body.Close()

	var spotify structures.SpotifyResponse

	if err := json.NewDecoder(resp.Body).Decode(&spotify); err != nil {
		log.Println("Error en decodificar JSON: ", err)
	} else {
		total := len(spotify.Tracks.Items)
		if total > 0 {
			var response structures.ResponseSong
			for i := 0; i < total; i++ {
				var song structures.Song
				song.Title = spotify.Tracks.Items[i].Name
				song.Album = spotify.Tracks.Items[i].Album.Name
				song.Artist = spotify.Tracks.Items[i].Artists[0].Name
				song.Rank = spotify.Tracks.Items[i].Popularity
				song.Released = spotify.Tracks.Items[i].Album.ReleaseDate
				song.Explicit = spotify.Tracks.Items[i].Explicit
				song.Length = spotify.Tracks.Items[i].DurationMs
				song.URL = spotify.Tracks.Items[i].ExternalUrls.Spotify
				song.TrackNumber = spotify.Tracks.Items[i].TrackNumber

				response.Songs = append(response.Songs, song)
			}
			response.Next = spotify.Tracks.Next
			response.Previous = spotify.Tracks.Previous
			response.Total = spotify.Tracks.Total
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&response)
		} else {
			err := structures.Err{Code: 200, Messaje: "The song was not found..."}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&err)
		}

	}

}
func DeezerRequest(q string, w http.ResponseWriter) {
	deezer := map[string]string{"message": "Esta es la respuesta de la informacion de Deezer...", "q": q}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deezer)
}
