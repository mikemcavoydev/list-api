package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plainText = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type UserStore interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	UpdateUser(user *User) error
	GetUserToken(scope, token string) (*User, error)
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

func (s *PostgresUserStore) CreateUser(user *User) error {
	query :=
		`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	err := s.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}

	query :=
		`SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username = $1`

	err := s.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash.hash, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) UpdateUser(user *User) error {
	query :=
		`UPDATE users SET username = $1, email = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING updated_at`

	result, err := s.db.Exec(query, user.Username, user.Email, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *PostgresUserStore) GetUserToken(scope, token string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(token))

	query :=
		`SELECT u.id, u.username, u.email, u.password_hash, u.created_at, u.updated_at
		FROM users u 
		INNER JOIN tokens t ON t.user_id = u.id
		WHERE t.hash = $1 AND t.scope = $2 AND t.expiry > $3`

	user := &User{
		PasswordHash: password{},
	}

	err := s.db.QueryRow(query, tokenHash[:], scope, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
