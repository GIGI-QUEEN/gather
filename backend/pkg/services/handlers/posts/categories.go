package posts

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/helpers"
	"social-network/pkg/models"
	"strconv"
)

func Categories(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	_, err := sqlite.CheckSession(r)
	if err != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}
	categories, err := sqlite.GetAllCategories()
	if err != nil {
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(categories)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func PostsByCategoryId(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	_, err := sqlite.CheckSession(r)
	if err != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}
	catId, err := strconv.Atoi(r.URL.Path[10:])
	if err != nil || catId < 1 {
		helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
	p, err := sqlite.GetPostsByCategory(catId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			helpers.ErrorResponse(w, helpers.EmptyCategoryErrorMsg, http.StatusBadRequest)
		} else {
			helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}
