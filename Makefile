run:
	go build -o ./bin/goweb3
	./bin/goweb3
test:
	go test -v ./...