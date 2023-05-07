package model

import (
  "context"
  "errors"
  "testing"
  "reflect"
)

func TestMockModelSearch(t *testing.T) {
  t.Run("pass", func(t *testing.T) {
    exp := []Book { Book { Name: "foo" } }

    m := &MockModel {
      SearchResult: MockSearchResult {
        Books: exp,
      },
    }

    got, err := m.Search(context.Background(), nil, "")
    if err != nil {
      t.Fatal(err)
    }

    if !reflect.DeepEqual(got, exp) {
      t.Fatalf("got %#v, exp %#v", got, exp)
    }
  })

  t.Run("fail", func(t *testing.T) {
    m := &MockModel {
      SearchResult: MockSearchResult {
        Err: errors.New("some error"),
      },
    }

    got, err := m.Search(context.Background(), nil, "")
    if err == nil {
      t.Fatalf("got %#v, exp err", got)
    }
  })
}

func TestMockModelBody(t *testing.T) {
  t.Run("pass", func(t *testing.T) {
    exp := "some body"

    m := &MockModel {
      BodyResult: MockBodyResult {
        Body: exp,
      },
    }

    got, err := m.Body(context.Background(), nil, 1)
    if err != nil {
      t.Fatal(err)
    }

    if got != exp {
      t.Fatalf("got %#v, exp %#v", got, exp)
    }
  })

  t.Run("fail", func(t *testing.T) {
    m := &MockModel {
      BodyResult: MockBodyResult {
        Err: errors.New("some error"),
      },
    }

    got, err := m.Body(context.Background(), nil, 1)
    if err == nil {
      t.Fatalf("got %#v, exp err", got)
    }
  })
}

func TestMockModelUpload(t *testing.T) {
  t.Run("pass", func(t *testing.T) {
    m := &MockModel {}

    if err := m.Upload(context.Background(), nil, []UploadedFile{}); err != nil {
      t.Fatal(err)
    }
  })

  t.Run("fail", func(t *testing.T) {
    m := &MockModel {
      UploadResult: errors.New("some error"),
    }

    if err := m.Upload(context.Background(), nil, []UploadedFile{}); err == nil {
      t.Fatal("got success, exp err")
    }
  })
}

func TestMockModelEdit(t *testing.T) {
  t.Run("pass", func(t *testing.T) {
    m := &MockModel {}

    if err := m.Edit(context.Background(), nil, 1, "", ""); err != nil {
      t.Fatal(err)
    }
  })

  t.Run("fail", func(t *testing.T) {
    m := &MockModel {
      EditResult: errors.New("some error"),
    }

    if err := m.Edit(context.Background(), nil, 1, "", ""); err == nil {
      t.Fatal("got success, exp err")
    }
  })
}
