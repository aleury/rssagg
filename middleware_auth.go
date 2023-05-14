package main

import (
	"fmt"
	"net/http"

	"github.com/aleury/rssagg/internal/auth"
	"github.com/aleury/rssagg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (app *application) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("auth error: %s", err))
			return
		}

		user, err := app.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("user not found: %s", err))
			return
		}

		handler(w, r, user)
	}
}
