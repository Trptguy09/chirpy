package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUserUpgradeToChirpyRed(w http.ResponseWriter, r *http.Request) {

	type response struct {
		User
	}

	type data struct {
		UserID uuid.UUID `json:"user_id"`
	}

	type outerRequest struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := outerRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
