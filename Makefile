.Phony: start.server start.client build assets docker.start docker.logs docker.stop docker.client mock tests docker.lint lint

current_dir = $(shell pwd)

tools:
	go install github.com/francoispqt/gojay/gojay@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/vektra/mockery/v2@v2.20.1

start.server: assets build
	.bin/./alo

start.client:
	telnet localhost 8080

build:
	mkdir -p .bin
	go build -o .bin/alo ./cmd/server/main.go
	go build -o .bin/alo-client ./cmd/client/main.go

lint:
	golangci-lint run

assets:
	gojay -s ./pkg/model/ -t User -o ./pkg/model/user.gen.go

tests: mocks
	go test -count=1 ./...

docker.start: docker.build
	docker-compose up -d

docker.stop:
	docker-compose down

docker.build:
	docker-compose build

docker.logs:
	docker-compose logs -f

docker.lint:
	docker run --rm -v $(current_dir):/app -w /app golangci/golangci-lint:v1.23.7 golangci-lint run

docker.client:
	docker run praqma/network-multitool telnet localhost 8080
	

mocks:
	rm -r ./pkg/mocks && mockery --dir=./pkg --all --output=./pkg/mocks