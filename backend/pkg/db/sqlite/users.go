package sqlite

import (
	"database/sql"
	"errors"
	"social-network/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

// TEST FUNCTIONS FOR services/ws-useronline-test.go
// ===========================================================================

func CheckUsersOnline() ([]*models.User, error) {

	rows, err := DB.Query("select id, user_id from sessions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var onlineUsers []*models.User

	for rows.Next() {
		session := &models.Session{}
		session.User = &models.User{}
		err := rows.Scan(&session.Id, &session.User.Id)
		if err != nil {
			return nil, err
		}
		session.User, err = GetUserIdAndUsername(session.User.Id)
		if err != nil {
			return nil, err
		}
		onlineUsers = append(onlineUsers, session.User)
	}
	return onlineUsers, nil
}

func GetUserIdAndUsername(id int) (*models.User, error) {
	row := DB.QueryRow("select id, username from users where id = ?", id)
	u := &models.User{}
	err := row.Scan(&u.Id, &u.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}

// =========================================================================== TILL HERE

func GetAllUsers() ([]*models.User, error) {
	rows, err := DB.Query("select id, username, firstname, lastname, status from Users order by username")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*models.User

	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName,&user.Status)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return users, nil
	} else {
		return nil, models.ErrNoRecord
	}
}

func InsertUser(user models.User) (int, error) {
	firstName := user.FirstName
	lastName := user.LastName
	age := user.Age
	gender := user.Gender
	userName := user.Username
	email := user.Email
	about := user.About
	password := user.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	if userName != "" {
		row := DB.QueryRow("select id from users where username = ?", userName)
		err = row.Scan()
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrDuplicateUsername
		}
	}

	row := DB.QueryRow("select id from users where email = ?", email)
	err = row.Scan()
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, models.ErrDuplicateEmail
	}

	result, err := DB.Exec("insert into users (firstname, lastname, age, gender, username, email, about, password, created_date) values (?,?,?,?,?,?,?,?, strftime('%s','now'))",
		firstName, lastName, age, gender, userName, email, about, string(hashedPassword))
	if err != nil {
		return 0, err
	}
	userId, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(userId), nil
}

func Authenticate(credName, password string) (int, error) {
	var id int
	var hashedPassword []byte
	row := DB.QueryRow("select id, password from users where email = ? or username = ?", credName, credName)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func GetUserProfile(profile_id, user_id int) (*models.User, error) {

	row := DB.QueryRow("select id,firstname, lastname, age, gender, username, email, created_date, status, about, privacy, avatar_path from users where id = ?", profile_id)
	u := &models.User{}
	err := row.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Age, &u.Gender, &u.Username, &u.Email, &u.Created, &u.Status, &u.About, &u.Privacy, &u.Avatar)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	u.Posts, err = GetAllUserPosts(profile_id, user_id)
	u.Followers, err = GetUserFollowers(profile_id)
	u.Followings, err = GetUserFollowings(profile_id)
	u.FollowRequests, err = GetUserFollowRequests(user_id)
	u.MemberInGroups, err = GetUserGroupsWhereIsMember(profile_id)
	u.FollowRequested, err = CheckFollowRequested(profile_id, user_id)
	return u, nil
}

func GetUserForPostInfo(id int) (*models.User, error) {
	row := DB.QueryRow("select id, username, firstname, lastname from users where id = ?", id)
	u := &models.User{}
	err := row.Scan(&u.Id, &u.Username, &u.FirstName, &u.LastName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}

func SetUserStatusOnline(id int) error {
	_, err := DB.Exec("update users set status = 1 where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func SetUserStatusOffline(id int) error {
	_, err := DB.Exec("update users set status = 0 where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func ChangeAccountAbout(id int, about string) error {
	_, err := DB.Exec("UPDATE users SET about = ? where id = ?", about, id)
	if err != nil {
		return err
	}
	return err
}

func ChangeAccountPrivacy(id int, privacy string) error {
	_, err := DB.Exec("update users set privacy = ? where id = ?", privacy, id)
	if err != nil {
		return err
	}
	return err
}

func ChangeAccountPrivate(id int) error {
	_, err := DB.Exec("UPDATE users SET privacy = \"private\" WHERE id = ?", id)
	if err != nil {
		return err
	}
	return err
}

func ChangeAccountPublic(id int) error {
	_, err := DB.Exec("UPDATE users SET privacy = \"public\" WHERE id = ?", id)
	if err != nil {
		return err
	}
	return err
}

func ChangeUsername(id int, username string) error {
	_, err := DB.Exec(`UPDATE users SET username ="`+username+`" where id = ?`, id)
	if err != nil {
		return err
	}
	return nil
}

func CheckFollowRequested(followedId, followingId int) (bool, error) {
	row := DB.QueryRow("SELECT follow_status from follower WHERE (followed_user_id, following_user_id, follow_status) = (?,?,0)", followedId, followingId)
	followStatus := -1
	err := row.Scan(&followStatus)
	if err != nil {
		return false, nil
	}
	if followStatus == 0 {
		return true, nil
	}
	return false, nil
}
