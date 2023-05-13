
.PHONY=all test check clean

# build binary
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

# run unit tests
test:
	go test ./...

# run static analysis (vet, staticcheck, lint) and vulnerability scan
check: vet staticcheck lint vulncheck

vet:
	go vet ./...

# go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck:
	staticcheck ./...

# go install golang.org/x/vuln/cmd/govulncheck@latest
vulncheck:
	govulncheck ./...

# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
lint:
	golangci-lint run ./...
