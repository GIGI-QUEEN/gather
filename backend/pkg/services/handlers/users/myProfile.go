package users

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/helpers"
	"social-network/pkg/models"
	"social-network/pkg/services/handlers/images"
	"strconv"
	"strings"
)

func GetMyProfile(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)

	s, err := sqlite.CheckSession(r)
	if err != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}

	var u *models.User
	u, err = sqlite.GetUserProfile(s.User.Id, s.User.Id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
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

func ChangeAccountAvatar(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	s, err := sqlite.CheckSession(r)
	if err != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}
	r.ParseMultipartForm(10 << 20)

	fileType := r.Form.Get("file_type")
	path := `/images/avatars/` + strconv.Itoa(s.User.Id) + "." + fileType
	images.UploadFile(r, s.User.Id)
	sqlite.InsertUserAvatar(s.User.Id, path)
}

func AccountSettings(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	s, err := sqlite.CheckSession(r)
	if err != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}
	url := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(url) == 1 {
		helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
	if len(url) == 3 {

		if r.Method == http.MethodPost {

			if (url[2]) == "privacy" || url[2] == "about" || url[2] == "username" {

				var u models.User
				err := helpers.DecodeJSONBody(w, r, &u)
				if err != nil {
					var errMsg *helpers.ErrorMsg
					if errors.As(err, &errMsg) {
						helpers.ErrorResponse(w, *errMsg, http.StatusBadRequest)
					} else {
						helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
					}
					return
				}
				if (url[2]) == "privacy" {
					sqlite.ChangeAccountPrivacy(s.User.Id, u.Privacy)
					return
				}

				if (url[2]) == "about" {
					sqlite.ChangeAccountAbout(s.User.Id, u.About)
					return
				}

				if (url[2]) == "username" {
					sqlite.ChangeUsername(s.User.Id, u.Username)
					return
				}
			}

			if url[2] == "accept-follow" || url[2] == "reject-follow" {
				var follower models.Follower
				err := helpers.DecodeJSONBody(w, r, &follower)

				if err != nil {
					var errMsg *helpers.ErrorMsg
					if errors.As(err, &errMsg) {
						helpers.ErrorResponse(w, *errMsg, http.StatusBadRequest)
					} else {
						helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
					}
					return
				}
				switch url[2] {
				case "accept-follow":
					sqlite.AcceptFollowRequest(follower.Follower.Id, s.User.Id)
					return
				case "reject-follow":
					sqlite.RejectFollowRequst(follower.Follower.Id, s.User.Id)
					return
				}
			}

		}
	}
}