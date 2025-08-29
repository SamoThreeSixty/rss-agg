package handlers

import (
	"net/http"
	"github.com/samothreesixty/rss-agg/internal/auth"
	"github.com/samothreesixty/rss-agg/internal/db"
	"github.com/samothreesixty/rss-agg/internal/handlers/utils"
)

type ApiConfig struct {
	DB *db.Queries
}

type authedHandler func(http.ResponseWriter, *http.Request, db.User)

func (apiConfig *ApiConfig) MiddlewareAuth(next authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.ExtractAPIKey(r.Header)	
		if apiKey == "" {
			utils.RespondWithError(w, 403, "Invalid API key")
			return
		}
		user, err := apiConfig.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			utils.RespondWithError(w, 403, "Invalid API key")
			return
		}
		next(w, r, user)
	})
}