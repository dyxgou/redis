run: build
	@ ./bin/server

build:
	@ go build -o ./bin/server ./cmd/server/main.go
	@ go build -o ./bin/client ./cmd/client/main.go

