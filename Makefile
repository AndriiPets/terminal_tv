build:
	@go build -o bin/yt

run: build
	@./bin/yt

test:
	@go test ./... -v