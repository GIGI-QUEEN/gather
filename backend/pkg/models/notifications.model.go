package models

type Notifications struct {
	UsersWantToFollow  []*User                         `json:"users_want_to_follow"`
	GroupJoinInvites   []*GroupJoinInvites             `json:"group_join_invites"`
	GroupJoinRequests  []*GroupJoinRequest             `json:"group_join_requests"`
	GroupEventsCreated []*GroupEventCreated            `json:"group_events_created"`
	CommentOnUsersPost []*GroupPostCommentNotification `json:"group_post_comment_notification"`
}

type GroupEventCreated struct {
	Group *Group      `json:"group,omitempty"`
	Event *GroupEvent `json:"group_event,omitempty"`
}

type GroupPostCommentNotification struct {
	GroupId        int    `json:"group_id,omitempty"`
	PostId         int    `json:"post_id,omitempty"`
	Title          string `json:"title,omitempty"`
	CommentId      int    `json:"comment_id,omitempty"`
	CommentContent string `json:"comment_content,omitempty"`
	CreatedDate    int    `json:"created_date,omitempty"`
}

type GroupJoinRequest struct {
	Group              *Group  `json:"group,omitempty"`
	UsersRequestedJoin []*User `json:"users_requested_join,omitempty"`
}

type GroupJoinInvites struct {
	Group       *Group `json:"group,omitempty"`
	UserInvited *User  `json:"user_invited,omitempty"`
	CreatedDate int    `json:"created_date,omitempty"`
}
