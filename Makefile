build:
	@go build -o bin/website cmd/main.go

run: build
	@./bin/website