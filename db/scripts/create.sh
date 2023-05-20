#!/bin/bash
#
# First this script creates the following:
#
# * `bookman_sys` role: database owner
# * `bookman_web` role: web interface role
# * `bookman` database
# * `bookman` schema in `bookman` database
# * `books` table in `bookman` schema of the `bookman` database
#
# Then this script pipes `books.txt.gz` into psql.  `books.txt.gz` does
# the following:
#
# * Populate the `books` table with several books from Project
#   Gutenberg (https://gutenberg.org/)
# * Create a GIN index on the `books(ts_vec)` column.
# * `VACUUM ANALYZE` the `books` table.
#

# set sane behavior
set -euo pipefile

# source file path
SRC_PATH=/docker-entrypoint-initdb.d/books.txt.gz

# password for read-only "bookman_web" database role
BOOKMAN_WEB_PASSWORD="$(cat /run/secrets/bookman_web_password)"

# create roles, db, schema, and table
echo "
  -- create db owner
  CREATE ROLE bookman_sys;
  COMMENT ON ROLE bookman_sys IS 'Bookman database owner';

  CREATE ROLE bookman_web WITH LOGIN PASSWORD :'BOOKMAN_WEB_PASSWORD';
  COMMENT ON ROLE bookman_web IS 'Bookman web user';

  CREATE DATABASE bookman WITH OWNER bookman_sys;
  COMMENT ON DATABASE bookman IS 'Bookman database';
  GRANT CONNECT ON DATABASE bookman TO bookman_web;

  -- connect to bookman database, set role
  \\c bookman
  SET ROLE bookman_sys;

  -- create bookman schema, set privileges
  CREATE SCHEMA bookman;
  COMMENT ON SCHEMA bookman IS 'bookman tables';
  GRANT USAGE ON SCHEMA bookman TO bookman_web;

  -- set search path
  SET search_path = bookman;

  -- create books table
  CREATE TABLE books (
    -- book ID
    id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,

    -- book name
    name TEXT UNIQUE NOT NULL CHECK (LENGTH(name) > 0),

    -- book author
    author TEXT NOT NULL CHECK (LENGTH(author) > 0),

    -- book content
    body TEXT NOT NULL,

    -- fts vector
    ts_vec tsvector GENERATED ALWAYS AS (to_tsvector('english',
      COALESCE(name, '') || ' ' ||
      COALESCE(author, '') || ' ' ||
      COALESCE(body, ''))
    ) STORED
  );

  -- document table and columns
  COMMENT ON TABLE books IS 'Books';
  COMMENT ON COLUMN books.id IS 'Book ID';
  COMMENT ON COLUMN books.name IS 'Book title';
  COMMENT ON COLUMN books.author IS 'Author name';
  COMMENT ON COLUMN books.body IS 'Book contents';
  COMMENT ON COLUMN books.ts_vec IS 'Book FTS vector';

  -- set owner and privileges
  ALTER TABLE books ONWER TO bookman_sys;
  GRANT SELECT, INSERT, UPDATE, DELETE ON books TO bookman_web;
" | psql -v ON_ERROR_STOP=1 -v BOOKMAN_WEB_PASSWORD="$BOOKMAN_WEB_PASSWORD" --dbname "$POSTGRES_DB"

# populate books table, create index, and vacuum table
zcat "$SRC_PATH" | psql -v ON_ERROR_STOP=1 --dbname bookman
