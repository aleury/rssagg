package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aleury/rssagg/internal/database"
	"github.com/google/uuid"
)

func (app *application) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (app *application) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing json: %s", err))
		return
	}

	user, err := app.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't create user: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
