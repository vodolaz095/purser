#!/usr/bin/env bash

set -e

which oapi-codegen
oapi-codegen -version
echo "Чистим старые файлы"
rm -f client.go
rm -f types.go
echo "Генерируем типы..."
oapi-codegen --generate="types"  --package="purser_client" -o ./types.go  purser.yaml
echo "Генерируем клиент..."
oapi-codegen --generate="client" --package="purser_client" -o ./client.go purser.yaml
echo "Чиним возможные паники с закрытием пустого тела ответа"
sed -i "s/defer func() { _ = rsp.Body.Close() }()/if rsp.Body != nil { defer rsp.Body.Close()}/g" ./client.go
echo "Форматируем код"
gofmt -s -w client.go
gofmt -s -w types.go
echo "Кодогенерация окончена!"
