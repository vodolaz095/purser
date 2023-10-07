deps:
	go mod download
	go mod verify
	go mod tidy

start:
	go run main.go


grpc:
	protoc \
		--proto_path=api/grpc \
		--proto_path=api/src/google/protobuf \
		--go_out=address_book --go_opt=paths=import \
		--go-grpc_out=address_book --go-grpc_opt=paths=import \
        address_book/src/*.proto


build:
	go build -o build/purser main.go
