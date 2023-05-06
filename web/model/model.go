// Data model methods
package model

import (
  "context"
  "github.com/jackc/pgx/v5/pgxpool"
  _ "embed"
)

// Book search result.
type Book struct {
  Id int `db:"id" json:"id"` // book ID
  Name string `db:"name" json:"name"` // book name
  Author string `db:"author" json:"author"` // author name
  Rank float64 `db:"rank" json:"rank"` // search result rank
}

// uploaded file data
type UploadedFile struct {
  Name string // book name
  Body string // book contents
}

// Book storage model interface.
type Model interface {
  // Get a list of books.
  //
  // If `q` is not empty, then the book name, content, and author are
  // matched against the search string, and the list of results is sorted
  // by relevance.
  //
  // If `q` is empty, then the return value is the full list of books,
  // sorted by name.
  Search(ctx context.Context, pool *pgxpool.Pool, q string) ([]Book, error)

  // Get body of given book.
  Body(ctx context.Context, pool *pgxpool.Pool, id int64) (string, error)

  // Upload slice of books.
  Upload(ctx context.Context, pool *pgxpool.Pool, files []UploadedFile) error

  // Set the name and author of the given book.
  Edit(ctx context.Context, pool *pgxpool.Pool, id int64, name, author string) error
}
