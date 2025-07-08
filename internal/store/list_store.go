package store

import (
	"database/sql"
)

type List struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Entries     []ListEntry `json:"entries"`
	UserID      int         `json:"user_id"`
}

type ListEntry struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	OrderIndex int    `json:"order_index"`
}

type ListStore interface {
	CreateList(list *List) (*List, error)
	GetListByID(id int64) (*List, error)
	UpdateList(list *List) error
	DeleteList(id int64) error
	GetListOwner(id int64) (int, error)
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
		`INSERT INTO lists (user_id, title, description) VALUES ($1, $2, $3) RETURNING id`

	err = tx.QueryRow(query, list.UserID, list.Title, list.Description).Scan(&list.ID)
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

	query :=
		`SELECT id, title, description, user_id FROM lists WHERE id = $1`

	err := s.db.QueryRow(query, id).Scan(&list.ID, &list.Title, &list.Description, &list.UserID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	entryQuery :=
		`SELECT id, title, order_index FROM list_entries WHERE list_id = $1 ORDER BY order_index`

	rows, err := s.db.Query(entryQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry ListEntry
		err = rows.Scan(
			&entry.ID,
			&entry.Title,
			&entry.OrderIndex,
		)
		if err != nil {
			return nil, err
		}

		list.Entries = append(list.Entries, entry)
	}

	return list, nil
}

func (s *PostgresListStore) UpdateList(list *List) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query :=
		`UPDATE lists SET title = $1, description = $2 WHERE id = $3`

	result, err := tx.Exec(query, list.Title, list.Description, list.ID)
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

	_, err = tx.Exec(`DELETE FROM list_entries WHERE list_id = $1`, list.ID)
	if err != nil {
		return err
	}

	for _, entry := range list.Entries {
		query := `
			INSERT INTO list_entries (title, order_index, list_id, user_id) 
			VALUES ($1, $2, $3, $4)`

		_, err := tx.Exec(query, entry.Title, entry.OrderIndex, list.ID, list.UserID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *PostgresListStore) DeleteList(id int64) error {
	query :=
		`DELETE from lists WHERE id = $1`

	result, err := s.db.Exec(query, id)
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

func (s *PostgresListStore) GetListOwner(id int64) (int, error) {
	var userID int

	query :=
		`SELECT user_id FROM lists WHERE id = $1`

	err := s.db.QueryRow(query, id).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
