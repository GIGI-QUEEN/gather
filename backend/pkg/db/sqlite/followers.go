package sqlite

import (
	"database/sql"
	"social-network/pkg/models"
)

func GetUserFollowers(user_id int) ([]*models.User, error) {
	stmt := `SELECT id, username, firstname,lastname,avatar_path from users 
	INNER JOIN follower on follower.followed_user_id = ? 
	WHERE id = follower.following_user_id
	AND follower.follow_status = 1`

	rows, _ := DB.Query(stmt, user_id)
	defer rows.Close()
	var users []*models.User

	for rows.Next() {
		u := &models.User{}
		err := rows.Scan(&u.Id, &u.Username, &u.FirstName,&u.LastName,&u.Avatar)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}



func GetUserFollowings(user_id int) ([]*models.User, error) {
	rows, _ := DB.Query("SELECT id, username, avatar_path from users INNER JOIN follower on follower.following_user_id = ? WHERE id = follower.followed_user_id AND follower.follow_status = 1", user_id)
	defer rows.Close()
	var users []*models.User

	for rows.Next() {
		u := &models.User{}
		err := rows.Scan(&u.Id, &u.Username, &u.Avatar)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUserGroupsWhereIsMember(userId int) ([]*models.Group, error) {
	const IsApproved = 1
	rows, _ := DB.Query("SELECT g.group_id, gr.title FROM group_users g JOIN groups gr ON g.group_id = gr.group_id WHERE g.group_user_id = ? AND g.is_approved = ?", userId, IsApproved)
	defer rows.Close()
	var groups []*models.Group

	for rows.Next() {
		group := &models.Group{}
		err := rows.Scan(&group.Id, &group.Title)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

// Query for user private post to allow new followers see private posts/disallow ufollowed users see private posts
func GetUserPrivatePosts(user_id int) ([]*models.Post, error) {
	rows, err := DB.Query("select id, title, contents, create_date, user_id, privacy from posts where (user_id, privacy)= (?, \"private\") order by id desc", user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		post.Categories = []*models.Category{}
		post.User = &models.User{}
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.CreateDate, &post.User.Id, &post.Privacy)
		if err != nil {
			return nil, err
		}
		post.User, err = GetUserForPostInfo(post.User.Id)
		if err != nil {
			return nil, err
		}
		post.Categories, err = GetCategoriesByPost(post.Id)
		if err != nil {
			return nil, err
		}
		likes, err := GetPostLikeUsers(post.Id)
		if err != nil {
			return nil, err
		}
		dislikes, err := GetPostDislikeUsers(post.Id)
		if err != nil {
			return nil, err
		}
		for _, userId := range likes {
			post.Likes = append(post.Likes, &models.User{Id: userId})
		}
		for _, userId := range dislikes {
			post.Dislikes = append(post.Dislikes, &models.User{Id: userId})
		}
		posts = append(posts, post)
	}
	if err != nil {
		return nil, err
	}
	if len(posts) > 0 {
		return posts, nil

	} else {
		return nil, models.ErrNoRecord
	}
}

func FollowUser(userToFollow, userId int) error {
	row := DB.QueryRow("select following_user_id from follower where (followed_user_id, following_user_id) = (?,?)", userToFollow, userId)
	id := -1
	privatePosts, _ := GetUserPrivatePosts(userId)

	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		_, err := DB.Exec("insert into follower (followed_user_id, following_user_id) values (?,?)", userToFollow, userId)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		for _, p := range privatePosts {
			_, err := DB.Exec("insert into allowed_post_users (user_id, post_id) values  (?,?)", userId, p.Id)
			if err != nil && err != sql.ErrNoRows {
				return err
			}
		}
	case nil:
		_, err := DB.Exec("DELETE FROM follower WHERE  (following_user_id, followed_user_id) = (?,?)", userId, userToFollow)
		if err != nil {
			return err
		}
		_, err = DB.Exec("delete from allowed_post_users where user_id = ?", userId)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
	}
	return nil
}

func CheckUserPrivacy(userId int) (string, error) {
	row := DB.QueryRow("select privacy from users where id = ?", userId) 
	privacy := ""
	err := row.Scan(&privacy)
	if err != nil {
		return "", err
	}
	return privacy, nil
}

//follow status could be 0 and 1 (0 - follow request pending, 1 - follow request accepted)
//followerId is s.User.Id (current user)

func FollowUser_v2(userToFollow, followerId int) error {
	row := DB.QueryRow("select following_user_id from follower where (followed_user_id, following_user_id) = (?,?)", userToFollow, followerId)
	id := -1
	userToFollowPrivacy, _ := CheckUserPrivacy(userToFollow)
	privatePosts, _ := GetUserPrivatePosts(userToFollow)

	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		if userToFollowPrivacy == "public"{
			_, err := DB.Exec("insert into follower (followed_user_id, following_user_id, follow_status) values (?,?,?)", userToFollow, followerId, 1)
			if err != nil && err != sql.ErrNoRows {
			return err
		}
		} else {
			_, err := DB.Exec("insert into follower (followed_user_id, following_user_id, follow_status) values (?,?,?)", userToFollow, followerId, 0)
			if err != nil && err != sql.ErrNoRows {
			return err
		}
		for _, p := range privatePosts {
			_, err := DB.Exec("insert into allowed_post_users (user_id, post_id) values  (?,?)", userToFollow, p.Id)
			if err != nil && err != sql.ErrNoRows {
				return err
			}
		}
		}
	}
	return nil
}

func UnfollowUser(userToUnfollow, followerId int) error {
	row := DB.QueryRow("select following_user_id from follower where (followed_user_id, following_user_id) = (?,?)", userToUnfollow, followerId)
	id := -1
	switch err := row.Scan(&id); err {
	case nil:
		_, err := DB.Exec("DELETE FROM follower WHERE  (following_user_id, followed_user_id) = (?,?)", followerId, userToUnfollow)
		if err != nil {
			return err
		}
		_, err = DB.Exec("delete from allowed_post_users where user_id = ?", userToUnfollow)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
	}
	return nil
}


func AcceptFollowRequest(userToAccept, userId int) error {
	_, err := DB.Exec("UPDATE follower SET follow_status = 1 where (following_user_id, followed_user_id) = (?,?)", userToAccept, userId)
	if err != nil {
		return err
	}
	return err
}

func RejectFollowRequst(userToReject, userId int) error {
	_, err := DB.Exec("DELETE FROM follower WHERE  (following_user_id, followed_user_id) = (?,?)", userToReject, userId)
		if err != nil {
			return err
		}
		return nil
}

func GetUserFollowRequests(userId int) ([]*models.User, error) {
	stmt := `SELECT id,username, firstname, lastname FROM users
	WHERE id IN (
	  SELECT following_user_id FROM follower
	  WHERE followed_user_id = ? AND follow_status = 0
	)`
	rows, _ := DB.Query(stmt, userId)

	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		u := &models.User{}
		err := rows.Scan(&u.Id, &u.Username, &u.FirstName, &u.LastName)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}