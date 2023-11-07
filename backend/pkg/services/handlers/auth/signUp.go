package auth

import (
	"errors"
	"net/http"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/helpers"
	"social-network/pkg/models"
	"social-network/pkg/services/handlers/images"
	"strconv"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	helpers.EnableCors(&w)
	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20)
		var firstname string
		var lastname string
		var email string
		var username string
		var about string
		var age int
		var gender string
		var password string
		var fileType string
		for key, value := range r.Form {

			switch key {
			case "firstname":
				firstname = value[0]
			case "lastname":
				lastname = value[0]
			case "email":
				email = value[0]
			case "username":
				username = value[0]
			case "about":
				about = value[0]
			case "age":
				age, _ = strconv.Atoi(value[0])
			case "gender":
				gender = value[0]
			case "password":
				password = value[0]
			
			case "file_type":
				fileType = value[0]
			}
		}
		
		if helpers.ValidateUserData(w, firstname, lastname, email, username, gender, password, age) {
			var u models.User
			u.FirstName = firstname
			u.LastName = lastname
			u.Email = email
			u.Username = username
			u.About = about
			u.Age = age
			u.Gender = gender
			u.Password = password

			id, err := sqlite.InsertUser(u)
			if fileType != "undefined" && id != 0 {
				images.UploadFile(r, id)
				path := `/images/avatars/` + strconv.Itoa(id) + "." + fileType
				sqlite.InsertUserAvatar(id, path)
			}
			if err != nil {
				var errMsg helpers.ErrorMsg
				if errors.Is(err, models.ErrDuplicateUsername) {
					errMsg.ErrorDescription = "Username already taken."
					errMsg.ErrorType = "USERNAME_ALREADY_TAKEN"
					helpers.ErrorResponse(w, errMsg, http.StatusUnsupportedMediaType)
					return
				}
				if errors.Is(err, models.ErrTooManySpaces) {
					errMsg.ErrorDescription = "Too many spaces."
					errMsg.ErrorType = "DOUBLE_SPACES_IF_FIELDS"
					helpers.ErrorResponse(w, errMsg, http.StatusUnsupportedMediaType)
					return
				}
				if errors.Is(err, models.ErrDuplicateEmail) {
					errMsg.ErrorDescription = "Email already taken."
					errMsg.ErrorType = "EMAIL_ALREADY_TAKEN"
					helpers.ErrorResponse(w, errMsg, http.StatusNotAcceptable)
					return
				}
				
				helpers.ErrorResponse(w, helpers.InternalErrorMsg, http.StatusInternalServerError)
				return
			} 
		}
	}
}
