package groups

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/helpers"
)

func ShowGroupList(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	_, err := sqlite.CheckSession(r)
	if err != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}

	groups, err := sqlite.GetGroupsList()
	if err != nil {
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(groups)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}
