package models

type Group struct {
	Id          int    `json:"id,omitempty"`
	Admin       *User  `json:"admin,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	CreatedDate int    `json:"created_date,omitempty"`
}

type GroupViewForUser struct {
	Id           int    `json:"id,omitempty"`
	Admin        *User  `json:"admin,omitempty"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	CreatedDate  int    `json:"created_date"`
	GroupMembers int    `json:"group_members,omitempty"`

	UserInvited   bool `json:"user_invited"`
	JoinRequested bool `json:"join_requested"`
	JoinApproved  bool `json:"join_approved"`

	// will be available only for admin of the group
	UsersJoinRequests []*User `json:"join_requests,omitempty"`

	// will be available only for group members approved by admin
	GroupPosts  []*GroupPost  `json:"group_posts,omitempty"`
	GroupEvents []*GroupEvent `json:"group_events,omitempty"`
}

type GroupPost struct {
	PostId           int                 `json:"post_id,omitempty"`
	Title            string              `json:"title,omitempty"`
	Content          string              `json:"content,omitempty"`
	GroupId          int                 `json:"group_id,omitempty"`
	User             *User               `json:"user,omitempty"`
	CreatedDate      int                 `json:"created_date,omitempty"`
	PostComments     []*GroupPostComment `json:"comments,omitempty"`
	Likes            []*User             `json:"likes,omitempty"`
	Dislikes         []*User             `json:"dislikes,omitempty"`
	IsLikedByUser    bool                `json:"-"`
	IsDislikedByUser bool                `json:"-"`
	Image            string              `json:"image,omitempty"`
}

type GroupEvent struct {
	EventId       int    `json:"event_id,omitempty"`
	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
	EventDate     int    `json:"event_date,omitempty"`
	MembersGoing  int    `json:"members_going"`
	GoingDecision int    `json:"going_decision"`
	CreatedDate   string `json:"created_date,omitempty"`
}

type GroupPostComment struct {
	Id               int    `json:"id,omitempty"`
	Content          string `json:"content,omitempty"`
	PostId           int    `json:"postId,omitempty"`
	User             *User  `json:"user,omitempty"`
	IsLikedByUser    bool   `json:"isLikedByUser"`
	IsDislikedByUser bool   `json:"isDislikedByuser"`
	LikeUsers        []int  `json:"likeUsers"`
	DislikeUsers     []int  `json:"dislikeUsers"`
	LikeCount        int    `json:"likeCount"`
	DislikeCount     int    `json:"dislikeCount"`
	Image            string `json:"image,omitempty"`
	ImageType        string `json:"image_type,omitempty"`
	CreatedDate      int    `json:"created_date,omitempty"`
}

type GroupPostCommentIds struct {
	CommentId []int `json:"comment_id,omitempty"`
}
