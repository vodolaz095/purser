FROM golang:1.20.8 as builder

RUN mkdir /purser/
WORKDIR /purser

# deps
ADD go.mod /purser/go.mod
ADD go.sum /purser/go.sum
RUN go mod download && go mod verify

# build
ADD . /purser/
ARG VER=development
ARG SUBVER=development
RUN go build -ldflags "-X main.Subversion=$SUBVER -X main.Version=$VER" -o build/purser main.go


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

