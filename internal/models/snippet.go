package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expired time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expired int) (int, error) {
	// Parameter placeholders in prepared statements vary depending on the DBMS and driver you’re using.
	// For example, the pq driver for Postgres requires a placeholder like $1 instead of ?.

	// sql standard require the interval value inside a quote (e.g IINTERVAL '7 DAYS') which make argument for the prepared statement ignored
	// so we concat it and cast it as interval
	stmt := `INSERT INTO snippet (title, content, created, expired) 
             VALUES ($1, $2, localtimestamp, (localtimestamp + ($3 || ' DAYS')::INTERVAL)) RETURNING id;`

	var id int
	err := m.DB.QueryRow(stmt, title, content, expired).Scan(&id)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `
    SELECT id, title, content, created, expired FROM snippet
    WHERE expired > localtimestamp and id = $1;
    `

	s := &Snippet{}

	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expired)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `
    SELECT id, title, content, created, expired FROM snippet
    WHERE expired > localtimestamp
    ORDER BY id DESC
    LIMIT 10;
    `

	snippets := []*Snippet{} // kenapa pake '*' ?

	rows, err := m.DB.Query(stmt)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defer
	// statement should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	// defer rows.Close()

	for rows.Next() {
		s := &Snippet{} // kenapa pake '&' ?
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expired)
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
