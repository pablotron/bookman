package model

import (
  "context"
  "github.com/jackc/pgx/v5/pgxpool"
)

// Mock result from Search() method.
type MockSearchResult struct {
  Books []Book
  Err   error
}

// Mock result from Body() method
type MockBodyResult struct {
  Body string
  Err  error
}

type MockModel struct {
  SearchResult MockSearchResult // Search() method result
  BodyResult MockBodyResult // Body() method result
  UploadResult error // Upload() method result
  EditResult error // Edit() method result
}

func (m *MockModel) Search(_ context.Context, _ *pgxpool.Pool, _ string) ([]Book, error) {
  return m.SearchResult.Books, m.SearchResult.Err
}

func (m *MockModel) Body(_ context.Context, _ *pgxpool.Pool, _ int64) (string, error) {
  return m.BodyResult.Body, m.BodyResult.Err
}

func (m *MockModel) Upload(_ context.Context, _ *pgxpool.Pool, _ []UploadedFile) error {
  return m.UploadResult
}

func (m *MockModel) Edit(_ context.Context, _ *pgxpool.Pool, _ int64, _, _ string) error {
  return m.EditResult
}
