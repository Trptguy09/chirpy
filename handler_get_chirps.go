package main

import (
	"database/sql"
	"errors"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	var respChirps []chirpResponse
	var authorID uuid.UUID
	var err error

	chirps, err := cfg.db.RetriveChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt retrive chirps", err)
		return
	}

	a := r.URL.Query().Get("author_id")
	if a != "" {
		authorID, err = uuid.Parse(a)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "unable to parse author_id", err)
			return
		}
	}

	for _, c := range chirps {

		if authorID != uuid.Nil && authorID != c.UserID {
			continue
		}
		respChirps = append(respChirps,
			chirpResponse{
				ID:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				Body:      c.Body,
				UserID:    c.UserID,
			})
	}
	s := r.URL.Query().Get("sort")
	if s == "asc" || s == "" {
		sort.Slice(respChirps, func(i, j int) bool {
			return respChirps[i].CreatedAt.Before(respChirps[j].CreatedAt)
		})
	}
	if s == "desc" {
		sort.Slice(respChirps, func(i, j int) bool {
			return respChirps[i].CreatedAt.After(respChirps[j].CreatedAt)
		})
	}
	respondWithJSON(w, http.StatusOK, respChirps)
}

func (cfg *apiConfig) handlerGetSingleChirp(w http.ResponseWriter, r *http.Request) {

	chirpIDStr := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id", err)
	}

	chirp, err := cfg.db.RetriveSingleChirp(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "404", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "500", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
