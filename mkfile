all: build

build:V:
	go build -o exe/ ./cmd/...

clean:
	rm -rf exe/*

