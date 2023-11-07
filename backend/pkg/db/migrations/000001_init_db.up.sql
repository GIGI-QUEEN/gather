CREATE TABLE  allowed_post_users(
		user_id INTEGER NOT NULL,
		post_id INTEGER,
		FOREIGN KEY("post_id") REFERENCES "posts"("id"),
		FOREIGN KEY("user_id") REFERENCES "users"("id")
	);
	CREATE TABLE  users(
		id INTEGER not null primary key autoincrement,
		firstname TEXT not null,
		lastname TEXT not null,
		age INTEGER not null,
		gender TEXT not null,
		username TEXT not null unique,
		email TEXT not null unique,
		password BLOB not null,
		created_date INTEGER not null,
		status INTEGER NOT NULL DEFAULT 0,
		about TEXT NOT NULL,
		privacy TEXT NOT NULL DEFAULT "public",
		avatar_path TEXT NOT NULL DEFAULT "/images/avatars/no-avatar.png"
	);
	

	CREATE TABLE  posts(
		id    INTEGER NOT NULL primary key autoincrement,
		title    TEXT NOT NULL,
		contents    TEXT NOT NULL,
		create_date    INTEGER NOT NULL,
		user_id    INTEGER NOT NULL,
		privacy    TEXT NOT NULL DEFAULT 'public',
		img_path    TEXT NOT NULL DEFAULT ""
	);
	CREATE TABLE  follower(
		following_user_id INTEGER NOT NULL,
		followed_user_id INTEGER NOT NULL,
		follow_status INTEGER NOT NULL DEFAULT 0,
		FOREIGN KEY (following_user_id) REFERENCES users(id)
		FOREIGN key (followed_user_id) REFERENCES users(id)
	);
	CREATE TABLE   posts_categories(
		post_id int not null,
		category_id int not null
	);
	CREATE TABLE   categories(
		id INTEGER primary key autoincrement,
		name TEXT not null
	);
        INSERT INTO categories (id, name) VALUES (1, 'Cars');
        INSERT INTO categories (id, name) VALUES (2, 'Animals');
        INSERT INTO categories (id, name) VALUES (3, 'Art');
        INSERT INTO categories (id, name) VALUES (4, 'Games');
        INSERT INTO categories (id, name) VALUES (5, 'Movies');
        INSERT INTO categories (id, name) VALUES (6, 'Misc');

	CREATE TABLE   likes(
		post_id INTEGER not null,
		user_id INTEGER not null
	);
	CREATE TABLE   dislikes(
		post_id INTEGER not null,
		user_id TEXT not null
	);
	CREATE TABLE   comments(
		id    INTEGER NOT NULL primary key autoincrement,
		post_id    INTEGER NOT NULL,
		content    TEXT NOT NULL,
		user_id    TEXT NOT NULL,
		img_path    TEXT NOT NULL DEFAULT ""
	);
	CREATE TABLE   comment_likes(
		id integer not null constraint comment_likes_pk primary key autoincrement,
		comment_id integer not null,
		user_id integer not null
	);
	CREATE TABLE   comment_dislikes(
		id INTEGER not null constraint comment_dislikes_pk primary key autoincrement,
		comment_id INTEGER not null,
		user_id INTEGER not null
	);
	CREATE TABLE   sessions(
		id TEXT not null primary key,
		user_id INTEGER not null unique,
		created_date INTEGER not null
	);
	CREATE TABLE   groups(
		group_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		admin_user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		created_date INTEGER not null,
		FOREIGN KEY (admin_user_id) REFERENCES users(id)
	);
	CREATE TABLE   group_posts(
		post_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		post_title TEXT NOT NULL,
		post_content TEXT NOT NULL,
		group_id INTEGER,
		user_id INTEGER,
		created_date INTEGER not null,
		img_path    TEXT NOT NULL DEFAULT "",
		FOREIGN KEY (group_id) REFERENCES groups(group_id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	CREATE TABLE  group_post_likes(
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL
	);
	CREATE TABLE  group_post_dislikes(
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL
	);
	CREATE TABLE   group_post_comments(
		comment_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		comment_content TEXT NOT NULL,
		group_post_id INTEGER,
		user_id INTEGER,
		is_unread INTEGER NOT NULL,
		img_path    TEXT NOT NULL DEFAULT "",
		created_date  INTEGER not null
	);
	CREATE TABLE group_post_comment_dislikes(
		id INTEGER not null constraint group_post_comment_dislikes_pk primary key autoincrement,
		comment_id INTEGER not null,
		user_id INTEGER not null
	);
	CREATE TABLE group_post_comment_likes(
		id INTEGER not null constraint group_post_comment_dislikes_pk primary key autoincrement,
		comment_id INTEGER not null,
		user_id INTEGER not null
	);
	CREATE TABLE   group_users(
		group_id INTEGER NOT NULL,
		group_user_id INTEGER NOT NULL,
		is_admin INTEGER NOT NULL,
		is_approved INTEGER NOT NULL,
		created_date INTEGER NOT NULL,
		FOREIGN KEY (group_id) REFERENCES groups(group_id),
		FOREIGN KEY (group_user_id) REFERENCES users(id)
	);
	CREATE TABLE  group_events(
		event_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		group_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		event_date INTEGER NOT NULL,
		created_date INTEGER NOT NULL,
		FOREIGN KEY (group_id) REFERENCES groups(group_id) ON DELETE CASCADE
	);

	CREATE TABLE group_events_members_dependency (
		event_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		going_decision INTEGER DEFAULT 0,
		is_seen INTEGER DEFAULT 0,
		FOREIGN KEY (event_id) REFERENCES group_events(event_id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	CREATE TABLE conversation_rooms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
	CREATE TABLE conversation_messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender_id INTEGER NOT NULL,
		recipient_id INTEGER NOT NULL,
		conversation_id INTEGER NOT NULL,
		message TEXT NOT NULL,
		created_date INTEGER NOT NULL,
		FOREIGN KEY(sender_id) REFERENCES users(id),
		FOREIGN KEY(recipient_id) REFERENCES users(id),
		FOREIGN KEY(conversation_id) REFERENCES conversation_rooms(id)
	);
	CREATE TABLE conversation_rooms_users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		conversation_id INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id),
		FOREIGN KEY(conversation_id) REFERENCES conversation_rooms(id)
	);