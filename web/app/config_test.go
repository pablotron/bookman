package app

import (
  "reflect"
  "testing"
)

func TestNewConfigFromEnv(t *testing.T) {
  var tests = []struct {
    name string // test name
    env map[string]string // test env vars
    exp Config // expected config
  } {{
    name: "default",
    exp: Config {
      PasswordPath: "/run/secrets/bookman_web_password",
      Dsn: "host=db dbname=bookman user=bookman_web",
      HttpAddr: ":3000",
    },
  }, {
    name: "password",
    env: map[string]string {
      "BOOKMAN_PASSWORD_PATH": "foo bar baz",
    },
    exp: Config {
      PasswordPath: "foo bar baz",
      Dsn: "host=db dbname=bookman user=bookman_web",
      HttpAddr: ":3000",
    },
  }, {
    name: "dsn",
    env: map[string]string {
      "BOOKMAN_DATABASE_DSN": "foo bar baz",
    },
    exp: Config {
      PasswordPath: "/run/secrets/bookman_web_password",
      Dsn: "foo bar baz",
      HttpAddr: ":3000",
    },
  }, {
    name: "dsn",
    env: map[string]string {
      "BOOKMAN_HTTP_ADDRESS": "foo bar baz",
    },
    exp: Config {
      PasswordPath: "/run/secrets/bookman_web_password",
      Dsn: "host=db dbname=bookman user=bookman_web",
      HttpAddr: "foo bar baz",
    },
  }}

  for _, test := range(tests) {
    t.Run(test.name, func(t *testing.T) {
      // set custom env vars
      for k, v := range(test.env) {
        t.Setenv(k, v)
      }

      // load config from environment
      got := NewConfigFromEnv()

      if !reflect.DeepEqual(got, test.exp) {
        t.Fatalf("got %#v, exp %#v", got, test.exp)
      }
    })
  }
}
