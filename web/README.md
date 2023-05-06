# Bookman

Simple web interface to upload and search through text files.

## Build

    # (optional, only needed if assets have changed)
    make assets

    # (optional) run vet, staticcheck, and golangci-lint
    make check

    # build `bookman` executable
    make all

## Run

Quick setup:

    # populate a password file
    # (note: the contents of this file should contain the password for
    # the database role specified in the DSN)
    echo -n 'some password' > pass.txt

    # set password file path (defaults to
    # `/run/secrets/bookman_web_password` if unspecified)
    export BOOKMAN_PASSWORD_PATH=./pass.txt

    # set postgres database DSN, sans password (defaults to `host=db
    # database=bookman user=bookman_web` if unspecified)
    export BOOKMAN_DATABASE_DSN='host=db database=bookman user=bookman_web'

    # run web server on port :3000
    ./bookman
