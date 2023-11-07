package sqlite

import (
	"social-network/pkg/models"
)

// function for GetGroupById
func getGroupEvents(groupId, userId int) ([]*models.GroupEvent, error) {

	_, err := DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	stmt := `
	SELECT
		group_events.event_id,
		group_events.title,
		group_events.description,
		group_events.event_date,
		group_events.created_date,
		COALESCE(group_events_members_dependency.going_decision, 0) AS going_decision,
		(SELECT COUNT(*) FROM group_events_members_dependency WHERE event_id = group_events.event_id AND going_decision = 1) AS members_going
	FROM
		group_events
	LEFT JOIN 
		group_events_members_dependency ON group_events.event_id = group_events_members_dependency.event_id AND group_events_members_dependency.user_id = ?
	WHERE
		group_events.group_id = ?
	ORDER BY
		event_date ASC;
	`

	rows, err := DB.Query(stmt, userId, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupEvents []*models.GroupEvent

	for rows.Next() {
		event := &models.GroupEvent{}
		err := rows.Scan(&event.EventId, &event.Title, &event.Description, &event.EventDate, &event.CreatedDate, &event.GoingDecision, &event.MembersGoing)
		if err != nil {
			return nil, err
		}

		groupEvents = append(groupEvents, event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groupEvents, nil
}

func InsertGroupEvent(groupId int, groupEvent models.GroupEvent, userId int) error {
	res, err := DB.Exec("INSERT INTO group_events(group_id, title, description, event_date, created_date) values (?,?,?,?, strftime('%s','now'))",
		groupId, groupEvent.Title, groupEvent.Description, groupEvent.EventDate)
	if err != nil {
		return err
	}

	eventId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	// insert straight away event creator with is_seen column as 1
	isSeen := 1
	_, err = DB.Exec("INSERT INTO group_events_members_dependency(event_id, user_id, is_seen) values (?,?,?)", eventId, userId, isSeen)
	if err != nil {
		return err
	}
	// insert other group memebers with field is_seen = 0, for the notification purposes
	if err := insertAllGroupMembersToEvent(groupId, int(eventId), userId); err != nil {
		return err
	}

	return nil
}

func insertAllGroupMembersToEvent(groupId, eventId, eventCreator int) error {
	const isSeen = 0        // 0 false, 1 true
	const goingDesicion = 0 // 0 not decided, 1 going, 2 not going
	query := `INSERT INTO group_events_members_dependency (event_id, user_id, going_decision, is_seen)
			SELECT ?, group_users.group_user_id, ?, ?
			FROM group_users
			WHERE group_users.group_id = ? AND group_users.group_user_id != ?;`

	_, err := DB.Exec(query, eventId, goingDesicion, isSeen, groupId, eventCreator)
	if err != nil {
		return err
	}

	return nil
}

func AcceptInviteToGroupEvent(eventId, userId int) error {
	// check for being present in group_events_members_dependency
	// if not, create new row with accept decision and is_seen = 1
	// if yes, just change decision

	const IsSeen = 1        // 0 false, 1 true
	const GoingDecision = 1 // 0 not decided, 1 going, 2 not going

	var exists int
	query := "SELECT EXISTS (SELECT 1 FROM group_events_members_dependency WHERE event_id = ? AND user_id = ?)"
	err := DB.QueryRow(query, eventId, userId).Scan(&exists)
	if err != nil {
		return err
	}

	if exists == 0 {
		_, err = DB.Exec("INSERT INTO group_events_members_dependency (event_id, user_id, going_decision, is_seen) VALUES (?, ?, ?, ?)", eventId, userId, GoingDecision, IsSeen)
		if err != nil {
			return err
		}
	} else {
		_, err = DB.Exec("UPDATE group_events_members_dependency SET going_decision = ?, is_seen = ? WHERE event_id = ? AND user_id = ?", GoingDecision, IsSeen, eventId, userId)
		if err != nil {
			return err
		}
	}

	return nil
}

func RejectInviteToGroupEvent(eventId, userId int) error {
	// check for being present in group_events_members_dependency
	// if not, create new row with accept decision and is_seen = 1
	// if yes, just change decision

	const IsSeen = 1        // 0 false, 1 true
	const GoingDecision = 2 // 0 not decided, 1 going, 2 not going

	var exists int
	query := "SELECT EXISTS (SELECT 1 FROM group_events_members_dependency WHERE event_id = ? AND user_id = ?)"
	err := DB.QueryRow(query, eventId, userId).Scan(&exists)
	if err != nil {
		return err
	}

	if exists == 0 {
		_, err = DB.Exec("INSERT INTO group_events_members_dependency (event_id, user_id, going_decision, is_seen) VALUES (?, ?, ?, ?)", eventId, userId, GoingDecision, IsSeen)
		if err != nil {
			return err
		}
	} else {
		_, err = DB.Exec("UPDATE group_events_members_dependency SET going_decision = ?, is_seen = ? WHERE event_id = ? AND user_id = ?", GoingDecision, IsSeen, eventId, userId)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteNotificationAboutCreatedGroupEvent(eventId, userId int) error {
	const IsSeen = 1
	_, err := DB.Exec("UPDATE group_events_members_dependency SET is_seen = ? WHERE event_id = ? AND user_id = ?", IsSeen, eventId, userId)
	if err != nil {
		return err
	}

	return nil
}
