package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf(
`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, 
		cfg.fileserverHits.Load())
	w.Write([]byte(body))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := "Hits reset to 0"
	cfg.fileserverHits.Store(0)
	w.Write([]byte(body))
}

func handleHealthz(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
