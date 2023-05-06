package app

import (
  "os"
)

// Configuration values.  Use NewConfigFromEnv() to create a new
// configuration from environment variables.
type Config struct {
  // file containing database password 
  PasswordPath string

  // database dsn
  Dsn string

  // http host and port
  HttpAddr string
}

// default configuration
var defaultConfig = Config {
  PasswordPath: "/run/secrets/bookman_web_password", // default password file path
  Dsn: "host=db dbname=bookman user=bookman_web", // default database dsn
  HttpAddr: ":3000", // default http listen address
}

// Create new configuration from environment variables
//
// Uses the following environment variables to override the default
// configuration, if they are provided:
//
// * BOOKMAN_PASSWORD_PATH: path to file containing database password
// * BOOKMAN_DATABASE_DSN: database dsn
// * BOOKMAN_HTTP_ADDR: host and port to listen for http requests
func NewConfigFromEnv() Config {
  // create config
  config := defaultConfig

  // get password file path
  passwordPath := os.Getenv("BOOKMAN_PASSWORD_PATH")
  if passwordPath != "" {
    config.PasswordPath = passwordPath
  }

  // get dsn
  dsn := os.Getenv("BOOKMAN_DATABASE_DSN")
  if dsn != "" {
    config.Dsn = dsn
  }

  // parse http host and port
  httpAddr := os.Getenv("BOOKMAN_HTTP_ADDRESS")
  if httpAddr != "" {
    config.HttpAddr = httpAddr
  }

  // return configuration
  return config
}
