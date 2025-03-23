install:
	@go get -v ./...

build: install
	@CGO_ENABLED=0 go build -v -ldflags="-s -w" -o bin/suacuna-cli ./main.go