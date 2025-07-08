package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zahnfee144p/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
}

// function definitions

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not open database, exiting: %s", err)
		return 
	}

	mux := http.NewServeMux()
	filesrv := http.FileServer(http.Dir(filepathRoot))
	srvCfg := apiConfig{}
	srvCfg.dbQueries = database.New(db)

	// app pages
	mux.Handle("/app/", http.StripPrefix("/app", srvCfg.middlewareMetricsInc(filesrv)))

	// api calls
	mux.HandleFunc("GET /api/healthz", handleHealthz)
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	// admin pages
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
