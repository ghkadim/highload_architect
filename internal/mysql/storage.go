package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

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

func toUserID(id int64) models.UserID {
	return models.UserID(strconv.FormatInt(id, 10))
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

	return toUserID(id), nil
}

func (s *Storage) UserGet(ctx context.Context, id models.UserID) (models.User, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT first_name, second_name, age, biography, city, password_hash FROM users WHERE id = ?", id)
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
		"SELECT id, first_name, second_name, age, biography, city FROM users "+
			"WHERE first_name LIKE ? AND second_name LIKE ?", firstName+"%", secondName+"%")
	if err != nil {
		return nil, fmt.Errorf("userSearch %s %s: %v", firstName, secondName, err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var id int64
		if err := rows.Scan(&id, &user.FirstName, &user.SecondName, &user.Age, &user.Biography, &user.City); err != nil {
			return nil, fmt.Errorf("userSearch %s %s: %v", firstName, secondName, err)
		}
		user.ID = toUserID(id)
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("userSearch %s %s: %v", firstName, secondName, err)
	}

	return users, nil
}

func (s *Storage) FriendAdd(ctx context.Context, userID1, userID2 models.UserID) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO friends (user1_id, user2_id) VALUES (?, ?)",
		userID1, userID2)
	if err != nil {
		return fmt.Errorf("frienAdd: %v", err)
	}
	return nil
}

func (s *Storage) FriendDelete(ctx context.Context, userID1, userID2 models.UserID) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM friends WHERE (user1_id, user2_id) = (?, ?)",
		userID1, userID2)
	if err != nil {
		return fmt.Errorf("frienDelete: %v", err)
	}
	return nil
}

func toPostID(id int64) models.PostID {
	return models.PostID(strconv.FormatInt(id, 10))
}

func (s *Storage) PostAdd(ctx context.Context, text string, author models.UserID) (models.Post, error) {
	result, err := s.db.ExecContext(ctx,
		"INSERT INTO posts (text, user_id) VALUES (?, ?)", text, author)
	if err != nil {
		return models.Post{}, fmt.Errorf("postAdd: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return models.Post{}, fmt.Errorf("postAdd: %v", err)
	}

	return models.Post{ID: toPostID(id), SequentialID: id, Text: text, AuthorID: author}, nil
}

func (s *Storage) PostUpdate(ctx context.Context, postID models.PostID, text string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE posts SET text = ? WHERE id = ?", text, postID)
	if err != nil {
		return fmt.Errorf("postUpdate: %v", err)
	}
	return nil
}

func (s *Storage) PostDelete(ctx context.Context, postID models.PostID) error {
	_, err := s.db.ExecContext(ctx,
		"DELETE FROM posts WHERE id = ?", postID)
	if err != nil {
		return fmt.Errorf("postDelete: %v", err)
	}
	return nil
}

func (s *Storage) PostGet(ctx context.Context, postID models.PostID) (models.Post, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT id, text, user_id FROM posts WHERE id = ?", postID)
	post := models.Post{ID: postID}
	var authorID int64
	if err := row.Scan(&post.SequentialID, &post.Text, &authorID); err != nil {
		if err == sql.ErrNoRows {
			return models.Post{}, models.ErrPostNotFound
		}
		return models.Post{}, fmt.Errorf("postGet %q: %v", postID, err)
	}
	post.AuthorID = toUserID(authorID)
	return post, nil
}

func (s *Storage) PostFeed(ctx context.Context, userID models.UserID, offset, limit int) ([]models.Post, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, text, user_id FROM posts "+
			"WHERE user_id IN ("+
			"	SELECT user2_id FROM friends"+
			"	WHERE user1_id = ?)"+
			"ORDER BY posts.id DESC LIMIT ? OFFSET ?", userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("postFeed: %v", err)
	}
	defer rows.Close()

	posts := make([]models.Post, 0)
	for rows.Next() {
		var post models.Post
		var authorID int64
		if err := rows.Scan(&post.SequentialID, &post.Text, &authorID); err != nil {
			return nil, fmt.Errorf("postFeed user %s: %v", userID, err)
		}
		post.ID = toPostID(post.SequentialID)
		post.AuthorID = toUserID(authorID)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postFeed user %s: %v", userID, err)
	}

	return posts, nil
}

func (s *Storage) UserPosts(ctx context.Context, user models.UserID, limit int) ([]models.Post, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, text FROM posts "+
			"WHERE user_id = ? "+
			"ORDER BY id DESC LIMIT ?", user, limit)
	if err != nil {
		return nil, fmt.Errorf("userPosts: %v", err)
	}
	defer rows.Close()

	posts := make([]models.Post, 0)
	for rows.Next() {
		post := models.Post{AuthorID: user}
		if err := rows.Scan(&post.SequentialID, &post.Text); err != nil {
			return nil, fmt.Errorf("userPosts user %s: %v", user, err)
		}
		post.ID = toPostID(post.SequentialID)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("userPosts user %s: %v", user, err)
	}

	return posts, nil
}

func (s *Storage) UserFriends(ctx context.Context, user models.UserID) ([]models.UserID, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT user2_id FROM friends WHERE user1_id = ?", user)
	if err != nil {
		return nil, fmt.Errorf("userFriends: %v", err)
	}
	defer rows.Close()

	users := make([]models.UserID, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("userFriends user %s: %v", user, err)
		}
		users = append(users, toUserID(id))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("userFriends user %s: %v", user, err)
	}

	return users, nil
}

const dialogMessageMaxSubID = 10

func dialogMessageShardingID(userID1, userID2 models.UserID) string {
	if strings.Compare(string(userID1), string(userID2)) < 0 {
		return string(userID1) + "/" + string(userID2)
	} else {
		return string(userID2) + "/" + string(userID1)
	}
}

func (s *Storage) DialogSend(ctx context.Context, message models.DialogMessage) error {
	shardingID := dialogMessageShardingID(message.From, message.To)
	shardingSubID := rand.Intn(dialogMessageMaxSubID)
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO dialogs (from_user_id, to_user_id, text, sharding_id, sharding_sub_id) VALUES (?,?,?,?,?)",
		message.From, message.To, message.Text, shardingID, shardingSubID)
	if err != nil {
		return fmt.Errorf("dialogSend: %v", err)
	}
	return nil
}

func (s *Storage) DialogList(ctx context.Context, userID1, userID2 models.UserID) ([]models.DialogMessage, error) {
	shardingID := dialogMessageShardingID(userID1, userID2)
	rows, err := s.db.QueryContext(ctx,
		"SELECT from_user_id, to_user_id, text FROM dialogs "+
			"WHERE sharding_id = ? ORDER BY id DESC", shardingID)
	if err != nil {
		return nil, fmt.Errorf("dialogList %s %s: %v", userID1, userID2, err)
	}
	defer rows.Close()

	var messages []models.DialogMessage
	for rows.Next() {
		var message models.DialogMessage
		var from_id, to_id int64
		if err := rows.Scan(&from_id, &to_id, &message.Text); err != nil {
			return nil, fmt.Errorf("dialogList %s %s: %v", userID1, userID2, err)
		}
		message.From = toUserID(from_id)
		message.To = toUserID(to_id)
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("dialogList %s %s: %v", userID1, userID2, err)
	}

	return messages, nil
}
