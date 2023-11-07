package sqlite

import (
	"database/sql"
	"errors"
)

func ChangeGroupPostLike(postId, userId int) error {
	row := DB.QueryRow("select user_id from group_post_likes where (post_id, user_id) = (?,?)", postId, userId)
	id := -1
	err := row.Scan(&id)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			_, err := DB.Exec("delete from group_post_dislikes where (post_id,user_id) = (?,?)", postId, userId)
			if err != nil && err != sql.ErrNoRows {
				return err
			}

			_, err = DB.Exec("insert into group_post_likes (post_id, user_id) values (?,?)", postId, userId)
			if err != nil {

				return err
			}
		} else {
			return err
		}
	} else {
		_, err := DB.Exec("delete from group_post_likes where (post_id,user_id) = (?,?)", postId, userId)
		if err != nil {

			return err
		}
	}
	return nil
}

func ChangeGroupPostDislike(postId int, userId int) error {
	row := DB.QueryRow("select user_id from group_post_dislikes where (post_id, user_id) = (?,?)", postId, userId)

	id := -1
	err := row.Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err := DB.Exec("delete from group_post_likes where (post_id,user_id) = (?,?)", postId, userId)
			if err != nil && err != sql.ErrNoRows {
				return err
			}

			_, err = DB.Exec("insert into group_post_dislikes (post_id, user_id) values (?,?)", postId, userId)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		_, err := DB.Exec("delete from group_post_dislikes where (post_id,user_id) = (?,?)", postId, userId)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetGroupPostLikeUsers(postID int) ([]int, error) {
	rows, err := DB.Query("select user_id from group_post_likes where post_id = ?", postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var likes []int
	for rows.Next() {
		userId := -1
		err = rows.Scan(&userId)
		if err != nil {
			return nil, err
		}
		likes = append(likes, userId)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return likes, nil
}

func GetGroupPostDislikeUsers(postID int) ([]int, error) {
	rows, err := DB.Query("select user_id from group_post_dislikes where post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dislikes []int

	for rows.Next() {
		userId := -1
		err = rows.Scan(&userId)
		if err != nil {
			return nil, err
		}
		dislikes = append(dislikes, userId)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return dislikes, nil
}
