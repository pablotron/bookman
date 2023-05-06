package main

import (
  "bookman/app"
  "bookman/web"
  "context"
  "net/http"
)

func main() {
  // read config from env
  config := app.NewConfigFromEnv()

  // create application context from context and config
  appCtx, err := app.NewContext(context.Background(), config)
  if err != nil {
    panic(err)
  }

  // create web router
  r, err := web.NewRouter(appCtx)
  if err != nil {
    panic(err)
  }

  // run http server
	panic(http.ListenAndServe(config.HttpAddr, r))
}
