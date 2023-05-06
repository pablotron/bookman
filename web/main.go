package main

import (
  "bookman/web"
  "context"
  "github.com/jackc/pgx/v5/pgxpool"
  "net/http"
  "os"
)

// Create database pool from config.
func newPool(ctx context.Context, config Config) (*pgxpool.Pool, error) {
  // read dsn password from secrets file
  password, err := os.ReadFile(config.PasswordPath)
  if err != nil {
    return nil, err
  }

  // parse dsn as pool config
  poolConfig, err := pgxpool.ParseConfig(config.Dsn)
  if err != nil {
    return nil, err
  }

  // set password from secret
  poolConfig.ConnConfig.Password = string(password)

  // connect to pool with config
  return pgxpool.NewWithConfig(ctx, poolConfig)
}

func main() {
  // read config from env
  config := newConfigFromEnv()

  // create pool
  pool, err := newPool(context.Background(), config)
  if err != nil {
    panic(err)
  }

  // create web router
  r, err := web.NewRouter(pool)
  if err != nil {
    panic(err)
  }

  // run http server
	panic(http.ListenAndServe(config.HttpAddr, r))
}
