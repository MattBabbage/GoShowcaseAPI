build: 
	@go build -o bin/GoShowcaseAPI

run: build
	@./bin/GoShowcaseAPI

test: 
	@go test -v ./...