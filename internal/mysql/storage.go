package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ghkadim/highload_architect/internal/models"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
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

func (s *Storage) UserRegister(ctx context.Context, user models.User) (string, error) {
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

	var uuid string
	row := s.db.QueryRowContext(ctx,
		"SELECT bin_to_uuid(uuid) FROM users WHERE id = ?", id)
	if err := row.Scan(&uuid); err != nil {
		return "", fmt.Errorf("userRegister get uuid: %v", err)
	}

	return uuid, nil
}

func (s *Storage) UserGet(ctx context.Context, id string) (models.User, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT first_name, second_name, age, biography, city, password_hash FROM users WHERE uuid = uuid_to_bin(?)", id)
	user := models.User{ID: id}
	if err := row.Scan(&user.FirstName, &user.SecondName, &user.Age, &user.Biography, &user.City, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, models.UserNotFound
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

func (s *Storage) FriendAdd(ctx context.Context, userID1, userID2 string) error {
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

func (s *Storage) FriendDelete(ctx context.Context, userID1, userID2 string) error {
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
