package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type parameters struct {
	Body string `json:"body"`
}

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

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "forbidden", errors.New("tried to call reset outside of dev environment"))
		return
	}
	
	cfg.fileserverHits.Store(0)
	cfg.dbQueries.DeleteAllUsers(req.Context())

	w.Header().Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := "Hits reset to 0, all users deleted"
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
	const chirpMaxLength = 140

	type returnVals struct {
		Valid bool   `json:"valid"`
		Body  string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding parameters", err)
		return
	}

	if len(params.Body) > chirpMaxLength {
		respondWithError(w, http.StatusBadRequest, "chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Valid: true,
		Body:  replaceProfanity(params.Body),
	})
}

func (cfg *apiConfig) handleUsers(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding parameters", err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
