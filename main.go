package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

// function definitions

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	filesrv := http.FileServer(http.Dir(filepathRoot))
	srvCfg := apiConfig{}

	mux.Handle("/app/", http.StripPrefix("/app", srvCfg.middlewareMetricsInc(filesrv)))
	mux.HandleFunc("GET /api/healthz", handleHealthz)
	mux.HandleFunc("GET /admin/metrics", srvCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", srvCfg.handleReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

// method definitions

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
