
src = $(wildcard *.go)

bin/mon: test $(src)
	go build -o bin/mon

.PHONY: test
test:
	go vet ./...
	go test -race ./...
