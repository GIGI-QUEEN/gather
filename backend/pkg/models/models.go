package models

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrDuplicateUsername  = errors.New("models: duplicate username")
	ErrTooManySpaces      = errors.New("inupt data: too many spaces in field")
	ErrInternalServer     = errors.New("INTERNAL_SERVER_ERROR")
	ErrWrongFileType      = errors.New("wrong file type")
	ErrBadRequest         = errors.New("BAD_REQUEST")
)

type Post struct {
	Id               int         `json:"id,omitempty"`
	Title            string      `json:"title,omitempty"`
	Content          string      `json:"content,omitempty"`
	CreateDate       int         `json:"creationdate,omitempty"`
	Categories       []*Category `json:"categories,omitempty"`
	User             *User       `json:"user,omitempty"`
	Likes            []*User     `json:"likes,omitempty"`
	Dislikes         []*User     `json:"dislikes,omitempty"`
	Comments         []*Comment  `json:"comments,omitempty"`
	IsLikedByUser    bool        `json:"-"`
	IsDislikedByUser bool        `json:"-"`
	Privacy          string      `json:"privacy,omitempty"`
	AllowedUsers     []*User     `json:"allowed_users,omitempty"`
	Image            string      `json:"image,omitempty"`
	ImageType        string      `json:"image_type,omitempty"`
}

type User struct {
	Id              int      `json:"id,omitempty"`
	FirstName       string   `json:"firstname,omitempty"`
	LastName        string   `json:"lastname,omitempty"`
	Age             int      `json:"age,omitempty"`
	Gender          string   `json:"gender,omitempty"`
	Username        string   `json:"username,omitempty"`
	Email           string   `json:"email,omitempty"`
	Password        string   `json:"password,omitempty"`
	HashedPassword  []byte   `json:"-"`
	Created         int      `json:"created,omitempty"`
	Status          int      `json:"status,omitempty"`
	Followers       []*User  `json:"followers,omitempty"`
	Followings      []*User  `json:"followings,omitempty"`
	Posts           []*Post  `json:"posts,omitempty"`
	About           string   `json:"about,omitempty"`
	Privacy         string   `json:"privacy,omitempty"`
	Avatar          string   `json:"avatar,omitempty"`
	FollowRequests  []*User  `json:"follow_requests,omitempty"`
	FollowRequested bool     `json:"follow_requested,omitempty"`
	MemberInGroups  []*Group `json:"member_in_groups,omitempty"`
	CreatedDate     int      `json:"created_date,omitempty"`
}

type Follower struct {
	Follower *User `json:"follower,omitempty"`
	Followed *User `json:"followed"`
}

type Session struct {
	Id   string
	User *User
}

type Comment struct {
	Id               int    `json:"id"`
	PostId           int    `json:"postId"`
	User             *User  `json:"user"`
	Content          string `json:"content"`
	IsLikedByUser    bool   `json:"isLikedByUser"`
	IsDislikedByUser bool   `json:"isDislikedByuser"`
	LikeUsers        []int  `json:"likeUsers"`
	DislikeUsers     []int  `json:"dislikeUsers"`
	LikeCount        int    `json:"likeCount"`
	DislikeCount     int    `json:"dislikeCount"`
	Image            string `json:"image,omitempty"`
	ImageType        string `json:"image_type,omitempty"`
}

type CommentReaction struct {
	CommentLikeDislike string `json:"commentlikedislike"`
}

type PostReaction struct {
	PostComment      string `json:"postcomment"`
	PostLikeDislike  string `json:"postlikedislike"`
	PostCommentImage string `json:"image"`
	ImageType        string `json:"image_type,omitempty"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// used for retriving messages from DB
type MessageWithUsersInfo struct {
	Sender      *User  `json:"sender"`
	Recipient   *User  `json:"recipient"`
	Message     string `json:"message"`
	CreatedDate int    `json:"created_date"`
}

type Image struct {
	ImageType string `json:"image_type,omitempty"`
}
