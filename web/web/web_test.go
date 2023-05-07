package web

import (
  "bookman/app"
  "bookman/model"
  "context"
  "errors"
  "fmt"
  "github.com/go-chi/chi/v5"
  "io"
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"
)

func TestDoApiSearch(t *testing.T) {
  t.Run("pass", func(t *testing.T) {
    // tests expected to pass
    var passTests = []struct {
      name string       // test name
      data []model.Book // data to return
      exp  string       // expected body
    } {{
      name: "empty",
      data: []model.Book{},
      exp: `[]`,
    }, {
      name: "empty",
      data: []model.Book {
        model.Book { Id: 1, Name: "foo" },
      },
      exp: `[{"id":1,"name":"foo","author":"","rank":0}]`,
    }}

    // run pass tests
    for _, test := range(passTests) {
      t.Run(test.name, func(t *testing.T) {
        // build app context w/ mock model
        appCtx := app.Context {
          Model: &model.MockModel {
            SearchResult: model.MockSearchResult {
              Books: test.data,
            },
          },
        }

        // create context, request, and response recorder
        ctx := context.WithValue(context.Background(), appCtxKey, &appCtx)
        req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
        if err != nil {
          t.Fatal(err)
        }
        resp := httptest.NewRecorder()

        // call handler
        doApiSearch(resp, req)

        // check response content-type
        t.Run("Content-Type", func(t *testing.T) {
          exp := "text/json"
          got := resp.Header().Get("Content-Type")
          if got != exp {
            t.Fatalf("got \"%s\", exp \"%s\"", got, exp)
          }
        })

        // check response body
        t.Run("body", func(t *testing.T) {
          // read response body
          body, err := io.ReadAll(resp.Result().Body)
          if err != nil {
            t.Fatal(err)
          }

          // check response body
          got := strings.TrimSpace(string(body))
          if got != test.exp {
            t.Fatalf("got \"%s\", exp \"%s\"", got, test.exp)
          }
        })
      })
    }
  })

  // test model.Search() failure
  t.Run("model search fail", func(t *testing.T) {
    // build app context w/ mock model
    appCtx := app.Context {
      Model: &model.MockModel {
        SearchResult: model.MockSearchResult {
          Err: errors.New("some error"),
        },
      },
    }

    // create context, request, and response recorder
    ctx := context.WithValue(context.Background(), appCtxKey, &appCtx)
    req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
    if err != nil {
      t.Fatal(err)
    }
    resp := httptest.NewRecorder()

    // doApiSearch() panics on error, so recover from it and log the
    // error here.  this is probably not the best way to do things, but
    // it was good enough for testing purposes
    defer func() {
      if err := recover(); err != nil {
        // log recovered error
        t.Log(err)
      }
    }()

    // call handler
    doApiSearch(resp, req)

    // shouldn't be reached
    t.Fatal("got success, exp err")
  })

  // TODO: test JSON encode write error
}

func TestDoBook(t *testing.T) {
  // note: because doBook uses chi.URLParam() in order to extract the
  // book ID from the request URL, we need to create a mock chi router
  // with a /book endpoint in order to test doBook()
  router := chi.NewRouter()
  router.Get("/book/{id:^\\d+$}", doBook)

  t.Run("pass", func(t *testing.T) {
    // tests expected to pass
    var passTests = []struct {
      id   int64 // fake book ID
      exp  string // expected body
    } {{
      id: 1,
      exp: `foo bar`,
    }, {
      id: 2,
      exp: `bar baz`,
    }}

    // run pass tests
    for _, test := range(passTests) {
      t.Run(fmt.Sprintf("%d", test.id), func(t *testing.T) {
        // build app context w/ mock model
        appCtx := app.Context {
          Model: &model.MockModel {
            BodyResult: model.MockBodyResult {
              Body: test.exp,
            },
          },
        }

        // create context, url path, request, and response recorder
        ctx := context.WithValue(context.Background(), appCtxKey, &appCtx)
        urlPath := fmt.Sprintf("/book/%d", test.id)
        req, err := http.NewRequestWithContext(ctx, "GET", urlPath, nil)
        if err != nil {
          t.Fatal(err)
        }
        resp := httptest.NewRecorder()

        // send request
        router.ServeHTTP(resp, req)

        // check response content-type
        t.Run("Content-Type", func(t *testing.T) {
          exp := "text/plain"
          got := resp.Header().Get("Content-Type")
          if got != exp {
            t.Fatalf("got \"%s\", exp \"%s\"", got, exp)
          }
        })

        // check response body
        t.Run("body", func(t *testing.T) {
          // read response body
          body, err := io.ReadAll(resp.Result().Body)
          if err != nil {
            t.Fatal(err)
          }

          // check response body
          got := strings.TrimSpace(string(body))
          if got != test.exp {
            t.Fatalf("got \"%s\", exp \"%s\"", got, test.exp)
          }
        })
      })
    }
  })


  // test strconv.ParseInt() failure
  t.Run("parseint fail", func(t *testing.T) {
    // build app context w/ mock model
    appCtx := app.Context {
      Model: &model.MockModel {},
    }

    // create context, url path, request, and response recorder
    ctx := context.WithValue(context.Background(), appCtxKey, &appCtx)
    req, err := http.NewRequestWithContext(ctx, "GET", "/book/36893488147419103232", nil)
    if err != nil {
      t.Fatal(err)
    }
    resp := httptest.NewRecorder()

    defer func() {
      if err := recover(); err != nil {
        // log recovered error
        t.Logf("got expected error: %s", err)
      }
    }()

    // send request
    router.ServeHTTP(resp, req)

    // read response body
    body, err := io.ReadAll(resp.Result().Body)
    if err != nil {
      t.Fatal(err)
    }

    // shouldn't be reached, log response body
    t.Fatalf("got success, exp err; body = \"%s\"", string(body))
  })

  // test model.Body() failure
  t.Run("body fail", func(t *testing.T) {
    // build app context w/ mock model
    appCtx := app.Context {
      Model: &model.MockModel {
        BodyResult: model.MockBodyResult {
          Err: errors.New("some error"),
        },
      },
    }

    // create context, url path, request, and response recorder
    ctx := context.WithValue(context.Background(), appCtxKey, &appCtx)
    req, err := http.NewRequestWithContext(ctx, "GET", "/book/1", nil)
    if err != nil {
      t.Fatal(err)
    }
    resp := httptest.NewRecorder()

    defer func() {
      if err := recover(); err != nil {
        // log recovered error
        t.Logf("got expected error: %s", err)
      }
    }()

    // send request
    router.ServeHTTP(resp, req)

    // read response body
    body, err := io.ReadAll(resp.Result().Body)
    if err != nil {
      t.Fatal(err)
    }

    // shouldn't be reached, log response body
    t.Fatalf("got success, exp err; body = \"%s\"", string(body))
  })
}

// TODO: TestDoUpload()
// TODO: TestDoEdit()
