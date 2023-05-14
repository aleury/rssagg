package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aleury/rssagg/internal/database"
	"github.com/google/uuid"
)

func (app *application) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := app.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't get feeds: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}

func (app *application) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing json: %s", err))
		return
	}

	feed, err := app.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't create feed: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedToFeed(feed))
}
