package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}
type inParams struct {
	Body string `json:"body"`
}

type errorResp struct {
	Error string `json:"error"`
}

type validResp struct {
	Valid bool `json:"valid"`
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d", cfg.fileserverHits.Load())

}

func (cfg *apiConfig) validate_chirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := inParams{}
	err := decoder.Decode(&params)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		resp, err := json.Marshal(errorResp{
			Error: "Something went wrong",
		})
		if err != nil {
			return
		}
		w.Write(resp)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		resp, err := json.Marshal(errorResp{
			Error: "Chirp is too long",
		})
		if err != nil {
			return
		}
		w.Write(resp)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp, err := json.Marshal(validResp{
		Valid: true,
	})
	if err != nil {
		return
	}
	w.Write(resp)
}

func main() {
	apiCfg := &apiConfig{}
	serverMux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	serverMux.Handle("/app", apiCfg.middlewareMetricsInc(handler))
	serverMux.HandleFunc("GET /api/healthz", apiCfg.healthz)
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.metrics)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.reset)
	serverMux.HandleFunc("POST /api/validate_chirp", apiCfg.validate_chirp)

	srv := http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

}
