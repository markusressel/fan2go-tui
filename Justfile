GO_FLAGS := ""
NAME := "fan2go-tui"
OUTPUT_BIN := "bin/" + NAME
PACKAGE := "github.com/markusressel/" + NAME
GIT_REV := `git rev-parse --short HEAD`
DATE := `date -u +"%Y-%m-%dT%H:%M:%SZ"`
VERSION := "0.4.0"

# Run all tests
test:
    go clean --testcache
    go test -v ./...

# Run all tests with coverage and show summary
coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out

# Run all tests with coverage and open HTML report
coverage-html:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out

# Builds the CLI
build:
    @go build {{GO_FLAGS}} \
        -ldflags "-w -s \
        -X {{NAME}}/cmd/global.Version={{VERSION}} \
        -X {{PACKAGE}}/cmd/global.Version={{VERSION}} \
        -X {{NAME}}/cmd/global.Commit={{GIT_REV}} \
        -X {{PACKAGE}}/cmd/global.Commit={{GIT_REV}} \
        -X {{NAME}}/cmd/global.Date={{DATE}} \
        -X {{PACKAGE}}/cmd/global.Date={{DATE}}" \
        -a -tags netgo -o {{OUTPUT_BIN}} main.go

# Builds the CLI for cross-compilation (ensures CGO is enabled)
build-cross:
    CGO_ENABLED=1 just build

# Build and run the CLI
run: build
    ./{{OUTPUT_BIN}}

# Deploy to custom bin directory
deploy-custom: clean build
    cp ./{{OUTPUT_BIN}} ~/.custom/bin/

# Deploy to /usr/local/bin
deploy: clean build
    sudo cp ./{{OUTPUT_BIN}} /usr/local/bin/
    sudo chmod ug+x /usr/local/bin/{{NAME}}

# Clean build artifacts
clean:
    go clean
    rm -f {{OUTPUT_BIN}}