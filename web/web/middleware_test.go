package web

import (
  "bookman/app"
  "context"
  "net/http"
  "net/http/httptest"
  "testing"
)

func TestAppContextFromContext(t *testing.T) {
  var exp app.Context
  ctx := context.WithValue(context.TODO(), appCtxKey, &exp)

  got := appContextFromContext(ctx)

  if got != &exp {
    t.Fatalf("got %v, exp %v", got, &exp)
  }
}

// fake HTTP request
var fakeRequest = httptest.NewRequest("GET", "/", nil)

func TestAppContextMiddleware(t *testing.T) {
  // expected application context
  var exp app.Context

  // create handler function which gets the pool from the request
  // context and checks it
  check := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
    // get and check pool from request context
    got := appContextFromContext(r.Context())

    // check expected value
    if got != &exp {
      t.Fatalf("got %v, exp %v", got, &exp)
    }
  })

  // wrap handler function with pool middleware, then send it a fake
  // request
  AppContextMiddleware(&exp)(check).ServeHTTP(nil, fakeRequest)
}

func TestSecurityHeadersMiddleware(t *testing.T) {
  // test content-security-policy value
  expCsp := "foo bar"

  tests := []struct {
    key string // header key
    exp string // expected value
  } {
    { "Access-Control-Allow-Methods", "GET, POST, HEAD, OPTIONS" },
    { "Content-Security-Policy", expCsp },
    { "Cross-Origin-Opener-Policy", "same-origin" },
    { "Cross-Origin-Resource-Policy", "same-origin" },
    { "Permissions-Policy", "camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), payment=(), usb=()" },
    { "Referrer-Policy", "strict-origin-when-cross-origin" },
    { "X-Content-Type-Options", "nosniff" },
    { "X-Frame-Options", "SAMEORIGIN" },

    // check that these two are NOT set
    // (they should be handled by an upstream reverse proxy)
    { "Access-Control-Allow-Origin", "" },
    { "Strict-Transport-Security", "" },
  }

  // create response recorder
  resp := httptest.NewRecorder()

  // minimal request handler which writes a "hi" to the response body
  hi := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
    if _, err := w.Write([]byte("hi")); err != nil {
      t.Fatal(err)
    }
  })

  // wrap handler with security headers middleware, send the combination
  // a fake request, and record the response
  SecurityHeadersMiddleware(expCsp)(hi).ServeHTTP(resp, fakeRequest)

  // check response headers
  for _, test := range(tests) {
    t.Run(test.key, func(t *testing.T) {
      got := resp.Header().Get(test.key)
      if got != test.exp {
        t.Fatalf("got \"%s\", exp \"%s\"", got, test.exp)
      }
    })
  }
}
