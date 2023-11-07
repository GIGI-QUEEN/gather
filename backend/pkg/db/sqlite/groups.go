package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"social-network/pkg/models"
	"strconv"
)

func InsertGroup(group models.Group, adminUserId int) (int, error) {
	result, err := DB.Exec("INSERT INTO groups(admin_user_id, title, description, created_date) values (?, ?, ?, strftime('%s','now'))", adminUserId, group.Title, group.Description)
	if err != nil {
		return -1, err
	}
	groupId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(groupId), nil
}

func InsertGroupConversationRoom(groupId, userId int) error {
	name := fmt.Sprintf("conversation_of_group_id_%d", groupId)

	result, err := DB.Exec("INSERT INTO conversation_rooms(name) values (?)", name)
	if err != nil {
		return err
	}

	conversationId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	_, err = DB.Exec("INSERT INTO conversation_rooms_users(user_id, conversation_id) values (?, ?)", userId, conversationId)
	if err != nil {
		return err
	}

	return nil
}

func InsertGroupPost(post models.GroupPost, groupId int, creatorID int) (int, error) {
	stmt := "INSERT INTO group_posts(post_title, post_content, group_id, user_id, created_date) values (?, ?, ?, ?, strftime('%s','now'))"
	result, err := DB.Exec(stmt, post.Title, post.Content, groupId, creatorID)

	if err != nil {
		return 0, err
	}
	postId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(postId), nil
}

func InsertUserAsAdminToGroup(groupId, adminUserId int) error {
	isAdmin := 1
	isApproved := 1
	_, err := DB.Exec("INSERT INTO group_users(group_id, group_user_id, is_admin, is_approved, created_date) values (?, ?, ?, ?, strftime('%s','now'))", groupId, adminUserId, isAdmin, isApproved)
	if err != nil {
		return err
	}
	return nil
}

func GetGroupsList() ([]*models.Group, error) {
	stmt := `SELECT * FROM groups ORDER BY created_date DESC`
	rows, err := DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		group := &models.Group{}
		group.Admin = &models.User{}
		err = rows.Scan(&group.Id, &group.Admin.Id, &group.Title, &group.Description, &group.CreatedDate)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func GetGroupById(groupId int, userId int) (*models.GroupViewForUser, error) {
	group := &models.GroupViewForUser{}
	group.GroupPosts = []*models.GroupPost{}
	group.UsersJoinRequests = []*models.User{}
	row := DB.QueryRow("select group_id, title, description, admin_user_id, created_date from groups where group_id = ?", groupId)
	group.Admin = &models.User{}
	err := row.Scan(&group.Id, &group.Title, &group.Description, &group.Admin.Id, &group.CreatedDate)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	group.GroupMembers, err = getGroupMembersCount(groupId)
	if err != nil {
		return nil, err
	}

	if userIsInvited(userId, groupId) {
		group.UserInvited = true
	}
	isApproved, err := CheckUserIsApproved(groupId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return group, nil
		}
	}
	if isApproved == 1 {
		group.GroupPosts, err = getGroupPosts(groupId)
		if err != nil {
			return nil, nil
		}

		group.GroupEvents, err = getGroupEvents(groupId, userId)
		if err != nil {
			return nil, nil
		}

		group.JoinRequested = true
		group.JoinApproved = true
	} else if isApproved == 0 {
		group.JoinRequested = true
	}

	if userIsAdmin(userId, group.Admin.Id) {
		if group.UsersJoinRequests, err = GetPossibleJoinRequests(groupId); err != nil {
			return nil, nil
		}
	}

	return group, nil
}

func GetGroupIdsWhereUserIsAdmin(userId int) []int {
	rows, err := DB.Query("SELECT group_id FROM group_users WHERE group_user_id = ? AND is_admin = 1", userId)
	if err != nil {
		return nil
	}
	defer rows.Close()

	// iterate over results
	var groupIds []int
	for rows.Next() {
		var groupId int
		err := rows.Scan(&groupId)
		if err != nil {
			return nil
		}
		groupIds = append(groupIds, groupId)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return groupIds
}

// Notifications related query
func GetPossibleJoinRequests(groupId int) ([]*models.User, error) {
	rows, err := DB.Query(`
		SELECT u.id, u.firstname, u.lastname, u.username, gu.created_date
		FROM users u
		INNER JOIN group_users gu ON gu.group_user_id = u.id
		WHERE gu.group_id = ? AND gu.is_admin = 0 AND gu.is_approved = 0
	`, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.CreatedDate)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetGroupInfoById(groupId int) (*models.Group, error) {
	stmt := `SELECT group_id, title FROM groups WHERE group_id = ?`
	group := &models.Group{}

	row := DB.QueryRow(stmt, groupId)
	err := row.Scan(&group.Id, &group.Title)

	if err != nil {
		return nil, err
	}

	return group, nil
}

func GetUsersWantToFollow(userId int) ([]*models.User, error) {
	const FollowStatus = 0
	rows, err := DB.Query(`
		SELECT 
			following_user_id,
			u.firstname,
			u.username
		FROM follower f
		INNER JOIN users u ON  u.id = f.following_user_id
		WHERE followed_user_id = ?
		AND follow_status = ?
	`, userId, FollowStatus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.Id, &user.FirstName, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetPossibleCreatedGroupEvents(userId int) ([]*models.GroupEventCreated, error) {
	query := `
	SELECT 
		group_events.event_id,
		group_events.title AS event_title,
		group_events.created_date,
		groups.group_id,
		groups.title AS group_title
	FROM 
		group_events_members_dependency 
	JOIN 
		group_events ON group_events_members_dependency.event_id = group_events.event_id 
	JOIN 
		groups ON group_events.group_id = groups.group_id 
	WHERE 
    	group_events_members_dependency.user_id = ? 
    AND 
		group_events_members_dependency.is_seen = ?	`

	const isSeen = 0
	unseenEvents, err := DB.Query(query, userId, isSeen)
	if err != nil {
		return nil, err
	}
	defer unseenEvents.Close()

	eventsCreated := make([]*models.GroupEventCreated, 0)
	for unseenEvents.Next() {
		eventCreated := &models.GroupEventCreated{Group: &models.Group{}, Event: &models.GroupEvent{}}
		err := unseenEvents.Scan(&eventCreated.Event.EventId, &eventCreated.Event.Title, &eventCreated.Event.CreatedDate, &eventCreated.Group.Id, &eventCreated.Group.Title)
		if err != nil {
			return nil, err
		}
		eventsCreated = append(eventsCreated, eventCreated)
	}
	if err = unseenEvents.Err(); err != nil {
		return nil, err
	}
	return eventsCreated, nil
}

// Notifications related query
func GetPossibleJoinInvites(userId int) ([]*models.GroupJoinInvites, error) {
	invitesRows, err := DB.Query(`
	SELECT ui.group_id, g.title, ui.from_user_id, ui.created_date
	FROM user_` + strconv.Itoa(userId) + `_groups_invitations ui
	INNER JOIN groups g ON g.group_id = ui.group_id
	WHERE ui.is_accepted = 0`)
	if err != nil {
		return nil, err
	}
	defer invitesRows.Close()

	invites := make([]*models.GroupJoinInvites, 0)
	for invitesRows.Next() {
		invite := &models.GroupJoinInvites{Group: &models.Group{}, UserInvited: &models.User{}}
		err := invitesRows.Scan(&invite.Group.Id, &invite.Group.Title, &invite.UserInvited.Id, &invite.CreatedDate)
		if err != nil {
			return nil, err
		}

		// Get user details for the invited user
		userRows, err := DB.Query(`SELECT id, firstname, lastname, username FROM users WHERE id = ?`, invite.UserInvited.Id)
		if err != nil {
			return nil, err
		}
		defer userRows.Close()

		if userRows.Next() {
			err := userRows.Scan(&invite.UserInvited.Id, &invite.UserInvited.FirstName, &invite.UserInvited.LastName, &invite.UserInvited.Username)
			if err != nil {
				return nil, err
			}
		}

		invites = append(invites, invite)
	}
	if err = invitesRows.Err(); err != nil {
		return nil, err
	}
	return invites, nil
}

func userIsInvited(userId, groupId int) bool {
	stmt := `SELECT is_accepted FROM user_` + strconv.Itoa(userId) + `_groups_invitations WHERE group_id = ?`
	var isAccepted int
	row := DB.QueryRow(stmt, groupId, userId)
	err := row.Scan(&isAccepted)

	if isAccepted == 0 && err == nil {
		return true
	}
	return false
}

func userIsAdmin(userId, adminId int) bool {
	return userId == adminId
}

func InsertJoinRequest(groupId, userId int) error {
	if _, err := CheckUserIsApproved(groupId, userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			isAdmin := 0
			isApproved := 0
			_, err := DB.Exec("INSERT INTO group_users(group_id, group_user_id, is_admin, is_approved, created_date) values (?, ?, ?, ?, strftime('%s','now'))", groupId, userId, isAdmin, isApproved)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func RemoveUserFromGroup(groupId, userId int) error {
	_, err := DB.Exec("DELETE FROM group_users WHERE group_id = ? AND group_user_id = ?", groupId, userId)
	if err != nil {
		return err
	}

	groupName := fmt.Sprintf("conversation_of_group_id_%d", groupId)
	_, err = DB.Exec("DELETE from conversation_rooms_users where user_id = ? AND conversation_id = (SELECT id FROM conversation_rooms WHERE name = ?)", userId, groupName)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserIsAdmin(groupId, userId int) bool {
	stmt := `SELECT is_admin FROM group_users WHERE group_id = ? and group_user_id = ?`
	var isAdmin int
	row := DB.QueryRow(stmt, groupId, userId)
	err := row.Scan(&isAdmin)

	if isAdmin == 1 && err == nil {
		return true
	}
	return false
}

func CheckUserIsMember(groupId, userId int) bool {
	stmt := `SELECT is_approved FROM group_users WHERE group_id = ? and group_user_id = ?`
	var isApproved int
	row := DB.QueryRow(stmt, groupId, userId)
	err := row.Scan(&isApproved)

	if isApproved == 1 && err == nil {
		return true
	}
	return false
}

func InsertNewUserToGroup(groupId, userId int) error {
	_, err := DB.Exec("UPDATE group_users SET is_approved = 1 where group_id = ? AND group_user_id = ?", groupId, userId)
	if err != nil {
		return err
	}

	_, err = DB.Exec("INSERT INTO conversation_rooms_users(user_id, conversation_id) values (?, (SELECT id FROM conversation_rooms WHERE name = ?))",
		userId, fmt.Sprintf("conversation_of_group_id_%d", groupId))
	if err != nil {
		return err
	}

	return nil
}

func RejectUserJoinRequest(groupId, userId int) error {
	_, err := DB.Exec("DELETE from group_users where group_id = ? AND group_user_id = ?", groupId, userId)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserIsApproved(groupId, userId int) (int, error) {
	stmt := `SELECT is_approved FROM group_users WHERE group_id = ? and group_user_id = ?`
	var isApproved int
	row := DB.QueryRow(stmt, groupId, userId)
	err := row.Scan(&isApproved)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, sql.ErrNoRows
		}
		return -1, err
	}
	return isApproved, err
}

func GetGroupPostById(postId int) (*models.GroupPost, error) {
	row := DB.QueryRow("select post_id, post_title, post_content, user_id, created_date, img_path from group_posts where post_id = ?", postId)
	post := &models.GroupPost{}
	post.User = &models.User{}
	err := row.Scan(&post.PostId, &post.Title, &post.Content, &post.User.Id, &post.CreatedDate, &post.Image)

	if err != nil {
		return nil, err
	}
	post.PostComments, err = GetGroupPostCommentsByPostId(post.PostId)
	if err != nil {
		return nil, err
	}

	post.User, err = GetUserForPostInfo(post.User.Id)
	if err != nil {
		return nil, err
	}

	likes, err := GetGroupPostLikeUsers(post.PostId)
	if err != nil {
		return nil, err
	}
	dislikes, err := GetGroupPostDislikeUsers(post.PostId)
	if err != nil {
		return nil, err
	}
	for _, userId := range likes {
		post.Likes = append(post.Likes, &models.User{Id: userId})
	}
	for _, userId := range dislikes {
		post.Dislikes = append(post.Dislikes, &models.User{Id: userId})
	}
	return post, nil
}

func getGroupPosts(groupId int) ([]*models.GroupPost, error) {

	rows, err := DB.Query("select post_id, post_title, post_content, user_id, created_date, img_path from group_posts where group_id = ? ORDER BY created_date DESC", groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groupPosts []*models.GroupPost

	for rows.Next() {

		post := &models.GroupPost{}
		post.User = &models.User{}
		err := rows.Scan(&post.PostId, &post.Title, &post.Content, &post.User.Id, &post.CreatedDate, &post.Image)

		if err != nil {
			return nil, err
		}
		post.User, err = GetUserForPostInfo(post.User.Id)
		if err != nil {
			return nil, err
		}

		likes, err := GetGroupPostLikeUsers(post.PostId)
		if err != nil {
			return nil, err
		}
		dislikes, err := GetGroupPostDislikeUsers(post.PostId)
		if err != nil {
			return nil, err
		}
		for _, userId := range likes {
			post.Likes = append(post.Likes, &models.User{Id: userId})
		}
		for _, userId := range dislikes {
			post.Dislikes = append(post.Dislikes, &models.User{Id: userId})
		}
		groupPosts = append(groupPosts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groupPosts, nil
}

func getGroupMembersCount(groupId int) (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM group_users WHERE is_approved = 1 AND group_id = ?", groupId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func InsertGroupPostImagePath(img_path string, post_id int) error {
	_, err := DB.Exec("UPDATE group_posts SET img_path = ? where post_id = ?", img_path, post_id)
	if err != nil {
		return err
	}
	return nil
}

// Invite User To Group
func InviteNewUserToGroup(groupId, userIdToInvite, userId int) error {
	if invitationsTableExists(userIdToInvite) {
		err := insertInvite(groupId, userIdToInvite, userId)
		if err != nil {
			return err
		}
	} else {
		err := createInviteTable(userIdToInvite)
		if err != nil {
			return err
		}
		err = insertInvite(groupId, userIdToInvite, userId)
		if err != nil {
			return err
		}
	}
	return nil
}

func invitationsTableExists(userIdToInvite int) bool {
	_, table_check := DB.Query("SELECT * FROM " + "user_" + strconv.Itoa(userIdToInvite) + "_groups_invitations;")
	if table_check == nil { // if no error, table exists
		return true
	} else {
		return false
	}
}

func insertInvite(groupId, userIdToInvite, userId int) error {
	isAccepted := 0
	var stmt = `INSERT INTO user_` + strconv.Itoa(userIdToInvite) + `_groups_invitations(group_id, from_user_id, is_accepted, created_date) values (?, ?, ?, strftime('%s','now'))`
	_, err := DB.Exec(stmt, groupId, userId, isAccepted)
	if err != nil {
		return err
	}
	return nil
}

func createInviteTable(userIdToInvite int) error {
	var create = `create table user_` + strconv.Itoa(userIdToInvite) + `_groups_invitations(
		id INTEGER not null primary key autoincrement,
		group_id INTEGER not null,
		from_user_id INTEGER not null,
		is_accepted INTEGER not null,
		created_date INTEGER not null)
		`
	_, err := DB.Exec(create)
	if err != nil {
		return err
	}
	if err = DB.Ping(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// Accept Group Invite
func AcceptInviteToGroup(groupId, userId int) error {
	// check if user realy been invited
	if err := checkUserForBeingInvited(groupId, userId); err != nil {
		return err
	}

	// insert him into group users
	isAdmin := 0
	isApproved := 1
	_, err := DB.Exec("INSERT INTO group_users(group_id, group_user_id, is_admin, is_approved, created_date) values (?, ?, ?, ?, strftime('%s','now'))", groupId, userId, isAdmin, isApproved)
	if err != nil {
		return err
	}

	_, err = DB.Exec("INSERT INTO conversation_rooms_users(user_id, conversation_id) values (?, (SELECT id FROM conversation_rooms WHERE name = ?))",
		userId, fmt.Sprintf("conversation_of_group_id_%d", groupId))
	if err != nil {
		return err
	}

	// after that delete row
	if err = deleteGroupInvite(groupId, userId); err != nil {
		return err
	}
	return nil
}

// Reject Group Invite
func RejectInviteToGroup(groupId, userId int) error {
	// check if user realy been invited
	if err := checkUserForBeingInvited(groupId, userId); err != nil {
		return err
	}
	// after that delete row
	if err := deleteGroupInvite(groupId, userId); err != nil {
		return err
	}
	return nil
}

func checkUserForBeingInvited(groupId, userId int) error {
	stmt := `SELECT is_accepted FROM ` + `user_` + strconv.Itoa(userId) + `_groups_invitations WHERE group_id = ?;`
	var isAccepted int
	row := DB.QueryRow(stmt, groupId)
	err := row.Scan(&isAccepted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}
	return nil
}

func deleteGroupInvite(groupId, userId int) error {
	stmt := `DELETE FROM ` + `user_` + strconv.Itoa(userId) + `_groups_invitations WHERE group_id = ?`
	_, err := DB.Exec(stmt, groupId, userId)
	if err != nil {
		return err
	}
	return nil
}
