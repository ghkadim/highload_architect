package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/ghkadim/highload_architect/internal/models"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(
	user string,
	password string,
	address string,
	database string,
) (*Storage, error) {
	cfg := mysql.Config{
		User:   user,
		Passwd: password,
		Net:    "tcp",
		Addr:   address,
		DBName: database,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) UserRegister(ctx context.Context, user models.User) (models.UserID, error) {
	result, err := s.db.ExecContext(ctx,
		"INSERT INTO users (first_name, second_name, age, biography, city, password_hash) VALUES (?, ?, ?, ?, ?, ?)",
		user.FirstName, user.SecondName, user.Age, user.Biography, user.City, user.PasswordHash)
	if err != nil {
		return "", fmt.Errorf("userRegister: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("userRegister: %v", err)
	}

	var uuid models.UserID
	row := s.db.QueryRowContext(ctx,
		"SELECT bin_to_uuid(uuid) FROM users WHERE id = ?", id)
	if err := row.Scan(&uuid); err != nil {
		return "", fmt.Errorf("userRegister get uuid: %v", err)
	}

	return uuid, nil
}

func (s *Storage) UserGet(ctx context.Context, id models.UserID) (models.User, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT first_name, second_name, age, biography, city, password_hash FROM users WHERE uuid = uuid_to_bin(?)", id)
	user := models.User{ID: id}
	if err := row.Scan(&user.FirstName, &user.SecondName, &user.Age, &user.Biography, &user.City, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("userGet %q: %v", id, err)
	}

	return user, nil
}

func (s *Storage) UserSearch(ctx context.Context, firstName, secondName string) ([]models.User, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT bin_to_uuid(uuid), first_name, second_name, age, biography, city FROM users "+
			"WHERE first_name LIKE ? AND second_name LIKE ?", firstName+"%", secondName+"%")
	if err != nil {
		return nil, fmt.Errorf("userSearch %s %s: %v", firstName, secondName, err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.SecondName, &user.Age, &user.Biography, &user.City); err != nil {
			return nil, fmt.Errorf("userSearch %s %s: %v", firstName, secondName, err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("userSearch %s %s: %v", firstName, secondName, err)
	}

	return users, nil
}

func (s *Storage) FriendAdd(ctx context.Context, userID1, userID2 models.UserID) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO friends (user1_id, user2_id) VALUES ("+
			"(SELECT id FROM users WHERE uuid = uuid_to_bin(?)),"+
			"(SELECT id FROM users WHERE uuid = uuid_to_bin(?)))",
		userID1, userID2)
	if err != nil {
		return fmt.Errorf("frienAdd: %v", err)
	}
	return nil
}

func (s *Storage) FriendDelete(ctx context.Context, userID1, userID2 models.UserID) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM friends WHERE (user1_id, user2_id) = ("+
			"(SELECT id FROM users WHERE uuid = uuid_to_bin(?)),"+
			"(SELECT id FROM users WHERE uuid = uuid_to_bin(?)))",
		userID1, userID2)
	if err != nil {
		return fmt.Errorf("frienDelete: %v", err)
	}
	return nil
}

func (s *Storage) PostAdd(ctx context.Context, text string, author models.UserID) (models.Post, error) {
	result, err := s.db.ExecContext(ctx,
		"INSERT INTO posts (text, user_id) VALUES "+
			"(?, (SELECT id FROM users WHERE uuid = uuid_to_bin(?)))", text, author)
	if err != nil {
		return models.Post{}, fmt.Errorf("postAdd: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return models.Post{}, fmt.Errorf("postAdd: %v", err)
	}

	var uuid models.PostID
	var seqID int64
	row := s.db.QueryRowContext(ctx,
		"SELECT bin_to_uuid(uuid), id FROM posts WHERE id = ?", id)
	if err := row.Scan(&uuid, &seqID); err != nil {
		return models.Post{}, fmt.Errorf("postAdd get uuid: %v", err)
	}

	return models.Post{ID: uuid, SequentialID: seqID, Text: text, AuthorID: author}, nil
}

func (s *Storage) PostUpdate(ctx context.Context, postID models.PostID, text string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE posts SET text = ? WHERE uuid = uuid_to_bin(?)", text, postID)
	if err != nil {
		return fmt.Errorf("postUpdate: %v", err)
	}
	return nil
}

func (s *Storage) PostDelete(ctx context.Context, postID models.PostID) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM posts WHERE uuid = uuid_to_bin(?)", postID)
	if err != nil {
		return fmt.Errorf("postDelete: %v", err)
	}
	return nil
}

func (s *Storage) PostGet(ctx context.Context, postID models.PostID) (models.Post, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT posts.id, text, bin_to_uuid(u.uuid) FROM posts "+
			"JOIN users u ON u.id = posts.user_id WHERE posts.uuid = uuid_to_bin(?)", postID)
	post := models.Post{ID: postID}
	if err := row.Scan(&post.SequentialID, &post.Text, &post.AuthorID); err != nil {
		if err == sql.ErrNoRows {
			return models.Post{}, models.ErrPostNotFound
		}
		return models.Post{}, fmt.Errorf("postGet %q: %v", postID, err)
	}
	return post, nil
}

func (s *Storage) PostFeed(ctx context.Context, userID models.UserID, offset, limit int) ([]models.Post, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT posts.id, bin_to_uuid(posts.uuid), text, bin_to_uuid(u.uuid) FROM posts "+
			"JOIN users u ON u.id = posts.user_id "+
			"WHERE u.id IN ("+
			"	SELECT user2_id FROM friends"+
			"	JOIN users u ON u.id = friends.user1_id "+
			"	WHERE u.uuid = uuid_to_bin(?))"+
			"ORDER BY posts.id DESC LIMIT ? OFFSET ?", userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("postFeed: %v", err)
	}
	defer rows.Close()

	posts := make([]models.Post, 0)
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.SequentialID, &post.ID, &post.Text, &post.AuthorID); err != nil {
			return nil, fmt.Errorf("postFeed user %s: %v", userID, err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postFeed user %s: %v", userID, err)
	}

	return posts, nil
}

func (s *Storage) UserPosts(ctx context.Context, user models.UserID, limit int) ([]models.Post, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT posts.id, bin_to_uuid(posts.uuid), text FROM posts "+
			"JOIN users u ON u.id = posts.user_id WHERE u.uuid = uuid_to_bin(?) "+
			"ORDER BY id DESC LIMIT ?", user, limit)
	if err != nil {
		return nil, fmt.Errorf("userPosts: %v", err)
	}
	defer rows.Close()

	posts := make([]models.Post, 0)
	for rows.Next() {
		post := models.Post{AuthorID: user}
		if err := rows.Scan(&post.SequentialID, &post.ID, &post.Text); err != nil {
			return nil, fmt.Errorf("userPosts user %s: %v", user, err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("userPosts user %s: %v", user, err)
	}

	return posts, nil
}

func (s *Storage) UserFriends(ctx context.Context, user models.UserID) ([]models.UserID, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT bin_to_uuid(u.uuid) FROM friends "+
			"JOIN users u ON u.id = user2_id "+
			"WHERE user2_id IN ( "+
			" 	SELECT user2_id FROM friends "+
			" 	JOIN users u ON u.id = friends.user1_id "+
			" 	WHERE u.uuid = uuid_to_bin(?))", user)
	if err != nil {
		return nil, fmt.Errorf("userFriends: %v", err)
	}
	defer rows.Close()

	users := make([]models.UserID, 0)
	for rows.Next() {
		var id models.UserID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("userFriends user %s: %v", user, err)
		}
		users = append(users, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("userFriends user %s: %v", user, err)
	}

	return users, nil
}
