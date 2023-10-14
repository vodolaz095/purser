export majorVersion=0
export minorVersion=1

export arch=$(shell uname)-$(shell uname -m)
export gittip=$(shell git log --format='%h' -n 1)
export patchVersion=$(shell git log --format='%h' | wc -l)
export ver=$(majorVersion).$(minorVersion).$(patchVersion).$(gittip)-$(arch)

export registry="reg.vodolaz095.ru"
export image="purser"


deps:
	go mod download
	go mod verify
	go mod tidy

oapi:
	which sed # dnf install sed
	which oapi-codegen # go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	cd ./api/http && ./generate.sh

protoc:
	which protoc
	which protoc-gen-go
	which protoc-gen-go-grpc

# https://go.dev/blog/govulncheck
# go install golang.org/x/vuln/cmd/govulncheck@latest
vuln:
	which govulncheck
	govulncheck ./...

start:
	which forego # https://github.com/ddollar/forego
	forego start purser

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

container:
	podman build --format=docker \
		--build-arg=VER=$(ver) \
		--build-arg=SUBVER=$(subver) \
		-t $(registry)/$(image):$(gittip) \
		-t $(registry)/$(image):latest \
		-f Dockerfile .
	podman push $(registry)/$(image):$(gittip)
	podman push $(registry)/$(image):latest

podman:
	/bin/podman run --network=host --expose=3000-3001 --env-file .env $(registry)/$(image):$(gittip)

build: grpc
	go build -o build/purser -ldflags "-X main.Version=$(ver)" main.go
