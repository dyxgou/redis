run: build
	@ ./bin/server

client: build
	@ ./bin/client

build:
	@ go build -o ./bin/server ./cmd/server/main.go
	@ go build -o ./bin/client ./cmd/client/main.go

