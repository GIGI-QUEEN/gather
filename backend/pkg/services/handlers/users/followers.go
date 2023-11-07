package users

import (
	"net/http"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/helpers"
	"strconv"
)

func Followers(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	_, err := sqlite.CheckSession(r)
	if err != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}
	userId, err := strconv.Atoi(r.URL.Path[6:])
	if err != nil || userId < 1 {
		helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
}
