all: build

run-air:V:
	air -c airfile

build:V:
	go build -o exe/ ./cmd/...

clean:
	rm -rf exe/*

