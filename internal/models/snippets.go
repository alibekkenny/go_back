package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Define a Snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table?
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *pgxpool.Pool
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES($1, $2, now(), NOW() + make_interval(days=>$3)) RETURNING id`
	var returnedId int
	err := m.DB.QueryRow(context.Background(), stmt, title, content, expires).Scan(&returnedId)
	if err != nil {
		return 0, err
	}
	fmt.Println(returnedId)
	// rowsAffected := result.RowsAffected()
	return int(returnedId), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	returnedInfo := &Snippet{}
	stmt := `SELECT * FROM snippets WHERE id = $1`
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&returnedInfo.ID, &returnedInfo.Title, &returnedInfo.Content, &returnedInfo.Created, &returnedInfo.Expires)
	if err != nil {
		return nil, err
	}
	return returnedInfo, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := "SELECT * FROM snippets ORDER BY id DESC LIMIT 10"
	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
