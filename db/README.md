# Bookman DB

Container image with postgres database containing `bookman` database.

The `bookman` database is owned by the `bookman_sys` database role with
a `bookman` schema containing a `books` table.

The `books` table is pre-populated with a selection of books from
Project Gutenberg.

The setup script creates a `bookman_web` role with limited access to
the database.

The image must be run with the following secrets:

* `bookman_postgres_password`: password of `postgres` user.
* `bookman_web_password`: password of `bookman_web` user.

To build the image and run an instance:

    # build image
    podman build -t bookman-db .

    # set secrets
    echo 'some secret password' | podman secret create bookman_postgres_password -
    echo 'another secret password' | podman secret create bookman_web_password -

    # run db container with the following options:
    #
    # 1. expose postgres port
    # 2. mount db_data volume as /data and save postgres data to /data/pgdata
    #    (PGDATA must be a subdirectory of a mounted volume or you will
    #    get permission errors)
    # 3. expose pair of database role password secrets and pass path to
    #    secret file to postgres.
    podman run -d -p 5432:5432 \
      -v db_data:/data \
      -e PGDATA=/data/pgdata \
      --secret bookman_postgres_password \
      --secret bookman_web_password \
      -e POSTGRES_PASSWORD_FILE=/run/secrets/bookman_postgres_password \
      bookman-db

Here's an example query which searches the `books` table using the FTS
index:

    -- find books matching phrase 'evil monster', sorted by relevance
    SELECT id,
           name,
           ts_rank_cd(ts_vec, websearch_to_tsquery('english', 'evil monster')) AS rank
      FROM bookman.books
     WHERE websearch_to_tsquery('english', 'evil monster') @@ ts_vec
     ORDER BY rank DESC;

Additional queries are available in `web/model/sql/`.
