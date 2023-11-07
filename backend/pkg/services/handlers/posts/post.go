package posts

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

func PostById(w http.ResponseWriter, r *http.Request, postId int) {
	helpers.EnableCors(&w)
	if postId < 1 {
		helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
	post, err := sqlite.GetPostById(postId)
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

func Post(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	user := &models.User{
		Id: -1,
	}
	s, authErr := sqlite.CheckSession(r)
	if authErr != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	} else {
		user = s.User
	}
	url := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(url) == 1 {
		helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
	postId, err := strconv.Atoi(url[1])
	if err != nil {
		helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
		return
	}

	if len(url) == 2 { // /post/%id%/
		PostById(w, r, postId)
		return
	}
	if len(url) == 3 { // /post/%id%/%something%

		
		if url[2] == "comment" {
			r.ParseMultipartForm(10 << 20)

			var content string
			var fileType string

			for key, value := range r.Form {

				switch key {

				case "content":
					content = value[0]
				case "file_type":
					fileType = value[0]

				}

			}
			commentId, _ := sqlite.InsertComment(postId, content, user.Id)
			if fileType != "undefined" {
				images.UploadFile(r, commentId)
				path := `/images/comments-images/` + strconv.Itoa(commentId) + "." + fileType
				sqlite.InsertCommentImagePath(commentId, path)
			}
			return
		}

		if r.Method == http.MethodPost {
			var pr models.PostReaction
			err := helpers.DecodeJSONBody(w, r, &pr)
			if err != nil {
				var errMsg *helpers.ErrorMsg
				if errors.As(err, &errMsg) {
					helpers.ErrorResponse(w, *errMsg, http.StatusBadRequest)
				} else {
					helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
				}
				return
			}
			if url[2] == "like" && pr.PostLikeDislike == "like" {
				sqlite.ChangePostLike(postId, user.Id)
				return
			}
			if url[2] == "dislike" && pr.PostLikeDislike == "dislike" {
				sqlite.ChangePostDislike(postId, user.Id)
				return
			}
		}
		helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
	}

	if len(url) == 5 {
		commentId, err := strconv.Atoi(url[3])
		if err != nil {
			helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
			return
		}
		if r.Method == http.MethodPost {
			var cr models.CommentReaction
			err := helpers.DecodeJSONBody(w, r, &cr)
			if err != nil {
				var errMsg *helpers.ErrorMsg
				if errors.As(err, &errMsg) {
					helpers.ErrorResponse(w, *errMsg, http.StatusBadRequest)
				} else {
					helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
				}
				return
			}
			if url[4] == "like" && cr.CommentLikeDislike == "like" {
				sqlite.ChangeCommentLike(commentId, user.Id)
				return
			}
			if url[4] == "dislike" && cr.CommentLikeDislike == "dislike" {
				sqlite.ChangeCommentDisLike(commentId, user.Id)
				return
			}
		}
	}
}
