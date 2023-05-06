package web

import (
  "bookman/app"
  "context"
  "net/http"
)

// appCtx context key
var appCtxKey struct{}

// get application context from context
func appContextFromContext(ctx context.Context) *app.Context {
  return ctx.Value(appCtxKey).(*app.Context)
}

// HTTP middleware which stores the given application context in the
// request context.
func AppContextMiddleware(appCtx *app.Context) func(next http.Handler) http.Handler {
  return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      // create new context with appCtx from request context
      ctx := context.WithValue(r.Context(), appCtxKey, appCtx)

      // call the next handler in the chain
      next.ServeHTTP(w, r.WithContext(ctx))
    })
  }
}

// HTTP middleware which adds the security headers to all responses:
//
// * Access-Control-Allow-Methods
// * Content-Security-Policy
// * Cross-Origin-Opener-Policy
// * Cross-Origin-Resource-Policy
// * Permissions-Policy
// * Referrer-Policy
// * X-Content-Type-Options
// * X-Frame-Options
//
// Because this site might be served locally or behind a reverse proxy,
// it does not set the following headers:
//
// * Access-Control-Allow-Origin
// * Strict-Transport-Security
func SecurityHeadersMiddleware(csp string) func(next http.Handler) http.Handler {
  return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      h := w.Header()

      // add security headers to response
      h.Add("Access-Control-Allow-Methods", "GET, POST, HEAD, OPTIONS")
      h.Add("Content-Security-Policy", csp)
      h.Add("Cross-Origin-Opener-Policy", "same-origin")
      h.Add("Cross-Origin-Resource-Policy", "same-origin")
      h.Add("Permissions-Policy", "camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), payment=(), usb=()")
      h.Add("Referrer-Policy", "strict-origin-when-cross-origin")
      h.Add("X-Content-Type-Options", "nosniff")
      h.Add("X-Frame-Options", "SAMEORIGIN")

      // TODO:
      // Access-Control-Allow-Origin
      // Strict-Transport-Security

      // call the next handler in the chain
      next.ServeHTTP(w, r)
    })
  }
}
