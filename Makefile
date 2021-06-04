all: run

build: .
	go run ./cmd/myreddit/main.go -o ./bin/myreddit

run: build
	go run ./bin/myreddit
