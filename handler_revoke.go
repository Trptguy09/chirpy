package main

import (
	"chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to get token", err)
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to revoke token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
