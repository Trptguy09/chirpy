package main

import (
	"chirpy/internal/auth"
	"net/http"
)

type response struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to get token", err)
		return
	}
	user, err := cfg.db.GetUserByRefreshToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to get user", err)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to make authentication token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}
