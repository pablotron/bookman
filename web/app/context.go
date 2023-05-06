// Application config, context, and database pool
package app

import (
  "bookman/model"
  "context"
  "github.com/jackc/pgx/v5/pgxpool"
  "os"
)

// Application context
type Context struct {
  // configuration
  Config Config

  // Storage model
  Model model.Model

  // database pool
  Pool *pgxpool.Pool
}

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

// create new application context
func NewContext(ctx context.Context, config Config) (*Context, error) {
  // create pool
  pool, err := newPool(ctx, config)
  if err != nil {
    return nil, err
  }

  // return application context
  return &Context {
    Config: config,
    Model: model.NewDbModel(),
    Pool: pool,
  }, nil
}

