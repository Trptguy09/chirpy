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
type ChirpRequest struct {
	Body string `json:"body"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}
type ValidResponse struct {
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req ChirpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Something went wrong",
		})
		return
	}

	if len(req.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Chirp is too long",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ValidResponse{
		Valid: true,
	})
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
