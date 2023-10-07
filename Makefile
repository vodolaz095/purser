deps:
	go mod download
	go mod verify
	go mod tidy


protoc:
	which protoc
	which protoc-gen-go
	which protoc-gen-go-grpc

# https://go.dev/blog/govulncheck
# install it by go install golang.org/x/vuln/cmd/govulncheck@latest
vuln:
	govulncheck ./...

start:
	go run main.go

cli_curl:
	./cmd/curl/create_secret.sh

cli_grpc:
	./cmd/purser_grpc_client/example.sh

grpc: protoc
	protoc \
		--proto_path=api/grpc \
		--proto_path=api/grpc/google/protobuf \
		--go_out=./ --go_opt=paths=import \
		--go-grpc_out=./ --go-grpc_opt=paths=import \
        api/grpc/*.proto

build: grpc
	go build -o build/purser main.go
