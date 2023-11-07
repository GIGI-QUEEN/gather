package services

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/helpers"
	"strconv"
	"strings"
)

func UserMessages(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	// Call GetMessages function with offset value
	s, err := sqlite.CheckSession(r)
	if err != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}

	url := strings.Split(strings.Trim(r.URL.Path, "/"), "/") // /chat/1/ -> [chat, 1]

	if len(url) == 1 && url[0] == "chat" {
		getAllUserChats(w, r)
		return
	}
	if len(url) == 2 {
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
			return
		}

		chatPersonId, err := strconv.Atoi(url[1])
		if err != nil {
			helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
			return
		}
		// show chat window with userId
		getChatWithUser(w, r, s.User.Id, chatPersonId, offset)
		return
	}
}

func getChatWithUser(w http.ResponseWriter, r *http.Request, userId, chatPersonId int, offset int) {
	helpers.EnableCors(&w)

	if chatPersonId < 1 {
		helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
		return
	}
	messages, err := sqlite.GetMessages(userId, chatPersonId, offset)
	if err != nil {
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(messages)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func getAllUserChats(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}

	s, err := sqlite.CheckSession(r)
	if err != nil {
		helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
		return
	}

	startedChatsLatestMessages, err := sqlite.GetStartedChatsLatestMessages(s.User.Id)
	if err != nil {
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(startedChatsLatestMessages)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetGroupChatMessages(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method == http.MethodGet {
		_, err := sqlite.CheckSession(r)
		if err != nil {
			helpers.ErrorResponse(w, helpers.UnauthorizedErrorMsg, http.StatusUnauthorized)
			return
		}

		url := strings.Split(strings.Trim(r.URL.Path, "/"), "/") // /group-chat/1/ -> [group-chat, 1]
		if len(url) == 2 {
			offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
			if err != nil {
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
				return
			}
			groupId, err := strconv.Atoi(url[1])
			if err != nil {
				helpers.ErrorResponse(w, helpers.NotFoundErrorMsg, http.StatusBadRequest)
				return
			}

			messages, err := sqlite.GetGroupMessages(groupId, offset)
			if err != nil {
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			jsonResp, err := json.Marshal(messages)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
				return
			}
			w.Write(jsonResp)
		}
	}

}
