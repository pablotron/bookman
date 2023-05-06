# Bookman

Minimal web interface for managing text books.

The main purpose of this tool is to exercise [podman][] and
[podman-compose][].

This repository contains the following containers:

* `db`: Database server.
* `web`: Web server.

## Usage

### Generate Database Role Passwords

First, generate database role passwords and save them as [Podman][] secrets:

    # generate passwords for the `postgres` and `bookman_web` database
    # roles and save them as secrets
    for i in bookman_{postgres,web}_password; do
      # generate random password and save it as a secret
      dd if=/dev/urandom bs=16 count=1 status=none | base64 | podman secret create $i -
    end

**Note:** Reading 18 bytes from /dev/urandom and base64-encoding them
produces a 25-digit password with 144 bits of entropy, which should be
enough for anybody ;).

### Start Bookman Service

To start the containers:

    # start db and web containers
    podman-compose up -d

(If you have not run this service before, this will build and tag images
for both containers).

Once the then expose port 3000 of the `web`
container on the host.

(Note: The first time the `db` container starts up it will seed the
database with a few dozen sample books from [Project Gutenberg][], so it
may take a moment).

### Stop Bookman Service

To stop the service:

    # shut down bookman containers
    podman-compose down

## Technical Details

### Frontend

* [Bulma][] (custom build to reduce size)
* Icons from [Bootstrap Icons][]
* Minified with [tdewolf/minify][] (see `web/assets/build.rb`).

Static web assets are minified and served compressed to keep the payload size
below 20 kB.

### Backend

* [Chi][]: routing and middleware
* [pgx][]: database driver

The web server itself is a staticly-linked binary built via a [multi-stage
build][] with web assets embedded into the binary via [go embed][].  As a
result, the `web` container consists of a single 9MB `/bookman` binary:

		> cd web && podman build -t bookman-web .
		...
		> podman unshare
		> cd $(podman image mount bookman-web)
		$ ls
		bookman
		$ du -h ./bookman
		9.0M	./bookman

The web interface sets a restrictive [Content-Security-Policy][csp].  The
complete list of security-related HTTP response headers is as follows:

* `Access-Control-Allow-Methods`
* `Content-Security-Policy`
* `Cross-Origin-Opener-Policy`
* `Cross-Origin-Resource-Policy`
* `Permissions-Policy`
* `Referrer-Policy`
* `X-Content-Type-Options`
* `X-Frame-Options`

Because this site might be served locally or behind a reverse proxy,
it does not set the following headers:

* `Access-Control-Allow-Origin`
* `Strict-Transport-Security`

See `SecurityHeadersMiddleware` in `web/middleware.go` for additional
details.

This site does not use cookies, local storage, or session storage.

### Database

The database server is [Postgres][].

The underlying `bookman` database and database objects are owned by the
`bookman_sys` database role rather than `postgres`.

Queries from the web interface run as the `bookman_web` database role,
which has relatively limited privileges.

The underlying `books` table has an indexed `tsvector` column which is
generated from the name, author, and content of each uploaded book.
Searches are performed against the index.

See `db/scripts/books.txt.gz` for additional information.  Note:
`books.txt.gz` also contains the contents of seed books from [Project
Gutenberg][], so it is quite large.

[podman]: https://podman.io/
  "Docker-compatible container engine."
[podman-compose]: https://github.com/containers/podman-compose
  "Podman-compatible clone of Docker Compose."
[project gutenberg]: https://www.gutenberg.org/
  "Library of free eBooks."
[bulma]: https://bulma.io/
  "Bulma CSS framework"
[bootstrap icons]: https://icons.getbootstrap.com/
  "Bootstrap icons"
[tdewolf/minify]: https://github.com/tdewolff/minify
  "Go minification library and command-line utility."
[go embed]: https://pkg.go.dev/embed
  "Embed files in Go binaries at build time."
[chi]: https://go-chi.io/
  "Lightweight router for building Go services."
[pgx]: https://github.com/jackc/pgx
  "Pure Go Postgres database driver."
[postgres]: https://www.postgresql.org/
  "Postgres database server."
[fts]: https://www.postgresql.org/docs/current/textsearch-intro.html
  "Full Text Search (FTS)"
[csp]: https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP
  "Content Security Policy".
[multi-stage build]: https://docs.docker.com/build/building/multi-stage/
  "Multi-stage build."
