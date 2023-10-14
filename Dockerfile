FROM golang:1.20.8 as builder

RUN mkdir /purser/

# protoc
# Устанавливаем инструментарий для генерации кода из proto файлов
# см. https://github.com/vodolaz095/go-investAPI/blob/master/Dockerfile
ARG PROTOC_VERSION="3.20.3"
ARG PROTOC_GEN_GO_VERSION="1.27.1"
ARG PROTOC_GEN_GO_GPRC_VERSION="1.2.0"
ARG GRPC_GATEWAY_VERSION="2.15.2"

# устанавливаем unzip и GNU make
RUN apt-get update && apt-get install unzip make -y

# устанавливаем плагины для protoc
ENV GOBIN=/usr/bin/
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VERSION}
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GEN_GO_GPRC_VERSION}
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v${GRPC_GATEWAY_VERSION}
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v${GRPC_GATEWAY_VERSION}

# Загружаем protoc и встроенные протоколы
WORKDIR /tmp
RUN wget "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip"

# Устанавливаем protoc
RUN unzip /tmp/protoc-${PROTOC_VERSION}-linux-x86_64.zip bin/protoc -d /usr/
RUN chown root:root /usr/bin/protoc
RUN chmod 0755 /usr/bin/protoc

# Устанавливаем встроенные протоколы
RUN unzip protoc-${PROTOC_VERSION}-linux-x86_64.zip include/* -d /usr/bin/
RUN chown root:root /usr/bin/include/google/ -R


# deps
WORKDIR /purser
ADD go.mod /purser/go.mod
ADD go.sum /purser/go.sum
RUN go mod download && go mod verify

# build
ADD . /purser/
ARG VER=development
RUN make build


# certs and shared libs
FROM alpine:3.17.3 AS alpine
RUN apk add -U --no-cache \
    ca-certificates \
    libc6-compat \
    curl

# Error loading shared library libresolv.so.2 on Alpine in Go 1.20
# after update to go version 1.20
# https://github.com/gohugoio/hugo/issues/10839#issuecomment-1499463944
# https://www.reddit.com/r/golang/comments/10te58n/error_loading_shared_library_libresolvso2_no_such/
#RUN ln -s /lib/libc.so.6 /usr/lib/libresolv.so.2

EXPOSE 3000
HEALTHCHECK --interval=10s --timeout=3s --retries=3 --start-period=10s \
    CMD curl --fail -v -L http://localhost:3000/healthcheck || exit 1

COPY --from=builder /purser/build/purser /bin/purser

CMD ["/bin/purser"]

