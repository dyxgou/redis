run: build
	@ echo "Building redis client..."
	@ ./bin/server

client: build
	@ echo "Building redis..."
	@ ./bin/client

build:
	@ go build -o ./bin/server ./cmd/server/main.go
	@ go build -o ./bin/client ./cmd/client/main.go

