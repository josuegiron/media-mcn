package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"./structures"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

func main() {
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
			SpotifyRequest(q, w)
		case "deezer":
			DeezerRequest(q, w)
		default:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&err)
		}
	}
}

func SpotifyRequest(q string, w http.ResponseWriter) {

	ctx := context.Background()

	conf := &oauth2.Config{
		ClientID:     "fea07c92beaf4a0ba83d42243461c997",
		ClientSecret: "0c2da7d6085f4ddba9819ffa3ddfe726",
		Scopes:       []string{"user-library-read"},
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://accounts.spotify.com/api/token",
			AuthURL:  "https://accounts.spotify.com/authorize",
		},
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	//url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	//fmt.Printf("Visit the URL for the auth dialog: %v", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	// Use the custom HTTP client when requesting a token.
	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)
	resss, _ := client.Get("https://api.spotify.com/v1/search?q=Oceans Hillsong&type=track&market=ES&limit=10&offset=11")
	log.Println(resss)

	url2 := "https://api.spotify.com/v1/search?q=" + q + "&type=track&market=ES&limit=10&offset=1"

	req, err := http.NewRequest("GET", url2, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer BQBPgazQ9K_rsiyGf4lM83z8ZTGgRMJq2JBXGVRzPlrckdzQqxruU_oPT-h93Nu4Wg_ANXuAF-hVuk4hm_RAphsMwa39SbFiEIN3iiS1ODI-Iflg-fRy7Pruln89UCLL81pgXAH8MX1yDLLnTPdVHl50_6FZOIUT8g")
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
