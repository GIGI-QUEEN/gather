package sqlite

import (
	"database/sql"
	"errors"
	"log"
	"social-network/pkg/models"
	"sort"
	"strconv"
)

func GetPosts(userId int) ([]*models.Post, error) {
	stmt := `
	SELECT * FROM posts WHERE privacy = "public"
	UNION
	SELECT posts.* from posts INNER JOIN allowed_post_users on allowed_post_users.post_id = posts.id WHERE allowed_post_users.user_id = ? AND posts.privacy = "almost private"
	UNION
	SELECT posts.* FROM posts INNER JOIN follower on follower.followed_user_id = posts.user_id where follower.following_user_id = ? AND posts.privacy = "private"
	UNION
	SELECT posts.* FROM posts WHERE privacy = "private" AND user_id = ?
	ORDER BY id desc
	`
	rows, err := DB.Query(stmt, userId, userId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		post.Categories = []*models.Category{}
		post.User = &models.User{}
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.CreateDate, &post.User.Id, &post.Privacy, &post.Image)
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
	}
	return []*models.Post{}, nil
}

func GetPostById(id int) (*models.Post, error) {
	post := &models.Post{}
	post.Categories = []*models.Category{}
	post.User = &models.User{}

	row := DB.QueryRow("select id, title, contents, create_date, user_id, privacy, img_path from Posts where id = ?", id)
	err := row.Scan(&post.Id, &post.Title, &post.Content, &post.CreateDate, &post.User.Id, &post.Privacy, &post.Image)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	post.Comments, err = GetCommentsByPostId(post.Id)
	if err != nil {
		return nil, err
	}

	post.Categories, err = GetCategoriesByPost(post.Id)
	if err != nil {
		return nil, err
	}
	post.User, err = GetUserForPostInfo(post.User.Id)
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

	return post, nil
}

func GetPostsByCategory(catID int) ([]*models.Post, error) {
	rows, err := DB.Query("select post_id from posts_categories where category_id =  ? ", catID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.Id)
		if err != nil {
			return nil, err
		}
		post, err = GetPostById(post.Id)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Id > posts[j].Id
	})

	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(posts) > 0 {
		return posts, nil
	} else {
		return nil, models.ErrNoRecord
	}
}

func InsertPost(title, contents string, categories []string, userId int, privacy string, users []string) (int, error) {
	result, err := DB.Exec("insert into posts (title, contents, create_date, user_id, privacy) values (?, ?, strftime('%s','now'), ?,?)", title, contents, userId, privacy)
	if err != nil {
		return -1, err
	}
	postId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	if privacy == "almost private" {
		users = append(users, strconv.Itoa(userId))
		for _, user := range users {
			usrId, _ := strconv.Atoi(user)
			_, err = DB.Exec("INSERT INTO allowed_post_users (user_id, post_id) VALUES (?,?)", usrId, postId)
			if err != nil {
				return -1, err
			}
		}
	}

	for _, catName := range categories {
		row := DB.QueryRow("select id from categories where name = ?", catName)
		var catId int
		err := row.Scan(&catId)
		if err != nil {
			return -1, err
		}
		_, err = DB.Exec("insert into posts_categories (post_id, category_id) values (?,?);", postId, catId)
		if err != nil {
			return -1, err
		}
	}
	return int(postId), nil
}

func InsertPostImagePath(post_id int, img_path string) error {
	_, err := DB.Exec("UPDATE posts SET img_path = ? where id = ?", img_path, post_id)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUserPosts(creator_id, user_id int) ([]*models.Post, error) {
	stmt := `
	SELECT * FROM posts WHERE (privacy, user_id) = ("public", ?)
	UNION 
	SELECT posts.* from posts INNER JOIN allowed_post_users on allowed_post_users.post_id = posts.id WHERE allowed_post_users.user_id = ? AND (posts.privacy,posts.user_id) = ("almost private", ?)
	UNION
	SELECT posts.* FROM posts INNER JOIN follower on follower.followed_user_id = posts.user_id where follower.following_user_id = ? AND (posts.privacy , posts.user_id) = (("private", ?)) AND follower.follow_status = 1
	UNION
	SELECT posts.* FROM posts WHERE privacy = "private" AND user_id = ?
	ORDER BY id desc
	`
	rows, err := DB.Query(stmt, creator_id, user_id, creator_id, user_id, creator_id, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		post.Categories = []*models.Category{}
		post.User = &models.User{}
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.CreateDate, &post.User.Id, &post.Privacy, &post.Image)
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

func GetLastInsertedPostId() int {
	res := DB.QueryRow("SELECT id from posts ORDER BY id DESC LIMIT 1")
	id := 0
	err := res.Scan(&id)

	if err != nil {
		log.Println(err)
	}
	return id
}
