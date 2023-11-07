package auth

import (
	"errors"
	"net/http"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/helpers"
	"social-network/pkg/models"

	uuid "github.com/gofrs/uuid"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	if r.Method == http.MethodPost {
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
		var credential string
		if u.Email == "" {
			credential = u.Username
		} else {
			credential = u.Email
		}
		id, err := sqlite.Authenticate(credential, u.Password)
		if err != nil {
			var errMsg helpers.ErrorMsg
			if errors.Is(err, models.ErrInvalidCredentials) {
				errMsg.ErrorDescription = "Email/username and password don't match."
				errMsg.ErrorType = "CREDENTIALS_DONT_MATCH"
				helpers.ErrorResponse(w, errMsg, http.StatusBadRequest)
			} else {
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			}
			return
		}
		err = sqlite.SetUserStatusOnline(id)
		if err != nil {
			helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			return
		}

		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:   "session",
			Value:  sID.String(),
			MaxAge: 60 * 60 * 24,
		}
		http.SetCookie(w, c)
		err = sqlite.InsertSession(c.Value, id)
		if err != nil {
			helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			return
		}
	}
}
