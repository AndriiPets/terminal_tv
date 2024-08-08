build:
	@go build -o bin/yt

run: build
	@./bin/yt --help

test:
	@go test ./... -v