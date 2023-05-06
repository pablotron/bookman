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

Quick setup

    # populate password file
    echo -n 'asdf' > pass.txt

    # set env vars
    export BOOKMAN_PASSWORD_PATH=./pass.txt
    export BOOKMAN_DATABASE_DSN='database=bookman host=worm user=bookman_web'

    # run web server on port :3000
    ./bookman

## Notes

no link: https://stackoverflow.com/questions/1321878/how-to-prevent-favicon-ico-requests
pg fts: https://www.postgresql.org/docs/current/textsearch-controls.html
(TODO: ts_headline())
bs icons:
https://icons.getbootstrap.com/icons/pencil-square/
bulma:
https://bulma.io/documentation/form/general/
