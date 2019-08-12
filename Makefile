
src = $(wildcard *.go)

bin/autosway: test cmd/autosway/main.go $(src)
	go build -o bin/autosway autosway/cmd/autosway

.PHONY: test
test:
	go vet ./...
	go test -race ./...
