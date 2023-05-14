package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aleury/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (app *application) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := app.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("could get feed follows: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (app *application) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing json: %s", err))
		return
	}

	feedFollow, err := app.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't create feed follow: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedFollowToFeedFollow(feedFollow))
}

func (app *application) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't parse feed follow id: %s", err))
		return
	}

	err = app.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't delete feed follow: %s", err))
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
