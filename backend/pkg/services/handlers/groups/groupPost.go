package groups

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/helpers"
	"social-network/pkg/models"
	"strconv"
	"strings"
)

func GroupPostById(w http.ResponseWriter, r *http.Request, postId int) {
	helpers.EnableCors(&w)
	if postId < 1 {
		helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
	post, err := sqlite.GetGroupPostById(postId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
		} else {
			helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(post)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GroupPost(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}

	_, authErr := sqlite.CheckSession(r)
	if authErr != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}
	url := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(url) == 1 {
		helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
	postId, err := strconv.Atoi(url[2])
	if err != nil {
		helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
	if len(url) == 3 { //  /group/post/5
		GroupPostById(w, r, postId)
		return
	}

}
