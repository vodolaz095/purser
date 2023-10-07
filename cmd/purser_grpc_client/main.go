package main

import (
	"context"
	"crypto/tls"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/vodolaz095/purser/internal/transport/grpc/proto"
)

// tokenAuth реализует интерфейс https://pkg.go.dev/google.golang.org/grpc/credentials#PerRPCCredentials
type tokenAuth struct {
	Token  string
	Secure bool
}

// GetRequestMetadata вызывается при каждом созданном запросе, чтобы добавить в него JWT токен из tokenAuth
func (t tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": "Bearer " + t.Token,
	}, nil
}

// RequireTransportSecurity вызывается при каждом созданном GPRC запросе, чтобы проверить, является ли соединение зашифрованным
func (t tokenAuth) RequireTransportSecurity() bool {
	return t.Secure
}

func main() {
	var address, token, body, id, del string
	var useTLS bool
	var insecureSkipVerify bool
	var res *proto.Secret
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// logging
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp}
	sink := zerolog.New(output).
		With().Timestamp().Caller().
		Logger().Level(zerolog.DebugLevel)
	log.Logger = sink
	flag.StringVar(&address, "addr", "127.0.0.1:3001", "purser gRPC connection string")
	flag.StringVar(&token, "token", "", "jwt token to use")
	flag.StringVar(&body, "body", "", "secret body, if left empty, STDIN is read")
	flag.StringVar(&id, "id", "", "id of secret, if left empty, new secret is created")
	flag.StringVar(&del, "del", "", "id of secret to be deleted")
	flag.BoolVar(&useTLS, "tls", false, "use tls")
	flag.BoolVar(&insecureSkipVerify, "insecure", false, "allow invalid TLS certificates")
	flag.Parse()

	opts := []grpc.DialOption{
		grpc.WithUserAgent("purser-grpc-cli"),
	}

	if useTLS {
		// включаем шифрование
		if insecureSkipVerify {
			log.Warn().Msgf("Шифруем соединение, но НЕ проверяем сертификат у %s", address)
		} else {
			log.Debug().Msgf("Шифруем соединение и проверяем сертификат у %s", address)
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			ServerName:         strings.Split(address, ":")[0], // нужно для SNI
			InsecureSkipVerify: insecureSkipVerify,             // если у удалённого сервера невалидный сертификат
		})))
	} else {
		log.Warn().Msgf("Внимание, используем соединение без шифрования с %s", address)
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if token != "" {
		// добавляем авторизацию по JWT токену
		log.Debug().Msgf("Используем токен для соединения с %s", address)
		opts = append(opts, grpc.WithPerRPCCredentials(tokenAuth{
			Token:  token,
			Secure: useTLS,
		}))
	} else {
		log.Debug().Msgf("Соединяемся с %s без авторизации.", address)
	}
	conn, err := grpc.DialContext(mainCtx, address, opts...)
	if err != nil {
		log.Fatal().Err(err).
			Msgf("Ошибка соединения с %s: %s", address, err)
	}
	defer conn.Close()

	log.Info().Msgf("Соединение с %s установлено!", address)

	client := proto.NewPurserClient(conn)

	if del != "" {
		_, err = client.DeleteSecretByID(mainCtx, &proto.SecretByIDRequest{Id: del})
		if err != nil {
			log.Error().Err(err).
				Msgf("Ошибка удаления секрета %s : %s", del, err)
		}
		log.Info().
			Msgf("Секрета %s удалён", del)
	}

	if del != "" {
		log.Debug().Msgf("Удаление секрета %s", del)
		_, err = client.DeleteSecretByID(mainCtx, &proto.SecretByIDRequest{Id: del})
		if err != nil {
			log.Error().Err(err).
				Msgf("Ошибка удаления секрета %s : %s", del, err)
			return
		}
		log.Info().
			Msgf("Секрет %s удалён", del)
		return
	}
	if id != "" {
		log.Debug().Msgf("Загружаем секрет %s", del)
		res, err = client.GetSecretByID(mainCtx, &proto.SecretByIDRequest{Id: id})
		if err != nil {
			log.Error().Err(err).
				Msgf("Ошибка получения секрета %s : %s", del, err)
			return
		}
		log.Info().Msgf("Секрет %s получен: %s", id, res.String())
		return
	}

	res, err = client.CreateSecret(mainCtx, &proto.NewSecretRequest{
		Body: body,
		Meta: []*proto.Meta{
			{Key: "User-Agent", Value: "purser-grpc-cli"},
		},
	})
	if err != nil {
		log.Error().Err(err).
			Msgf("Ошибка создания секрета %s", err)
		return
	}
	log.Info().Msgf("Секрет %s создан: %s", id, res.String())
}
