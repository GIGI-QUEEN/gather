package notifications

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/helpers"
	"social-network/pkg/models"
)

func UserNotifications(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	if r.Method == http.MethodGet {
		s, err := sqlite.CheckSession(r)
		if err != nil {
			helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
			return
		}
		notif := &models.Notifications{}
		notif.CommentOnUsersPost, err = sqlite.CheckForUnreadCommentsOnMyPostInGroup(s.User.Id)
		if err = checkUsersWantToFollow(notif, s.User.Id); err != nil {
			log.Println("ERROR IN UserNotifications in checkUsersWantToFollow() : ", err)
		}

		if err = checkGroupJoinRequests(notif, s.User.Id); err != nil {
			log.Println("ERROR IN UserNotifications in checkGroupJoinRequests() : ", err)
		}

		if err = checkReceivedGroupInvites(notif, s.User.Id); err != nil {
			log.Println("ERROR IN UserNotifications in checkReceivedGroupInvites() : ", err)
		}
		if err = checkGroupEventCreated(notif, s.User.Id); err != nil {
			log.Println("ERROR IN UserNotifications in checkGroupEventCreated() : ", err)
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResp, err := json.Marshal(notif)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			return
		}
		w.Write(jsonResp)
	} else {
		helpers.ErrorResponse(w, helpers.MethodNotAllowedMsg, http.StatusMethodNotAllowed)
	}
}

func checkUsersWantToFollow(notif *models.Notifications, userId int) error {
	var err error
	notif.UsersWantToFollow, err = sqlite.GetUsersWantToFollow(userId)
	if err != nil {
		return err
	}
	return nil
}

func checkGroupJoinRequests(notif *models.Notifications, userId int) error {
	groupIds := sqlite.GetGroupIdsWhereUserIsAdmin(userId)

	if len(groupIds) > 0 {
		groupJoinRequests := make([]*models.GroupJoinRequest, 0)
		for _, groupId := range groupIds {
			users, err := sqlite.GetPossibleJoinRequests(groupId)
			if err != nil {
				return err
			}
			groupInfo, err := sqlite.GetGroupInfoById(groupId)
			if err != nil {
				return err
			}

			if len(users) > 0 {
				groupJoinRequests = append(groupJoinRequests, &models.GroupJoinRequest{
					Group:              groupInfo,
					UsersRequestedJoin: users,
				})
			}
		}
		notif.GroupJoinRequests = groupJoinRequests
	}
	return nil
}

func checkReceivedGroupInvites(notif *models.Notifications, userId int) error {
	var err error
	notif.GroupJoinInvites, err = sqlite.GetPossibleJoinInvites(userId)
	if err != nil {
		return err
	}
	return nil
}

func checkGroupEventCreated(notif *models.Notifications, userId int) error {
	var err error
	notif.GroupEventsCreated, err = sqlite.GetPossibleCreatedGroupEvents(userId)
	if err != nil {
		return err
	}
	return nil
}

func ClearAllGroupPostCommentsNotifications(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	if r.Method == http.MethodPut {
		_, err := sqlite.CheckSession(r)
		if err != nil {
			helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
			return
		}
		var commentsIds models.GroupPostCommentIds
		err = helpers.DecodeJSONBody(w, r, &commentsIds)
		if err != nil {
			var errMsg *helpers.ErrorMsg
			if errors.As(err, &errMsg) {
				helpers.ErrorResponse(w, *errMsg, http.StatusBadRequest)
			} else {
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			}
			return
		}
		if len(commentsIds.CommentId) > 0 {
			err = sqlite.MarkAllUnreadCommentsAsRead(commentsIds.CommentId)
			if err != nil {
				if errors.Is(err, models.ErrNoRecord) {
					helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
				} else {
					helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
				}
				return
			}
		}
	}
}
