package groups

import (
	"errors"
	"log"
	"net/http"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/helpers"
	"social-network/pkg/models"
)

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	s, err := sqlite.CheckSession(r)
	if err != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodPost {
		var group models.Group

		err := helpers.DecodeJSONBody(w, r, &group)
		if err != nil {
			var errMsg *helpers.ErrorMsg
			if errors.As(err, &errMsg) {
				helpers.ErrorResponse(w, *errMsg, http.StatusBadRequest)
			} else {
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			}
			return
		}
		if helpers.ValidateUserGroupCreationData(w, group) {
			groupId, err := sqlite.InsertGroup(group, s.User.Id)
			if err != nil {
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			}

			err = sqlite.InsertGroupConversationRoom(groupId, s.User.Id)
			if err != nil {
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			}

			err = sqlite.InsertUserAsAdminToGroup(groupId, s.User.Id)
			if err != nil {
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			}
			log.Printf("UserId[%d] created groupId[%d]", s.User.Id, groupId)
		}
		return
	}
	helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
}
