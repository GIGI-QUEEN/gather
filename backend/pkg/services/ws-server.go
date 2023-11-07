package services

import (
	"encoding/json"
	"log"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

// NewWebsocketServer creates a new WsServer type
func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

// Run our websocket server, accepting various requests
func (server *WsServer) Run() {
	for {
		select {
		case client := <-server.register:
			server.registerClient(client)
		case client := <-server.unregister:
			server.unregisterClient(client)
		case message := <-server.broadcast:
			server.broadcastToClients(message)
		}
	}
}

func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	delete(server.clients, client)
}

func (server *WsServer) broadcastToClients(message []byte) {
	var rawMessage map[string]interface{}
	err := json.Unmarshal(message, &rawMessage)
	if err != nil {
		return
	}

	switch rawMessage["event_type"].(string) {
	case "ws_msg_event":
		handleMessageEvent(message, rawMessage, server.clients)
	case "ws_group_msg_event":
		handleGroupMessageEvent(message, rawMessage, server.clients)
	case "ws_post_comment_event":
		handleGroupPostCommentedEvent(message, rawMessage, server.clients)
	case "ws_group_event_created_event":
		handleGroupEventCreated(message, rawMessage, server.clients)
	case "ws_group_join_request_event":
		handleGroupJoinRequest(message, rawMessage, server.clients)
	case "ws_group_join_invite_event":
		handleGroupJoinInvite(message, rawMessage, server.clients)
	case "ws_follow_request_event":
		handleFollowRequest(message, rawMessage, server.clients)
	}

}

func handleMessageEvent(message []byte, rawMessage map[string]interface{}, serverClients map[*Client]bool) {
	newMessage := &models.WsMessage{}
	newMessage.Sender = int(rawMessage["sender"].(float64))
	newMessage.Recipient = int(rawMessage["recipient"].(float64))
	newMessage.Message = rawMessage["message"].(string)
	if err := sqlite.SaveMessage(*newMessage); err != nil {
		log.Println("ws_handleMessageEvent(): ", err)
	}
	for client := range serverClients {
		if client.id == newMessage.Recipient {
			client.send <- message
		}
	}
}

func handleGroupMessageEvent(message []byte, rawMessage map[string]interface{}, serverClients map[*Client]bool) {
	newMessage := &models.WsMessage{}
	newMessage.Sender = int(rawMessage["sender"].(float64)) // same will be recipient
	groupId := int(rawMessage["group_id"].(float64))
	newMessage.Message = rawMessage["message"].(string)

	members, _ := sqlite.SaveGroupMessageAndReturnGroupMemberIds(*newMessage, groupId)
	for client := range serverClients {
		for _, groupMemberId := range members {
			if client.id == groupMemberId {
				client.send <- message
			}
		}
	}
}

func handleGroupPostCommentedEvent(message []byte, rawMessage map[string]interface{}, serverClients map[*Client]bool) {
	clientId := int(rawMessage["post_author"].(float64))

	for client := range serverClients {
		if client.id == clientId {
			client.send <- message
		}
	}
}

func handleGroupEventCreated(message []byte, rawMessage map[string]interface{}, serverClients map[*Client]bool) {
	groupId := int(rawMessage["group_id"].(float64))
	eventCreatorId := int(rawMessage["event_creator"].(float64))

	userIds, err := sqlite.GetGroupMembersIds(groupId)
	if err != nil {
		return
	}

	for client := range serverClients {
		for _, userId := range userIds {
			if client.id == userId {
				if client.id != eventCreatorId {
					client.send <- message
				}
			}
		}
	}
}

func handleGroupJoinRequest(message []byte, rawMessage map[string]interface{}, serverClients map[*Client]bool) {
	groupId := int(rawMessage["group_id"].(float64))

	adminId, err := sqlite.GetGroupAdminId(groupId)
	if err != nil {
		return
	}
	for client := range serverClients {
		if client.id == adminId {
			client.send <- message
		}
	}
}

func handleGroupJoinInvite(message []byte, rawMessage map[string]interface{}, serverClients map[*Client]bool) {
	userIdToInvite := int(rawMessage["user_to_invite"].(float64))

	for client := range serverClients {
		if client.id == userIdToInvite {
			client.send <- message
		}
	}
}

func handleFollowRequest(message []byte, rawMessage map[string]interface{}, serverClients map[*Client]bool) {
	followed_user_id := int(rawMessage["user_to_follow"].(float64))

	for client := range serverClients {
		if client.id == followed_user_id {
			client.send <- message
		}
	}
}
