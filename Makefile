build:
	@go build -o bin/goWEB cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/goWEB
