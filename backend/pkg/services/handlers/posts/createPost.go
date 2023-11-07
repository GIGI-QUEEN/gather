package posts

import (
	"net/http"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/helpers"
	"social-network/pkg/services/handlers/images"
	"strconv"
	"strings"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}

	s, err := sqlite.CheckSession(r)
	if err != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}
	r.ParseMultipartForm(10 << 20)

	var title string
	var content string
	var privacy string
	var categories string
	var allowedFollowers string
	var fileType string

	for key, value := range r.Form {
		switch key {
		case "title":
			title = value[0]
		case "content":
			content = value[0]
		case "privacy":
			privacy = value[0]
		case "categories":
			categories = value[0]
		case "allowed":
			allowedFollowers = value[0]
		case "file_type":
			fileType = value[0]
		}
	}
	postCategories := strings.Split(categories, ",")
	postUsers := strings.Split(allowedFollowers, ",")
	if helpers.ValidateUserPostData(w, title, content) {
		id, err := sqlite.InsertPost(title, content, postCategories, s.User.Id, privacy, postUsers)
		if err != nil {
			helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		}
		if fileType != "undefined" {
			images.UploadFile(r, id)
			path := `/images/posts-images/` + strconv.Itoa(id) + "." + fileType
			err = sqlite.InsertPostImagePath(id, path)
			if err != nil {
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			}
		}
	}

}
