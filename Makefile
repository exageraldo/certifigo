install:
	@go get -v ./...

build-cli: install
	@CGO_ENABLED=0 go build -v -ldflags="-s -w" -o bin/certifigo-cli ./cmd/cli/*.go

build-microservice: install
	@CGO_ENABLED=1 go build -v -ldflags="-s -w" -o bin/certifigo-microservice ./cmd/microservice/*.go

build: build-cli build-microservice