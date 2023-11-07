package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"regexp"
	"social-network/pkg/models"
	"strings"
)

type ErrorMsg struct {
	ErrorDescription string `json:"error_description"`
	ErrorType        string `json:"error_type"`
}

func (mr *ErrorMsg) Error() string {
	return mr.ErrorDescription
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errDescription := "Content Type is not application/json"
		errType := "WRONG_CONTENCT_TYPE"
		return &ErrorMsg{ErrorDescription: errDescription, ErrorType: errType}
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&dst)
	if err != nil {
		// big error handling done here
		// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
		var unmarshalTypeError *json.UnmarshalTypeError
		switch {

		case errors.Is(err, io.EOF):
			errDescription := "Request body must not be empty."
			errType := "REQUEST_BODY_EMPTY"
			return &ErrorMsg{ErrorDescription: errDescription, ErrorType: errType}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			errDescription := msg
			errType := "INVALID_VALUE_FOR_FIELD"
			return &ErrorMsg{ErrorDescription: errDescription, ErrorType: errType}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			errDescription := "Request body contains unknown field " + fieldName
			errType := "UNKNOWN_FIELD"
			return &ErrorMsg{ErrorDescription: errDescription, ErrorType: errType}

		default:
			return err
		}

	}
	return nil
}

func ErrorResponse(w http.ResponseWriter, errMessage ErrorMsg, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	jsonResp, _ := json.Marshal(errMessage)
	w.Write(jsonResp)
}

func ValidateUserData(w http.ResponseWriter, firstname, lastname, email, username, gender, password string, age int) bool {
	var errMsg ErrorMsg
	space := regexp.MustCompile(`\s+`)
	firstname = space.ReplaceAllString(strings.TrimSpace(firstname), " ")
	lastname = space.ReplaceAllString(strings.TrimSpace(lastname), " ")
	username = space.ReplaceAllString(username, "")
	email = space.ReplaceAllString(strings.TrimSpace(email), "")

	if firstname == "" {
		errMsg.ErrorDescription = "Firstname is missing."
		errMsg.ErrorType = "FIRSTNAME_FIELD_EMPTY"
		ErrorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	if lastname == "" {
		errMsg.ErrorDescription = "Lastname is missing."
		errMsg.ErrorType = "LASTNAME_FIELD_EMPTY"
		ErrorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	/* if username == "" {
		errMsg.ErrorDescription = "Username is missing."
		errMsg.ErrorType = "USERNAME_FIELD_EMPTY"
		ErrorResponse(w, errMsg, http.StatusBadRequest)
		return false
	} */
	if email == "" {
		errMsg.ErrorDescription = "Email is missing."
		errMsg.ErrorType = "EMAIL_FIELD_EMPTY"
		ErrorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}

	if gender == "" {
		errMsg.ErrorDescription = "Gender is missing."
		errMsg.ErrorType = "GENDER_FIELD_EMPTY"
		ErrorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}

	if password == "" {
		errMsg.ErrorDescription = "Password is missing."
		errMsg.ErrorType = "PASSWORD_FIELD_EMPTY"
		ErrorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}

	if len(password) < 6 || len(password) > 20 {
		errMsg.ErrorDescription = "Password is too short - 6 chars min."
		errMsg.ErrorType = "PASSWORD_TOO_SHORT"
		ErrorResponse(w, errMsg, http.StatusNotAcceptable)
		return false
	}

	if age <= 0 || age > 120 {
		errMsg.ErrorDescription = "Age is not valid"
		errMsg.ErrorType = "AGE_NOT_VALID"
		ErrorResponse(w, errMsg, http.StatusNotAcceptable)
		return false
	}

	_, errMail := mail.ParseAddress(email)
	if errMail != nil {
		errMsg.ErrorDescription = "Email is not valid"
		errMsg.ErrorType = "EMAIL_INVALID"
		ErrorResponse(w, errMsg, http.StatusNotAcceptable)
		return false
	}
	return true

}

func ValidateUserPostData(w http.ResponseWriter, title, content string) bool {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" || content == "" {
		var errMsg ErrorMsg
		if title == "" {
			errMsg.ErrorDescription = "Title field is empty"
			errMsg.ErrorType = "TITLE_FIELD_EMPTY"

		} else if content == "" {
			errMsg.ErrorDescription = "Content field is empty"
			errMsg.ErrorType = "CONTENT_FIELD_EMPTY"

		}
		ErrorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	return true
}

func ValidateUserGroupPostData_v1(w http.ResponseWriter, post models.GroupPost) bool {
	title := strings.TrimSpace(post.Title)
	content := strings.TrimSpace(post.Content)

	if title == "" || content == "" {
		var errMsg ErrorMsg
		if title == "" {
			errMsg.ErrorDescription = "Title field is empty"
			errMsg.ErrorType = "TITLE_FIELD_EMPTY"

		} else if content == "" {
			errMsg.ErrorDescription = "Content field is empty"
			errMsg.ErrorType = "CONTENT_FIELD_EMPTY"

		}
		ErrorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	return true
}

func ValidateUserGroupCreationData(w http.ResponseWriter, group models.Group) bool {
	title := strings.TrimSpace(group.Title)
	content := strings.TrimSpace(group.Description)

	if title == "" || content == "" {
		var errMsg ErrorMsg
		if title == "" {
			errMsg.ErrorDescription = "Title field is empty"
			errMsg.ErrorType = "TITLE_FIELD_EMPTY"

		} else if content == "" {
			errMsg.ErrorDescription = "Content field is empty"
			errMsg.ErrorType = "CONTENT_FIELD_EMPTY"

		}
		ErrorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	return true
}

func ValidateGroupEventCreationData(w http.ResponseWriter, groupEvent models.GroupEvent) bool {
	title := strings.TrimSpace(groupEvent.Title)
	content := strings.TrimSpace(groupEvent.Description)

	if title == "" || content == "" {
		var errMsg ErrorMsg
		if title == "" {
			errMsg.ErrorDescription = "Title field is empty"
			errMsg.ErrorType = "TITLE_FIELD_EMPTY"

		} else if content == "" {
			errMsg.ErrorDescription = "Content field is empty"
			errMsg.ErrorType = "CONTENT_FIELD_EMPTY"

		}
		ErrorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	return true
}

func ValidateUserMessage(w http.ResponseWriter, content string) bool {
	content = strings.TrimSpace(content)
	if content == "" {
		var errMsg ErrorMsg
		errMsg.ErrorDescription = "Content field is empty"
		errMsg.ErrorType = "CONTENT_FIELD_EMPTY"
		ErrorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	return true
}

/* func ValidateUserMessage(w http.ResponseWriter, message models.Message) bool {
	content := strings.TrimSpace(message.Content)
	if content == "" {
		var errMsg ErrorMsg
		errMsg.ErrorDescription = "Content field is empty"
		errMsg.ErrorType = "CONTENT_FIELD_EMPTY"
		ErrorResponse(w, errMsg, http.StatusBadRequest)
		return false
	}
	return true
} */
