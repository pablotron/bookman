// web interface
package web

import (
  "bookman/model"
  "embed"
  "encoding/json"
  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
  "github.com/jackc/pgx/v5/pgxpool"
  "io"
  io_fs "io/fs"
  "net/http"
  "strconv"
  "strings"
)

// Get a list of books.
//
// If the `q` request parameter is not empty, then the returned list of
// books is a list of books which match the query string, sorted in
// descending order by relevance.
//
// If the `q` request parameter is empty, then the returned list of
// books is a complete list of books, sorted by name.
func doApiSearch(w http.ResponseWriter, r *http.Request) {
  // get context from request and pool from context
  ctx := r.Context()
  pool := poolFromContext(ctx)

  // set response header
  w.Header().Add("Content-Type", "text/json")

  // get books
  books, err := model.Search(ctx, pool, r.FormValue("q"))
  if err != nil {
    panic(err)
  }

  // write JSON-encoded list of books
  if err := json.NewEncoder(w).Encode(books); err != nil {
    panic(err)
  }
}

// Route handler which shows contents of given book.
func doBook(w http.ResponseWriter, r *http.Request) {
  // get context from request and pool from context
  ctx := r.Context()
  pool := poolFromContext(ctx)

  // parse book ID
  bookId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
  if err != nil {
    panic(err)
  }

  // get book body
  body, err := model.Body(ctx, pool, bookId)
  if err != nil {
    panic(err)
  }

  // set response header, write body
  w.Header().Add("Content-Type", "text/plain")
  if _, err := w.Write([]byte(body)); err != nil {
    panic(err)
  }
}

// Route handler for file uploads
func doApiUpload(w http.ResponseWriter, r *http.Request) {
  // get context from request and pool from context
  ctx := r.Context()
  pool := poolFromContext(ctx)

  // get multipart reader from request
  mpr, err := r.MultipartReader()
  if err != nil {
    panic(err)
  }

  // build list of uploaded files
  var files []model.UploadedFile
  for {
    part, err := mpr.NextPart()
    if err == io.EOF {
      break;
    } else if err != nil {
      panic(err)
    }

    // read part data
    data, err := io.ReadAll(part)
    if err != nil {
      panic(err)
    }

    // add to list of files
    files = append(files, model.UploadedFile {
      Name: strings.TrimSuffix(part.FileName(), ".txt"),
      Body: string(data),
    })
  }

  // upload files
  if err := model.Upload(ctx, pool, files); err != nil {
    panic(err)
  }

  // send response
  w.Header().Add("Content-Type", "text/json")
  if _, err := w.Write([]byte("null")); err != nil {
    panic(err)
  }
}

// Edit book route handler.
func doApiEdit(w http.ResponseWriter, r *http.Request) {
  // get context from request and pool from context
  ctx := r.Context()
  pool := poolFromContext(ctx)

  // parse book ID
  id, err := strconv.ParseInt(r.FormValue("id"), 10, 32)
  if err != nil {
    panic(err)
  }

  // get new name and author
  name := r.FormValue("name")
  author := r.FormValue("author")

  // edit book
  if err := model.Edit(ctx, pool, id, name, author); err != nil {
    panic(err)
  }

  // send response
  w.Header().Add("Content-Type", "text/json")
  if _, err := w.Write([]byte("null")); err != nil {
    panic(err)
  }
}

// Route handler which panics.
func doApiPanic(w http.ResponseWriter, r *http.Request) {
  panic("this is a test panic")
}

//go:embed public
var publicFs embed.FS

// list of content types to compress with the Compress middleware
var compressContentTypes = []string {
  "text/html",
  "text/plain",
  "text/css",
  "text/javascript",
  "text/json",
}

// Content security policy.
//
// We allow `data:` images because the favicon is a `data:` URL.
//
// Passed to the SecurityHeaders middleware.
var contentSecurityPolicy = "default-src 'self'; img-src 'self' data:"

func NewRouter(pool *pgxpool.Pool) (*chi.Mux, error) {
  // get public directory
  public, err := io_fs.Sub(publicFs, "public")
  if err != nil {
    return nil, err
  }

  // create router, attach middleware
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
  r.Use(middleware.Compress(5, compressContentTypes...))
  r.Use(SecurityHeadersMiddleware(contentSecurityPolicy))
	r.Use(PoolMiddleware(pool))

  // bind routes
	r.Get("/api/search", doApiSearch)
	r.Get("/api/panic", doApiPanic)
	r.Post("/api/upload", doApiUpload)
	r.Post("/api/edit", doApiEdit)
	r.Get("/book/{id:^\\d+$}", doBook)
  // bind static site (note the "/*" to match all files)
	r.Handle("/*", http.FileServer(http.FS(public)))

  // return router
  return r, nil
}
