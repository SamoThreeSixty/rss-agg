package handlers

import (
	"net/http"
	"github.com/samothreesixty/rss-agg/internal/handlers/utils"
)

func (apiCfg *ApiConfig) HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJson(w, 200, struct{}{})
}