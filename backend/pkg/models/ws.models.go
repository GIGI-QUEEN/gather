package models

type WsMessage struct {
	Sender          int    `json:"sender"`
	SenderUsername  string `json:"sender_username"`
	SenderFristname string `json:"sender_firstname"`

	Recipient   int    `json:"recipient"`
	Message     string `json:"message"`
	CreatedDate int    `json:"created_date"`
}

type WsNotification struct {
	GroupId     int    `json:"group_id"`
	PostId      int    `json:"post_id"`
	Comment     string `json:"comment"`
	CreatedDate int    `json:"created_date"`
}

// event for messages plus notifications

const (
	WebSocketMessageEvent          string = "ws_msg_event"
	WebSocketGroupPostCommentEvent string = "group_post_comment_event"
)
