package groups

import (
	"database/sql"
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

func Group(w http.ResponseWriter, r *http.Request) {
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
		groupId, err := convertIdFromUrl(url[1])
		if err != nil {
			helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			showGroupById(w, r, groupId, s.User.Id)
			return
		case http.MethodPost:
			createGroupPost(w, r, groupId, s.User.Id)
			return
		}
	case 3:
		if r.Method != http.MethodPost {
			helpers.ErrorResponse(w, helpers.MethodNotAllowedMsg, http.StatusMethodNotAllowed)
			return
		} else {
			groupId, err := convertIdFromUrl(url[1])
			if err != nil {
				helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
				return
			}
			switch url[2] {
			case "join":
				requestJoin(w, r, groupId, s.User.Id)
				return
			case "approve":
				approveJoiningRequest(w, r, groupId, s.User.Id)
				return
			case "reject":
				rejectJoiningRequest(w, r, groupId, s.User.Id)
			case "leave":
				leaveTheGroup(w, r, groupId, s.User.Id)
				return
			case "invite":
				inviteUserToGroup(w, r, groupId, s.User.Id)
				return
			case "accept-invite":
				acceptInvite(w, r, groupId, s.User.Id)
				return
			case "reject-invite":
				rejectInvite(w, r, groupId, s.User.Id)
				return
			case "create-event":
				createGroupEvent(w, r, groupId, s.User.Id)
				return
			default:
				helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
			}
		}
	case 4:
		if r.Method != http.MethodPut {
			helpers.ErrorResponse(w, helpers.MethodNotAllowedMsg, http.StatusMethodNotAllowed)
			return
		} else {
			if url[0] != "group" && url[2] != "event" {
				helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
				return
			}
			_, eventId, err := convertIdsFromUrl(url)
			if err != nil {
				helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
				return
			}
			deleteNotifAboutEvent(w, r, eventId, s.User.Id)
		}
	case 5:
		switch url[2] {
		case "post":
			if r.Method != http.MethodPost {
				helpers.ErrorResponse(w, helpers.MethodNotAllowedMsg, http.StatusMethodNotAllowed)
				return
			} else {
				groupId, postId, err := convertIdsFromUrl(url)
				if err != nil {
					helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
					return
				}
				switch url[4] {
				case "comment":
					commentOnGroupPost(w, r, groupId, postId, s.User.Id)
				case "like":
					likeOnGroupPost(w, r, groupId, postId, s.User.Id)
				case "dislike":
					dislikeOnGroupPost(w, r, groupId, postId, s.User.Id)
				default:
					helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
				}

			}
		case "event":
			if r.Method != http.MethodPost {
				helpers.ErrorResponse(w, helpers.MethodNotAllowedMsg, http.StatusMethodNotAllowed)
				return
			} else {
				_, eventId, err := convertIdsFromUrl(url) // same function can be used here
				if err != nil {
					helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
					return
				}
				switch url[4] {
				case "event-accept":
					acceptGroupEvent(w, r, eventId, s.User.Id)
				case "event-reject":
					rejectGroupEvent(w, r, eventId, s.User.Id)
				default:
					helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
				}
			}
		}

	case 7:
		if r.Method != http.MethodPost {
			helpers.ErrorResponse(w, helpers.MethodNotAllowedMsg, http.StatusMethodNotAllowed)
			return
		} else {
			commentId, err := strconv.Atoi(url[5])
			if err != nil {
				helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
				return
			}
			switch url[6] {
			case "like":
				likeOnGroupPostComment(w, r, commentId, s.User.Id)
			case "dislike":
				dislikeOnGroupPostComment(w, r, commentId, s.User.Id)
			default:
				helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
			}
		}
	default:
		helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
	}
}

func showGroupById(w http.ResponseWriter, r *http.Request, groupId int, userId int) {
	helpers.EnableCors(&w)

	group, err := sqlite.GetGroupById(groupId, userId)
	if err != nil {

		if errors.Is(err, models.ErrNoRecord) {
			helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
		} else {
			helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(group)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func createGroupPost(w http.ResponseWriter, r *http.Request, groupId int, userId int) {
	helpers.EnableCors(&w)

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

	// check if user isApproved in group_users to be able to create posts
	isApproved, err := sqlite.CheckUserIsApproved(groupId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.ErrorResponse(w, helpers.ForbiddenErrorMsg, http.StatusForbidden)
		} else {
			helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		}
		return
	}

	if isApproved == 1 {
		var groupPost models.GroupPost

		r.ParseMultipartForm(10 << 20)
		var title string
		var content string
		var fileType string

		for key, value := range r.Form {
			switch key {
			case "title":
				title = value[0]
			case "content":
				content = value[0]
			case "file_type":
				fileType = value[0]
			}
		}
		groupPost.Title = title
		groupPost.Content = content
		if helpers.ValidateUserGroupPostData_v1(w, groupPost) {
			postId, err := sqlite.InsertGroupPost(groupPost, groupId, user.Id)
			if err != nil {
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
				return
			}
			if fileType != "undefined" {
				images.UploadFile(r, postId)
				path := `/images/group-posts-images/` + strconv.Itoa(postId) + "." + fileType
				sqlite.InsertGroupPostImagePath(path, postId)
			}
		}

	} else {
		helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
	}
}

func requestJoin(w http.ResponseWriter, r *http.Request, groupId int, userId int) {
	helpers.EnableCors(&w)

	if err := sqlite.InsertJoinRequest(groupId, userId); err != nil {
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		return
	} else {
		showGroupById(w, r, groupId, userId)
		return
	}
}

func approveJoiningRequest(w http.ResponseWriter, r *http.Request, groupId int, userId int) {
	helpers.EnableCors(&w)

	if r.Method == http.MethodPost {
		adminId := userId
		if sqlite.CheckUserIsAdmin(groupId, adminId) {
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
			if err := sqlite.InsertNewUserToGroup(groupId, u.Id); err != nil {
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			}
		} else {
			helpers.ErrorResponse(w, helpers.ForbiddenErrorMsg, http.StatusForbidden)
		}
	}
}

func rejectJoiningRequest(w http.ResponseWriter, r *http.Request, groupId int, userId int) {
	helpers.EnableCors(&w)
	if r.Method == http.MethodPost {
		admin := userId
		if sqlite.CheckUserIsAdmin(groupId, admin) {
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

			if err := sqlite.RejectUserJoinRequest(groupId, u.Id); err != nil {
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			}
		} else {
			helpers.ErrorResponse(w, helpers.ForbiddenErrorMsg, http.StatusForbidden)
		}

	}
}

func leaveTheGroup(w http.ResponseWriter, r *http.Request, groupId int, userId int) {
	helpers.EnableCors(&w)

	if err := sqlite.RemoveUserFromGroup(groupId, userId); err != nil {
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		return
	} else {
		showGroupById(w, r, groupId, userId)
		return
	}
}

func inviteUserToGroup(w http.ResponseWriter, r *http.Request, groupId int, userId int) {
	helpers.EnableCors(&w)

	if sqlite.CheckUserIsMember(groupId, userId) {
		var userToInvite models.User
		err := helpers.DecodeJSONBody(w, r, &userToInvite)
		if err != nil {
			var errMsg *helpers.ErrorMsg
			if errors.As(err, &errMsg) {
				helpers.ErrorResponse(w, *errMsg, http.StatusBadRequest)
			} else {
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			}
			return
		}
		if err := sqlite.InviteNewUserToGroup(groupId, userToInvite.Id, userId); err != nil {
			helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		}
	} else {
		helpers.ErrorResponse(w, helpers.ForbiddenErrorMsg, http.StatusForbidden)
	}
}

func acceptInvite(w http.ResponseWriter, r *http.Request, groupId int, userId int) {
	helpers.EnableCors(&w)

	if err := sqlite.AcceptInviteToGroup(groupId, userId); err != nil {
	}
}

func rejectInvite(w http.ResponseWriter, r *http.Request, groupId int, userId int) {
	helpers.EnableCors(&w)

	if err := sqlite.RejectInviteToGroup(groupId, userId); err != nil {
	}
}

func createGroupEvent(w http.ResponseWriter, r *http.Request, groupId int, userId int) {
	helpers.EnableCors(&w)

	s, authErr := sqlite.CheckSession(r)
	if authErr != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}

	// check if user isApproved in group_users to be able to create posts
	isApproved, err := sqlite.CheckUserIsApproved(groupId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.ErrorResponse(w, helpers.ForbiddenErrorMsg, http.StatusForbidden)
		} else {
			helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		}
		return
	}
	if isApproved == 1 {
		var groupEvent models.GroupEvent

		err := helpers.DecodeJSONBody(w, r, &groupEvent)
		if err != nil {
			var errMsg *helpers.ErrorMsg
			if errors.As(err, &errMsg) {
				helpers.ErrorResponse(w, *errMsg, http.StatusBadRequest)
			} else {
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			}
			return
		}
		if helpers.ValidateGroupEventCreationData(w, groupEvent) {
			err := sqlite.InsertGroupEvent(groupId, groupEvent, s.User.Id)
			if err != nil {
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			}
		}
		return
	} else {
		helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
	}
}

func acceptGroupEvent(w http.ResponseWriter, r *http.Request, eventId int, userId int) {
	helpers.EnableCors(&w)

	if err := sqlite.AcceptInviteToGroupEvent(eventId, userId); err != nil {
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
	}
}

func rejectGroupEvent(w http.ResponseWriter, r *http.Request, eventId int, userId int) {
	helpers.EnableCors(&w)

	if err := sqlite.RejectInviteToGroupEvent(eventId, userId); err != nil {
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
	}
}

func deleteNotifAboutEvent(w http.ResponseWriter, r *http.Request, eventId int, userId int) {
	helpers.EnableCors(&w)

	if err := sqlite.DeleteNotificationAboutCreatedGroupEvent(eventId, userId); err != nil {
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
	}
}

func commentOnGroupPost(w http.ResponseWriter, r *http.Request, groupId int, postId int, userId int) {
	helpers.EnableCors(&w)

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
	commentId, _ := sqlite.InsertGroupPostComment(postId, userId, content)
	if fileType != "undefined" {
		images.UploadFile(r, commentId)
		path := "/images/group-posts-comments-images/" + strconv.Itoa(commentId) + "." + fileType
		sqlite.InsertGroupPostCommentImagePath(commentId, path)
	}
}

func likeOnGroupPost(w http.ResponseWriter, r *http.Request, groupId int, postId int, userId int) {
	helpers.EnableCors(&w)

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

	if err = sqlite.ChangeGroupPostLike(postId, userId); err != nil {
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
	}
}

func dislikeOnGroupPost(w http.ResponseWriter, r *http.Request, groupId int, postId int, userId int) {
	helpers.EnableCors(&w)

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

	if err = sqlite.ChangeGroupPostDislike(postId, userId); err != nil {
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
	}
}

func likeOnGroupPostComment(w http.ResponseWriter, r *http.Request, commentId int, userId int) {
	helpers.EnableCors(&w)

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

		if err = sqlite.ChangeGroupPostCommentLike(commentId, userId); err != nil {
			helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)

		}
	}
}

func dislikeOnGroupPostComment(w http.ResponseWriter, r *http.Request, commentId int, userId int) {
	helpers.EnableCors(&w)

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

		if err = sqlite.ChangeGroupPostCommentDislike(commentId, userId); err != nil {
			helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)

		}
	}
}

func convertIdsFromUrl(url []string) (int, int, error) {
	firstId, err := strconv.Atoi(url[1])
	if err != nil {
		return -1, -1, nil
	}
	secondId, err := strconv.Atoi(url[3])
	if err != nil {
		return -1, -1, nil
	}

	return firstId, secondId, nil
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
