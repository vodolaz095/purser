package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/vodolaz095/purser/api/openapi"
)

func main() {
	var address, token, body, id, del string
	var resp *http.Response
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// logging
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp}
	sink := zerolog.New(output).
		With().Timestamp().Caller().
		Logger().Level(zerolog.DebugLevel)
	log.Logger = sink
	flag.StringVar(&address, "addr", "http://127.0.0.1:3000", "purser http connection string")
	flag.StringVar(&token, "token", "", "jwt token to use")
	flag.StringVar(&body, "body", "", "secret body")
	flag.StringVar(&id, "id", "", "id of secret, if left empty, new secret is created")
	flag.StringVar(&del, "del", "", "id of secret to be deleted")
	flag.Parse()

	client, err := openapi.New(address, token)
	if err != nil {
		log.Fatal().Err(err).Msgf("Ошибка соединения с API через %s: %s", address, err)
	}

	if body != "" {
		resp, err = client.PostApiV1Secret(mainCtx, openapi.PostApiV1SecretJSONRequestBody{
			Body: &body,
		})
		if err != nil {
			log.Fatal().Err(err).Msgf("Ошибка создания секрета: %s", err)
		}
		if resp.StatusCode != http.StatusCreated {
			log.Fatal().Msgf("Неожиданный статус ответа %s", resp.Status)
		}

		log.Info().Msgf("%v", resp.Header)
		id = strings.TrimPrefix(resp.Header.Get("Location"), "/api/v1/secret/")
		log.Info().Msgf("Секрет %s создан", id)
	}
	if id != "" {
		resp, err = client.GetApiV1SecretId(mainCtx, id)
		if err != nil {
			log.Fatal().Err(err).Msgf("Ошибка получения секрета %s: %s", id, err)
		}
		secret, err := openapi.ParseGetApiV1SecretIdResponse(resp)
		if err != nil {
			log.Fatal().Err(err).Msgf("Ошибка получения секрета %s: %s", id, err)
		}
		if secret.StatusCode() != http.StatusOK {
			log.Fatal().Msgf("Неожиданный статус ответа %s", secret.Status())
		}
		log.Info().
			Str("body", *secret.JSON200.Body).
			Str("created_at", *secret.JSON200.CreatedAt).
			Str("expires_at", *secret.JSON200.ExpireAt).
			Str("expires_at", fmt.Sprint(*secret.JSON200.Fields)).
			Msgf("Секрет %s получен", id)
	}
	if del != "" {
		resp, err = client.DeleteApiV1SecretId(mainCtx, id)
		if err != nil {
			log.Fatal().Err(err).Msgf("Ошибка удаления секрета %s: %s", id, err)
		}
		if resp.StatusCode != http.StatusNoContent {
			log.Fatal().Msgf("Неожиданный статус ответа %s", resp.Status)
		}
		log.Info().Msgf("Секрет %s удалён", id)
	}
}
