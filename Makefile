all: run

build: .
	go build -o ./bin/myreddit ./cmd/myreddit/main.go 

run: build
	./bin/myreddit
