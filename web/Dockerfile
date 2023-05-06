# build stage
FROM docker.io/golang:1.20-alpine AS build
COPY . /src
WORKDIR /src
RUN ["go", "build", "-trimpath", "-ldflags=-s -w"]

# main stage
FROM scratch
COPY --from=build /src/bookman /bookman

# set config env var defaults
ENV BOOKMAN_PASSWORD_PATH="/run/secrets/bookman_web_password" \
    BOOKMAN_DATABASE_DSN="host=db dbname=bookman user=bookman_web" \
    BOOKMAN_HTTP_ADDR=":3000"

CMD ["/bookman"]
