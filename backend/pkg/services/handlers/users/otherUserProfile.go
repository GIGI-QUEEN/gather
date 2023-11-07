package users

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

func GetProfile(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}

	s, authErr := sqlite.CheckSession(r)
	if authErr != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}

	url := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	switch len(url) {
	case 2:
		userId, err := convertIdFromUrl(url[1])
		if err != nil {
			helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
			return
		}

		getUserProfileById(w, r, userId)
		return
	case 3:
		if r.Method == http.MethodPost {
			userToFollow, err := convertIdFromUrl(url[1])
			if err != nil {
				helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
				return
			}
			if url[2] == "follow" {
				sqlite.FollowUser_v2(userToFollow, s.User.Id)
				return
			}
			if url[2] == "unfollow" {
				sqlite.UnfollowUser(userToFollow, s.User.Id)
				return
			}
		}

	default:
		helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
}

func getUserProfileById(w http.ResponseWriter, r *http.Request, userProfileId int) {
	helpers.EnableCors(&w)
	s, err := sqlite.CheckSession(r)
	if err != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}

	var u *models.User
	u, err = sqlite.GetUserProfile(userProfileId, s.User.Id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			var errMsg helpers.ErrorMsg
			errMsg.ErrorDescription = "User not found"
			errMsg.ErrorType = "STATUS_BAD_REQUEST"
			helpers.ErrorResponse(w, errMsg, http.StatusBadRequest)
		} else {
			helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(u)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func convertIdFromUrl(idAsStringFromUrl string) (int, error) {
	groupId, err := strconv.Atoi(idAsStringFromUrl)
	if err != nil {
		return 0, err
	}
	if groupId < 1 {
		return 0, err
	}
	return groupId, nil
}
