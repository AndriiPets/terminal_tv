build:
	@go build -o bin/yt

run: build
	@./bin/yt search decino

test:
	@go test ./... -v