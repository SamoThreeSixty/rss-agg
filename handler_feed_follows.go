package main

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/samothreesixty/rss-agg/internal/db"
)

func (apiConfig *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Invalid request payload")
		return
	}

	feedFollow, err := apiConfig.DB.CreateFeedFollow(r.Context(), db.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 500, "Cannot create feed")
		return
	}

	respondWithJson(w, 201, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiConfig *apiConfig) handlerGetFeedFollowsByUser(w http.ResponseWriter, r *http.Request, user db.User) {
	feedFollows, err := apiConfig.DB.GetFeedFollowsByUserID(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 500, "Cannot get feed follows")
		return
	}

	feedFollowModels := make([]FeedFollow, len(feedFollows))
	for i, feedFollow := range feedFollows {
		feedFollowModels[i] = databaseFeedFollowToFeedFollow(feedFollow)
	}

	respondWithJson(w, 200, feedFollowModels)
}