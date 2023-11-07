package sqlite

import (
	"database/sql"
	"errors"
	"social-network/pkg/models"
	"strings"
)

func CheckForUnreadCommentsOnMyPostInGroup(userId int) ([]*models.GroupPostCommentNotification, error) {
	stmt := `
		SELECT 	
			group_posts.group_id,
			group_posts.post_id,
			group_posts.post_title,
			group_post_comments.comment_id,
			group_post_comments.comment_content,
			group_post_comments.created_date
		FROM group_post_comments
		INNER JOIN group_posts ON group_posts.post_id = group_post_comments.group_post_id
		WHERE group_posts.user_id = ? AND group_post_comments.is_unread = 1 ORDER BY group_post_comments.created_date DESC;
	`
	rows, err := DB.Query(stmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []*models.GroupPostCommentNotification
	for rows.Next() {
		gpc := &models.GroupPostCommentNotification{}
		err = rows.Scan(&gpc.GroupId, &gpc.PostId, &gpc.Title, &gpc.CommentId, &gpc.CommentContent, &gpc.CreatedDate)
		if err != nil {
			return nil, err
		}
		comments = append(comments, gpc)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func MarkAllUnreadCommentsAsRead(commentIds []int) error {
	if len(commentIds) > 0 {
		query := "UPDATE group_post_comments SET is_unread = 0 WHERE is_unread = 1 AND comment_id IN ("
		placeholders := strings.Repeat("?, ", len(commentIds)-1) + "?"
		query += placeholders + ")"

		stmt, err := DB.Prepare(query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		params := make([]interface{}, len(commentIds))
		for i, id := range commentIds {
			params[i] = id
		}

		_, err = stmt.Exec(params...)
		if err != nil {
			return err
		}
		return nil
	}
	return models.ErrNoRecord
}

func InsertGroupPostComment(postId, userId int, content string) (int, error) {
	// Check if the post is not written by the user
	var isUnread int
	err := DB.QueryRow("SELECT CASE WHEN user_id = ? THEN 0 ELSE 1 END FROM group_posts WHERE post_id = ?", userId, postId).Scan(&isUnread)
	if err != nil {
		return -1, err
	}

	// Insert the comment with the correct is_unread value
	result, err := DB.Exec("INSERT INTO group_post_comments (comment_content, group_post_id, user_id, is_unread, created_date) values (?,?,?,?, strftime('%s','now'))",
		content, postId, userId, isUnread)
	if err != nil {
		return -1, err
	}
	commentId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(commentId), err
}

func InsertGroupPostCommentImagePath(comment_id int, img_path string) error {
	_, err := DB.Exec("UPDATE group_post_comments SET img_path = ? where comment_id = ?", img_path, comment_id)
	if err != nil {
		return err
	}
	return nil
}

func GetGroupPostCommentsByPostId(postID int) ([]*models.GroupPostComment, error) {
	rows, err := DB.Query("select comment_id, comment_content, user_id, group_post_id, img_path from group_post_comments where group_post_id = ? order by comment_id desc", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []*models.GroupPostComment
	for rows.Next() {
		c := &models.GroupPostComment{User: &models.User{}}
		err = rows.Scan(&c.Id, &c.Content, &c.User.Id, &c.PostId, &c.Image)
		if err != nil {
			return nil, err
		}
		c.User, err = GetUserForPostInfo(c.User.Id)
		if err != nil {
			return nil, err
		}
		c.LikeUsers, err = GetGroupPostCommentLikes(c.Id)
		if err != nil {
			return nil, err
		}
		c.DislikeUsers, err = GetGroupPostCommentDislikes(c.Id)
		if err != nil {
			return nil, err
		}
		c.LikeCount = len(c.LikeUsers)
		c.DislikeCount = len(c.DislikeUsers)
		comments = append(comments, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func GetGroupPostCommentLikes(commentId int) ([]int, error) {
	rows, err := DB.Query("select user_id from group_post_comment_likes where comment_id = ?", commentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var likes []int
	for rows.Next() {
		userId := -1
		err = rows.Scan(&userId)
		likes = append(likes, userId)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return likes, nil
}

func GetGroupPostCommentDislikes(commentId int) ([]int, error) {
	rows, err := DB.Query("select user_id from group_post_comment_dislikes where comment_id = ?", commentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var dislikes []int
	for rows.Next() {
		userId := -1
		err = rows.Scan(&userId)
		dislikes = append(dislikes, userId)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return dislikes, nil
}

func ChangeGroupPostCommentLike(commentId int, userId int) error {
	row := DB.QueryRow("select user_id from group_post_comment_likes where (comment_id, user_id)=(?,?)", commentId, userId)
	id := -1
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err := DB.Exec("delete from group_post_comment_dislikes where (comment_id,user_id) = (?,?)", commentId, userId)
			if err != nil && err != sql.ErrNoRows {
				return err
			}

			_, err = DB.Exec("insert into group_post_comment_likes (comment_id, user_id) values (?,?)", commentId, userId)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		_, err := DB.Exec("delete from group_post_comment_likes where (comment_id,user_id) = (?,?)", commentId, userId)
		if err != nil {
			return err
		}
	}
	return nil
}

func ChangeGroupPostCommentDislike(commentId int, userId int) error {
	row := DB.QueryRow("select user_id from group_post_comment_dislikes where (comment_id, user_id)=(?,?)", commentId, userId)
	id := -1
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err := DB.Exec("delete from group_post_comment_likes where (comment_id,user_id) = (?,?)", commentId, userId)
			if err != nil && err != sql.ErrNoRows {
				return err
			}

			_, err = DB.Exec("insert into group_post_comment_dislikes (comment_id, user_id) values (?,?)", commentId, userId)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		_, err := DB.Exec("delete from group_post_comment_dislikes where (comment_id,user_id) = (?,?)", commentId, userId)
		if err != nil {
			return err
		}
	}
	return nil
}
