package handlers

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/samothreesixty/rss-agg/internal/db"
	"github.com/samothreesixty/rss-agg/internal/models"
	"fmt"
)

func (apiConfig *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body);
	params := &parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Invalid request payload")
		return
	}

	user, err := apiConfig.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})
	if err != nil {
		respondWithError(w, 500, "Cannot create user")
		return
	}

	respondWithJson(w, 201, models.DatabaseUserToUser(user))
}

func (apiConfig *ApiConfig) HandlerGetUser(w http.ResponseWriter, r *http.Request, user db.User) {
	respondWithJson(w, 200, models.DatabaseUserToUser(user))
}

func (apiConfig *ApiConfig) HandlerGetUserPosts(w http.ResponseWriter, r *http.Request, user db.User) {
	posts, err := apiConfig.DB.GetPostsForUser(r.Context(), db.GetPostsForUserParams{
		user.ID,
		10,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get posts: %w", err))
		return
	}

	respondWithJson(w, 200, models.DatabasePostsToPosts(posts))
}