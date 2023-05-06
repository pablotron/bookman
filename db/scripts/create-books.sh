#!/bin/bash
set -eu

# source file path
SRC_PATH=/docker-entrypoint-initdb.d/books.txt.gz

# password for read-only "bookman_web" database role
BOOKMAN_WEB_PASSWORD="$(cat /run/secrets/bookman_web_password)"

# create books database and roles
# (see header of books.txt.gz for details)
zcat "$SRC_PATH" | psql -v ON_ERROR_STOP=1 -v BOOKMAN_WEB_PASSWORD="$BOOKMAN_WEB_PASSWORD" --dbname "$POSTGRES_DB"
