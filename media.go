package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"./structures"

	"github.com/gorilla/mux"
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

	url := "https://api.spotify.com/v1/search?q=" + url.QueryEscape(q) + "&type=track&market=ES&limit=10&offset=1"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer BQA241RQnONvt6d5fyMf0umxB_L_025PVk-J844exXEWl7vWMS10U3XoJKWQFmH-pjcJVhtQ8dDL_hMeJlNPUts8AufalDWrXMwWRdanx9nK78erYh54p1_porMCoxspsBrENubew29VaU2O0oda70kRA9VxGxOqbb_8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer resp.Body.Close()

	var spotify structures.SpotifyResponse

	if err := json.NewDecoder(resp.Body).Decode(&spotify); err != nil {
		log.Println("Error: ", err)
	} else {
		total := len(spotify.Tracks.Items)
		if total > 0 {
			for i := 0; i < total; i++ {

				fmt.Println("Album: ", spotify.Tracks.Items[i].Name, i)
			}
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
