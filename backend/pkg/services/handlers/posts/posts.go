package posts

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/helpers"
)

func Posts(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	s, err := sqlite.CheckSession(r)
	if err != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}

	posts, err := sqlite.GetPosts(s.User.Id)
	if err != nil {
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(posts)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}
