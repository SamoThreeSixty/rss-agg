package handlers

import "net/http"

func (apiCfg *ApiConfig) HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, 200, struct{}{})
}