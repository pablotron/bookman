// Data model methods
package model

import (
  "context"
  "fmt"
  "github.com/jackc/pgx/v5"
  "github.com/jackc/pgx/v5/pgxpool"
  _ "embed"
)

//go:embed sql/list.sql
var listSql string

//go:embed sql/search.sql
var searchSql string

//go:embed sql/text.sql
var textSql string

//go:embed sql/upload.sql
var uploadSql string

//go:embed sql/edit.sql
var editSql string

// book list item
type Book struct {
  Id int `db:"id" json:"id"` // book ID
  Name string `db:"name" json:"name"` // book name
  Author string `db:"author" json:"author"` // author name
  Rank float64 `db:"rank" json:"rank"` // search result rank
}

// Get a list of books.
//
// If `q` is not empty, then the book name, content, and author are
// matched against the search string, and the list of results is sorted
// by relevance.
//
// If `q` is empty, then the return value is the full list of books,
// sorted by name.
func Search(ctx context.Context, pool *pgxpool.Pool, q string) ([]Book, error) {
  if len(q) > 0 {
    // search books by query string

    // build query args
    args := pgx.NamedArgs {
      "q": q,
    }

    // exec query, get rows
    rows, err := pool.Query(ctx, searchSql, args)
    if err != nil {
      return []Book{}, fmt.Errorf("Query(): %w", err)
    }

    // build results
    books, err := pgx.CollectRows(rows, pgx.RowToStructByName[Book])
    if err != nil {
      return []Book{}, fmt.Errorf("CollectRows(): %w", err)
    }

    // return success
    return books, nil
  } else {
    // list books by name

    // exec query, get rows
    rows, err := pool.Query(ctx, listSql)
    if err != nil {
      return []Book{}, fmt.Errorf("Query(): %w", err)
    }

    // build results
    books, err := pgx.CollectRows(rows, pgx.RowToStructByName[Book])
    if err != nil {
      return []Book{}, fmt.Errorf("CollectRows(): %w", err)
    }

    // return success
    return books, nil
  }
}

// book list item
type FullBook struct {
  Id int `db:"id" json:"id"` // book ID
  Name string `db:"name" json:"name"` // book name
  Author string `db:"author" json:"author"` // author name
  Body string `db:"body" json:"body"` // book contents
}

// Get body of given book.
func Body(ctx context.Context, pool *pgxpool.Pool, id int64) (string, error) {
  // build query args
  args := pgx.NamedArgs {
    "id": id,
  }

  // exec query, get rows
  rows, err := pool.Query(ctx, textSql, args)
  if err != nil {
    return "", fmt.Errorf("Query(): %w", err)
  }

  // build results
  book, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[FullBook])
  if err != nil {
    return "", fmt.Errorf("CollectOneRow: %w", err)
  }

  // return success
  return book.Body, nil
}

// uploaded file data
type UploadedFile struct {
  Name string // book name
  Body string // book contents
}

// Upload slice of books.
func Upload(ctx context.Context, pool *pgxpool.Pool, files []UploadedFile) error {
  // begin transaction
  tx, err := pool.Begin(ctx)
  if err != nil {
    return err
  }

  for i := range(files) {
    // build query args
    args := pgx.NamedArgs {
      "name": files[i].Name,
      "body": files[i].Body,
    }

    // upload file
    _, err := tx.Exec(ctx, uploadSql, args)
    if err != nil {
      // rollback transaction
      if rollback_err := tx.Rollback(ctx); rollback_err != nil {
        // FIXME: should probably just log the rollback error
        return rollback_err
      } else {
        return err
      }
    }
  }

  // commit changes, return result
  return tx.Commit(ctx)
}

// Set the name and author of the given book.
func Edit(ctx context.Context, pool *pgxpool.Pool, id int64, name, author string) error {
  // build query args
  args := pgx.NamedArgs {
    "id": id,
    "name": name,
    "author": author,
  }

  // exec query
  _, err := pool.Exec(ctx, editSql, args)
  return err
}
