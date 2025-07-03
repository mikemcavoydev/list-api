package store

import (
	"database/sql"
)

type List struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Entries     []ListEntry `json:"entries"`
}

type ListEntry struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	OrderIndex int    `json:"order_index"`
}

type ListStore interface {
	CreateList(list *List) (*List, error)
	GetListByID(id int64) (*List, error)
}

type PostgresListStore struct {
	db *sql.DB
}

func NewPostgresListStore(db *sql.DB) *PostgresListStore {
	return &PostgresListStore{db: db}
}

func (s *PostgresListStore) CreateList(list *List) (*List, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query :=
		`INSERT INTO lists (title, description) VALUES ($1, $2) RETURNING id`

	err = tx.QueryRow(query, list.Title, list.Description).Scan(&list.ID)
	if err != nil {
		return nil, err
	}

	for _, entry := range list.Entries {
		query :=
			`INSERT INTO list_entries (list_id, title, order_index) VALUES ($1, $2, $3) RETURNING id`
		err = tx.QueryRow(query, list.ID, entry.Title, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *PostgresListStore) GetListByID(id int64) (*List, error) {
	list := &List{}
	return list, nil
}
