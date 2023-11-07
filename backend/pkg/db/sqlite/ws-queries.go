package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"social-network/pkg/models"
)

func SaveMessage(message models.WsMessage) error {

	recipientId := message.Recipient
	senderId := message.Sender
	messageFromStruct := message.Message

	name := fmt.Sprintf("conversation_of_user_ids_%d_%d", senderId, recipientId)
	name2 := fmt.Sprintf("conversation_of_user_ids_%d_%d", recipientId, senderId)
	// Check if there is an existing conversation room between the sender and recipient.
	var conversationID int64
	err := DB.QueryRow("SELECT id FROM conversation_rooms WHERE name = ? OR name = ?", name, name2).Scan(&conversationID)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		// If there is no existing conversation room, create a new one.
		result, err := DB.Exec("INSERT INTO conversation_rooms (name) VALUES (?)", fmt.Sprintf("conversation_of_user_ids_%d_%d", senderId, recipientId))
		if err != nil {
			return err
		}
		conversationID, err = result.LastInsertId()
		if err != nil {
			return err
		}

		// Associate both users with the new conversation room.
		_, err = DB.Exec("INSERT INTO conversation_rooms_users (user_id, conversation_id) VALUES (?, ?)", senderId, conversationID)
		if err != nil {
			return err
		}
		_, err = DB.Exec("INSERT INTO conversation_rooms_users (user_id, conversation_id) VALUES (?, ?)", recipientId, conversationID)
		if err != nil {
			return err
		}
	}

	// Insert the first message into the conversation room.
	_, err = DB.Exec("INSERT INTO conversation_messages (sender_id, recipient_id, conversation_id, message, created_date) VALUES (?, ?, ?, ?, strftime('%s','now'))",
		senderId, recipientId, conversationID, messageFromStruct)
	if err != nil {
		return err
	}
	return nil
}

func SaveGroupMessageAndReturnGroupMemberIds(message models.WsMessage, groupId int) ([]int, error) {

	_, err := DB.Exec("INSERT INTO conversation_messages (sender_id, recipient_id, conversation_id, message, created_date) VALUES (?, ?, (SELECT id FROM conversation_rooms WHERE name = ?), ?, strftime('%s','now'))",
		message.Sender, message.Sender, fmt.Sprintf("conversation_of_group_id_%d", groupId), message.Message)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query("SELECT group_user_id from group_users WHERE group_id = ? AND group_user_id != ?", groupId, message.Sender)
	if err != nil {
	}
	defer rows.Close()
	var groupMemberIds []int

	for rows.Next() {
		var memberId int
		err := rows.Scan(&memberId)
		if err != nil {
			return nil, err
		}
		groupMemberIds = append(groupMemberIds, memberId)
	}
	return groupMemberIds, nil

}

func GetMessages(senderID int, recipientID int, offset int) ([]models.WsMessage, error) {
	var messages []models.WsMessage

	name := fmt.Sprintf("conversation_of_user_ids_%d_%d", senderID, recipientID)
	name2 := fmt.Sprintf("conversation_of_user_ids_%d_%d", recipientID, senderID)

	var conversationId int64
	err := DB.QueryRow("SELECT id FROM conversation_rooms WHERE name = ? OR name = ?", name, name2).Scan(&conversationId)
	if err != nil {
		if err == sql.ErrNoRows {
			// No conversation room between the two users
			return messages, nil
		}
		return nil, err
	}

	rows, err := DB.Query("SELECT sender_id, recipient_id, conversation_id, message, created_date FROM conversation_messages WHERE conversation_id = ?"+
		" ORDER BY created_date DESC LIMIT 20 OFFSET ?", conversationId, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var senderID int
		var recipientID int
		var conversationID int
		var message string
		var createdDate int
		if err := rows.Scan(&senderID, &recipientID, &conversationID, &message, &createdDate); err != nil {
			return nil, err
		}
		messages = append(messages, models.WsMessage{Sender: senderID, Recipient: recipientID, Message: message, CreatedDate: createdDate})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func GetGroupMessages(groupId, offset int) ([]models.WsMessage, error) {
	var messages []models.WsMessage

	groupName := fmt.Sprintf("conversation_of_group_id_%d", groupId)

	var conversationId int64
	err := DB.QueryRow("SELECT id FROM conversation_rooms WHERE name = ?", groupName).Scan(&conversationId)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`
	SELECT 
		cm.sender_id, 
		cm.conversation_id, 
		cm.message, 
		cm.created_date,
		u.username,
		u.firstname
	FROM conversation_messages cm
	JOIN users u ON cm.sender_id = u.id
	WHERE cm.conversation_id = ?
	ORDER BY cm.created_date DESC
	LIMIT 20 OFFSET ?
`, conversationId, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var senderID int
		var conversationID int
		var senderUsername string
		var senderFirstname string
		var message string
		var createdDate int
		if err := rows.Scan(&senderID, &conversationID, &message, &createdDate, &senderUsername, &senderFirstname); err != nil {
			return nil, err
		}
		messages = append(messages, models.WsMessage{Sender: senderID, Message: message, CreatedDate: createdDate, SenderUsername: senderUsername, SenderFristname: senderFirstname})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil

}

func GetStartedChatsLatestMessages(userId int) ([]*models.MessageWithUsersInfo, error) {
	rows, err := DB.Query(`
	SELECT
	sender_id, 
	sender_firstname, 
	sender_lastname, 
	sender_username, 
	sender_status, 
	sender_avatar_path,
	recipient_id, 
	recipient_firstname,
	recipient_lastname, 
	recipient_username, 
	recipient_status,
	recipient_avatar_path,
	message, 
	created_date
	FROM

	(SELECT 
	u1.id AS sender_id, 
	u1.firstname AS sender_firstname, 
	u1.lastname AS sender_lastname, 
	u1.username AS sender_username, 
	u1.status as sender_status, 
	u1.avatar_path as sender_avatar_path,
	u2.id AS recipient_id, 
	u2.firstname AS recipient_firstname,
	u2.lastname AS recipient_lastname, 
	u2.username AS recipient_username, 
	u2.status as recipient_status,
	u2.avatar_path as recipient_avatar_path,
	m.message AS message, 
	m.created_date AS created_date
	FROM conversation_messages m
	INNER JOIN (
		SELECT MAX(created_date) AS max_date, conversation_id 
		FROM conversation_messages 
		WHERE sender_id = ? OR recipient_id = ? 
		GROUP BY conversation_id
	) max_dates ON m.created_date = max_dates.max_date AND m.conversation_id = max_dates.conversation_id
	INNER JOIN conversation_rooms_users cu1 ON m.sender_id = cu1.user_id AND cu1.conversation_id = m.conversation_id
	INNER JOIN conversation_rooms_users cu2 ON m.recipient_id = cu2.user_id AND cu2.conversation_id = m.conversation_id
	INNER JOIN users u1 ON m.sender_id = u1.id
	INNER JOIN users u2 ON m.recipient_id = u2.id
	WHERE cu1.user_id = ? OR cu2.user_id = ?
	ORDER BY m.created_date DESC)

	WHERE sender_id != recipient_id
	`, userId, userId, userId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*models.MessageWithUsersInfo{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	messages := []*models.MessageWithUsersInfo{}
	for rows.Next() {
		var sender models.User
		var recipient models.User
		var message models.MessageWithUsersInfo

		err := rows.Scan(
			&sender.Id,
			&sender.FirstName,
			&sender.LastName,
			&sender.Username,
			&sender.Status,
			&sender.Avatar,
			&recipient.Id,
			&recipient.FirstName,
			&recipient.LastName,
			&recipient.Username,
			&recipient.Status,
			&recipient.Avatar,
			&message.Message,
			&message.CreatedDate,
		)
		if err != nil {
			return nil, err
		}

		message.Sender = &sender
		message.Recipient = &recipient
		messages = append(messages, &message)
	}

	return messages, nil
}

func GetGroupMembersIds(groupId int) ([]int, error) {
	// execute the SQL query to fetch group_user_id values
	rows, err := DB.Query("SELECT group_user_id FROM group_users WHERE group_id = ?", groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// loop through the rows and collect the group_user_id values
	var userIds []int
	for rows.Next() {
		var userId int
		if err := rows.Scan(&userId); err != nil {
			return nil, err
		}
		userIds = append(userIds, userId)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIds, nil
}

func GetGroupAdminId(groupId int) (int, error) {

	const IsAdmin = 1
	var adminId int

	err := DB.QueryRow("SELECT group_user_id FROM group_users WHERE group_id = ? AND is_admin = ?",
		groupId, IsAdmin).Scan(&adminId)
	if err != nil {
		return -1, err
	}

	return adminId, nil
}
