package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db error: %v", err)
	}

	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("migrating test db error: %v", err)
	}

	_, err = db.Exec(`TRUNCATE lists, list_entries CASCADE`)
	if err != nil {
		t.Fatalf("truncating tables error: %v", err)
	}

	return db
}

func TestCreateList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresListStore(db)

	tests := []struct {
		name    string
		list    *List
		wantErr bool
	}{
		{
			name: "valid list",
			list: &List{
				Title:       "Test list",
				Description: "Test description",
				Entries: []ListEntry{
					{
						Title:      "Test entry",
						OrderIndex: 0,
					},
					{
						Title:      "Test entry 2",
						OrderIndex: 0,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdList, err := store.CreateList(tt.list)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.list.Title, createdList.Title)
			assert.Equal(t, tt.list.Description, createdList.Description)

			retrieved, err := store.GetListByID(int64(createdList.ID))
			require.NoError(t, err)

			assert.Equal(t, createdList.ID, retrieved.ID)
			assert.Equal(t, len(tt.list.Entries), len(retrieved.Entries))

			for i, entry := range retrieved.Entries {
				assert.Equal(t, tt.list.Entries[i].Title, entry.Title)
				assert.Equal(t, tt.list.Entries[i].OrderIndex, entry.OrderIndex)
			}
		})
	}
}
