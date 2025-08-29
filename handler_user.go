package main

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/samothreesixty/rss-agg/internal/db"
	"fmt"
)

func (apiConfig *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
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

	respondWithJson(w, 201, databaseUserToUser(user))
}

func (apiConfig *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user db.User) {
	respondWithJson(w, 200, databaseUserToUser(user))
}

func (apiConfig *apiConfig) handlerGetUserPosts(w http.ResponseWriter, r *http.Request, user db.User) {
	posts, err := apiConfig.DB.GetPostsForUser(r.Context(), db.GetPostsForUserParams{
		user.ID,
		10,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get posts: %w", err))
		return
	}

	respondWithJson(w, 200, databasePostsToPosts(posts))
}