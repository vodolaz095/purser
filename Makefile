deps:
	go mod download
	go mod verify
	go mod tidy


protoc:
	which protoc
	which protoc-gen-go
	which protoc-gen-go-grpc


start:
	go run main.go


grpc: protoc
	protoc \
		--proto_path=api/grpc \
		--proto_path=api/grpc/google/protobuf \
		--go_out=./ --go_opt=paths=import \
		--go-grpc_out=./ --go-grpc_opt=paths=import \
        api/grpc/*.proto


build: grpc
	go build -o build/purser main.go
