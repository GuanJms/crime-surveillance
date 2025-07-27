package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	dbTimeout = time.Second * 3
)

var ErrNoContent error

type PostgresRepository struct {
	Conn *sql.DB
}

func NewPostgresRepository(conn *sql.DB) *PostgresRepository {
	ErrNoContent = errors.New("no content error")
	return &PostgresRepository{
		Conn: conn,
	}
}

func (p *PostgresRepository) CreateUser(user *User) error {
	if user.Password == nil {
		return fmt.Errorf("password is required")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(*user.Password), 12)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING id`

	var userID string
	err = p.Conn.QueryRowContext(ctx, query, user.Username, passwordHash, *user.Role).Scan(&userID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return fmt.Errorf("username already exists")
			}
		}
		return err
	}

	user.ID = &userID
	return nil
}

func (p *PostgresRepository) getPasswordHash(username string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT password_hash FROM users WHERE username = $1`

	var passwordHash string
	err := p.Conn.QueryRowContext(ctx, query, username).Scan(&passwordHash)
	if err != nil {
		return "", err
	}
	return passwordHash, nil
}

func (p *PostgresRepository) AuthenticateUser(cred *UserLoginCredentials) (bool, error) {
	if passwordHash, err := p.getPasswordHash(cred.Username); err != nil {
		return false, err
	} else if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(cred.Password)); err != nil {
		return false, errors.New("invalid credentials")
	}
	// successful authenticated
	return true, nil
}

func (p *PostgresRepository) GetUserInfo(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 	id, username, role, created_at, updated_at, last_login, last_activity
		FROM users
		WHERE username = $1
	`

	row := p.Conn.QueryRowContext(ctx, query, username)

	var u User
	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.LastLogin,
		&u.LastActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (p *PostgresRepository) UpdateLoginTime(user *User) error {

	if user == nil {
		return errors.New("user is empty")
	}

	if user.ID == nil {
		return errors.New("user ID is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `UPDATE users SET last_login = $1, last_activity = $2 WHERE id = $3`

	_, err := p.Conn.ExecContext(ctx, query, time.Now(), time.Now(), *user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) ChangeUserRoleTo(id string, role string) error {
	if id == "" {
		return errors.New("id should not be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `UPDATE users SET role = $1, updated_at = $2 WHERE id = $3`

	res, err := p.Conn.ExecContext(ctx, query, role, time.Now(), id)
	if err != nil {
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
