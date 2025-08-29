package handlers

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/samothreesixty/rss-agg/internal/db"
	"github.com/samothreesixty/rss-agg/internal/models"
)

func (apiConfig *ApiConfig) HandlerCreateFeed(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Invalid request payload")
		return
	}

	feed, err := apiConfig.DB.CreateFeed(r.Context(), db.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: 	   params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 500, "Cannot create feed")
		return
	}

	respondWithJson(w, 201, models.DatabaseFeedToFeed(feed))
}

func (apiConfig *ApiConfig) HandlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiConfig.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 500, "Cannot get feeds")
		return
	}

	feedModels := make([]models.Feed, len(feeds))
	for i, feed := range feeds {
		feedModels[i] = models.DatabaseFeedToFeed(feed)
	}

	respondWithJson(w, 200, feedModels)
}