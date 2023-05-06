
.PHONY=all check clean

all:
	go build -trimpath -ldflags='-s -w'

clean:
	go clean

# needs to be run before `make all` if any of the following have
# changed:
#
# * web/public/script.js 
# * web/assets/bookman.sass
assets:
	ruby web/assets/build.rb

test:
	go test ./...

check: vet staticcheck lint

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

lint:
	golangci-lint run ./...
