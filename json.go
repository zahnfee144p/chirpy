package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("content-type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func replaceProfanity(body string) string {
	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	const profanityRepl = "****"

	split := strings.Split(body, " ")
	for idx, word := range split {
		if slices.Contains(profanity, strings.ToLower(word)) {
			split[idx] = profanityRepl
		}
	}
	return strings.Join(split, " ")
}
