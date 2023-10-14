package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/purser/internal/service"
	"github.com/vodolaz095/purser/internal/transport/grpc/proto"
	"github.com/vodolaz095/purser/model"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// PurserGrpcServer реализует grpc сервер
type PurserGrpcServer struct {
	proto.UnimplementedPurserServer
	SecretService  *service.SecretService
	CounterService *service.CounterService
}

func (pgs *PurserGrpcServer) extractJwtSubject(ctx context.Context) (string, error) {
	// также тут можно вызывать некое хранилище отозванных токенов, чтобы проверить,
	// что этот ещё работает, а также можно проверить роли и разрешения пользователя
	raw := ctx.Value(TokenSubjectKey)
	switch raw.(type) {
	case string:
		return raw.(string), nil
	default:
		return "", fmt.Errorf("wrong type for %s: %v", TokenSubjectKey, raw)
	}
}

// GetSecretByID загружает секрет по его идентификатору
func (pgs *PurserGrpcServer) GetSecretByID(ctx context.Context, request *proto.SecretByIDRequest) (*proto.Secret, error) {
	ctx2, span := pgs.SecretService.Tracer.Start(ctx, "transport/grpc/GetSecretByID")
	defer span.End()

	subject, err := pgs.extractJwtSubject(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	span.AddEvent("JWT token validated")
	span.SetAttributes(attribute.String("subject", subject))
	pgs.CounterService.Increment(ctx2, "grpc_get_secret_called", 1)
	secret, err := pgs.SecretService.FindByID(ctx2, request.GetId())
	if err != nil {
		if errors.Is(err, model.ErrSecretNotFound) {
			pgs.CounterService.Increment(ctx2, "grpc_get_secret_not_found", 1)
			log.Debug().
				Str("trace_id", span.SpanContext().TraceID().String()).
				Str("secret_id", request.GetId()).
				Str("subject", ctx2.Value(TokenSubjectKey).(string)).
				Msgf("Пользователь %s не нашёл секрет %s",
					ctx.Value(TokenSubjectKey).(string), secret.ID,
				)
			return nil, status.Errorf(codes.NotFound, "secret %s is not found", request.GetId())
		}
		pgs.CounterService.Increment(ctx2, "grpc_get_secret_error", 1)
		log.Error().Err(err).
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("secret_id", request.GetId()).
			Str("subject", ctx2.Value(TokenSubjectKey).(string)).
			Msgf("Ошибка при поиске секрета %s: %s", request.GetId(), err)
		return nil, err
	}
	pgs.CounterService.Increment(ctx2, "grpc_get_secret_success", 1)
	log.Info().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("secret_id", request.GetId()).
		Str("subject", ctx2.Value(TokenSubjectKey).(string)).
		Msgf("Пользователь %s нашёл секрет %s",
			ctx.Value(TokenSubjectKey).(string), secret.ID,
		)
	return convertModelToDto(secret), nil
}

// DeleteSecretByID удаляет секрет по его идентификатору
func (pgs *PurserGrpcServer) DeleteSecretByID(ctx context.Context, request *proto.SecretByIDRequest) (*proto.Nothing, error) {
	ctx2, span := pgs.SecretService.Tracer.Start(ctx, "transport/grpc/DeleteSecretByID")
	defer span.End()

	subject, err := pgs.extractJwtSubject(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	span.AddEvent("JWT token validated")
	span.SetAttributes(attribute.String("subject", subject))
	pgs.CounterService.Increment(ctx2, "grpc_delete_secret_called", 1)

	err = pgs.SecretService.DeleteByID(ctx2, request.GetId())
	if err != nil {
		if errors.Is(err, model.ErrSecretNotFound) {
			pgs.CounterService.Increment(ctx2, "grpc_delete_secret_not_found", 1)
			return nil, status.Errorf(codes.NotFound, "secret %s is not found", request.GetId())
		}
		pgs.CounterService.Increment(ctx2, "grpc_delete_secret_error", 1)
		log.Error().Err(err).
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("secret_id", request.GetId()).
			Str("subject", ctx2.Value(TokenSubjectKey).(string)).
			Msgf("Ошибка при удалении секрета %s : %s", request.GetId(), err)
		return nil, err
	}
	pgs.CounterService.Increment(ctx2, "grpc_delete_secret_success", 1)
	log.Info().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("secret_id", request.GetId()).
		Str("subject", ctx2.Value(TokenSubjectKey).(string)).
		Msgf("Пользователь %s удалил секрет %s",
			ctx.Value(TokenSubjectKey).(string), request.GetId(),
		)
	return nil, nil
}

// CreateSecret создаёт новый секрет и возвращает его идентификатор
func (pgs *PurserGrpcServer) CreateSecret(ctx context.Context, request *proto.NewSecretRequest) (*proto.Secret, error) {
	ctx2, span := pgs.SecretService.Tracer.Start(ctx, "transport/grpc/CreateSecret")
	defer span.End()
	subject, err := pgs.extractJwtSubject(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	span.AddEvent("JWT token validated")
	span.SetAttributes(attribute.String("subject", subject))
	pgs.CounterService.Increment(ctx2, "grpc_create_secret_called", 1)
	meta := convertMetaDTO(request.Meta)
	meta["subject"] = subject
	md, found := metadata.FromIncomingContext(ctx)
	if found {
		meta["User-Agent"] = md.Get("User-Agent")[0]
	}
	secret, err := pgs.SecretService.Create(ctx2, request.Body, meta)
	if err != nil {
		pgs.CounterService.Increment(ctx2, "grpc_create_secret_error", 1)
		log.Error().Err(err).
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("subject", subject).
			Msgf("Ошибка при создании секрета : %s", err)
		return nil, err
	}
	pgs.CounterService.Increment(ctx2, "grpc_create_secret_success", 1)
	log.Info().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("secret_id", secret.ID).
		Str("subject", subject).
		Msgf("Пользователь %s создал секрет %s", subject, secret.ID)
	return convertModelToDto(secret), nil
}
