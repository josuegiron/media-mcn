package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"./structures"
	"github.com/gorilla/mux"
	"github.com/mediocregopher/radix.v2/redis"
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

	token, err := GetTokenToCache()
	if err != nil {
		err := structures.Err{Code: 200, Messaje: "The song was not found..."}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&err)
	} else {

		url := "https://api.spotify.com/v1/search?q=" + url.QueryEscape(q) + "&type=track&market=US&limit=10&offset=1"

		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
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
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(&response)
			} else {
				err := structures.Err{Code: 200, Messaje: "The song was not found..."}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(&err)
			}

		}

	}

}

func DeezerRequest(q string, w http.ResponseWriter) {
	url := "https://api.deezer.com/search/track?q=" + url.QueryEscape(q)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	client2 := &http.Client{}
	resp, err := client2.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer resp.Body.Close()

	var deezer structures.DeezerResponse
	if err := json.NewDecoder(resp.Body).Decode(&deezer); err != nil {
		log.Println("Error en decodificar JSON: ", err)
	} else {
		total := len(deezer.Data)
		if total > 0 {
			var response structures.ResponseSong
			for i := 0; i < total; i++ {
				var song structures.Song
				song.Title = deezer.Data[i].Title
				song.Album = deezer.Data[i].Album.Title
				song.Artist = deezer.Data[i].Artist.Name
				song.Rank = deezer.Data[i].Rank
				song.Explicit = deezer.Data[i].ExplicitLyrics
				song.Length = deezer.Data[i].Duration * 1000
				song.URL = deezer.Data[i].Link
				song.TrackNumber = i

				response.Songs = append(response.Songs, song)
			}
			response.Next = deezer.Next
			response.Total = deezer.Total
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&response)
		} else {
			err := structures.Err{Code: 200, Messaje: "The song was not found..."}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&err)
		}

	}
}

//	Recoge el token de la base de datos en cache.

func GetTokenToCache() (string, error) {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	token, err := conn.Cmd("HGET", "Token", "access_token").Str()
	if err != nil {
		return "", err
	}

	return token, nil
}
