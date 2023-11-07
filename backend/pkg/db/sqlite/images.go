package sqlite

func InsertUserAvatar(user_id int, path string) error {
	_, err := DB.Exec("UPDATE users SET avatar_path = ? where id = ?", path, user_id)
	if err != nil {
		return err
	}
	return nil
}

func InsertPostImage(post_id int, path string) error {
	_, err := DB.Exec("insert into posts_images (post_id, path) values (?,?)", post_id, path)
	if err != nil {
		return err
	}
	return nil
}
