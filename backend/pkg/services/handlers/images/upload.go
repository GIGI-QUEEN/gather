package images

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"social-network/pkg/db/sqlite"
	"strconv"
)

func UploadFile(r *http.Request, id int) {
	file, imageType := RetrieveFile(r)
	fileType := CheckFileType(file)
	switch imageType {
	case "avatar":
		tempFile, _ := os.CreateTemp("images/avatars", "*"+fileType)

		os.Rename(tempFile.Name(), "images/avatars/"+strconv.Itoa(id)+fileType)

		defer tempFile.Close()

		filePath := `/images/avatars/user-` + strconv.Itoa(id) + `-avatar` + fileType

		_ = sqlite.InsertUserAvatar(id, filePath)

		tempFile.Write(file)
	case "post-image":
		tempFile, _ := os.CreateTemp("images/posts-images", "*"+fileType)
		os.Rename(tempFile.Name(), "images/posts-images/"+strconv.Itoa(id)+fileType)
		defer tempFile.Close()
		tempFile.Write(file)
	case "group-post-image":
		tempFile, _ := os.CreateTemp("images/group-posts-images", "*"+fileType)
		os.Rename(tempFile.Name(), "images/group-posts-images/"+strconv.Itoa(id)+fileType)
		defer tempFile.Close()
		tempFile.Write(file)
	case "comment-image":
		tempFile, _ := os.CreateTemp("images/comments-images/", "*"+fileType)
		os.Rename(tempFile.Name(), "images/comments-images/"+strconv.Itoa(id)+fileType)
		defer tempFile.Close()
		tempFile.Write(file)
	case "group-post-comment-image":
		tempFile, _ := os.CreateTemp("images/group-posts-comments-images/", "*"+fileType)
		os.Rename(tempFile.Name(), "images/group-posts-comments-images/"+strconv.Itoa(id)+fileType)
		defer tempFile.Close()
		tempFile.Write(file)
	case "message-image":
		tempFile, _ := os.CreateTemp("images/messages-images/", "*"+fileType)
		os.Rename(tempFile.Name(), "images/messages-images/"+strconv.Itoa(id)+fileType)
		defer tempFile.Close()
		tempFile.Write(file)
	}
}

//returns file as []byte and type of image (based on form name specified on fronent) received to specify path where to store this image

func RetrieveFile(r *http.Request) ([]byte, string) {
	err := r.ParseMultipartForm(10 << 20)
	n := r.Form.Get("image_type")

	file, _, err := r.FormFile(n)
	if err != nil {
		fmt.Printf("Error retreiving the File: %s", err)
		return nil, ""
	}

	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	return fileBytes, n
}

func CheckFileType(file []byte) string {
	contetType := http.DetectContentType(file)
	var fileType string

	switch contetType {
	case "image/png":
		fileType = ".png"
	case "image/jpeg":
		fileType = ".jpg"
	case "image/gif":
		fileType = ".gif"
	}
	return fileType
}
