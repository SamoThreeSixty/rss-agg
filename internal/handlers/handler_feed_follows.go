package handlers

import (
	"net/http"
	"encoding/json"
	"time"
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/samothreesixty/rss-agg/internal/db"
	"github.com/samothreesixty/rss-agg/internal/models"
	"github.com/samothreesixty/rss-agg/internal/handlers/utils"
)

func (apiConfig *ApiConfig) HandlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 400, "Invalid request payload")
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
		utils.RespondWithError(w, 500, "Cannot create feed")
		return
	}

	utils.RespondWithJson(w, 201, models.DatabaseFeedFollowToFeedFollow(feedFollow))
}

func (apiConfig *ApiConfig) HandlerGetFeedFollowsByUser(w http.ResponseWriter, r *http.Request, user db.User) {
	feedFollows, err := apiConfig.DB.GetFeedFollowsByUserID(r.Context(), user.ID)
	if err != nil {
		utils.RespondWithError(w, 500, "Cannot get feed follows")
		return
	}

	feedFollowModels := make([]models.FeedFollow, len(feedFollows))
	for i, feedFollow := range feedFollows {
		feedFollowModels[i] = models.DatabaseFeedFollowToFeedFollow(feedFollow)
	}

	utils.RespondWithJson(w, 200, feedFollowModels)
}

func (apiConfig *ApiConfig) HandlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user db.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowId")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		utils.RespondWithError(w, 400, "Invalid feed follow ID")
		return
	}

	deletedID, err := apiConfig.DB.DeleteFeedFollow(r.Context(), db.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err == sql.ErrNoRows {
		utils.RespondWithError(w, 404, "Feed follow not found")
		return
	}
	if err != nil {
		utils.RespondWithError(w, 500, "Cannot delete feed follow")
		return
	}

	utils.RespondWithJson(w, 200, map[string]string{
		"result": "success",
		"id":    deletedID.String(),
	})
}