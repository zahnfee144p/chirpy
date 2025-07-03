package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	w.Header().Add("content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`,
		cfg.fileserverHits.Load())
	w.Write([]byte(body))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	w.Header().Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := "Hits reset to 0"
	cfg.fileserverHits.Store(0)
	w.Write([]byte(body))
}

func handleHealthz(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	w.Header().Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handleValidateChirp(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	const ChirpMaxLength = 140
	const ProfanityRepl = "****"

	type chirp struct {
		Body  string `json:"body"`
		Error string `json:"error"`
		Valid bool   `json:"valid"`
	}

	decoder := json.NewDecoder(req.Body)
	data := chirp{}

	if err := decoder.Decode(&data); err != nil {
		log.Printf("Error decoding request body: %v", err)
		data.Error = "Something went wrong"

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonData)
		return
	}

	if len(data.Body) > ChirpMaxLength {
		data.Error = "Chirp is too long"
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonData)
		return
	}

	data.Valid = true
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// func replaceProfanity(body string) string {
// 	return ""
// }
