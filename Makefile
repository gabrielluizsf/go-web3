build: 
	go build -o ./bin/goweb3

run: build
	./bin/goweb3
test:
	go test -v ./...