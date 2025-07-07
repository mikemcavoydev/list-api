package store

import (
	"database/sql"
	"time"

	"github.com/mikemcavoydev/list-api/internal/tokens"
)

type PostgresTokenStore struct {
	DB *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		DB: db,
	}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
}

func (s *PostgresTokenStore) CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = s.Insert(token)
	return token, err
}

func (s *PostgresTokenStore) Insert(token *tokens.Token) error {
	query :=
		`INSERT INTO tokens (hash, user_id, expiry, scope) VALUES ($1, $2, $3, $4)`

	_, err := s.DB.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)

	return err
}

func (s *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error {
	query :=
		`DELETE FROM tokens WHERE Scope = $1 AND user_id = $2`

	_, err := s.DB.Exec(query, scope, userID)

	return err
}
