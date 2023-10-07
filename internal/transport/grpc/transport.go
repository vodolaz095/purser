package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vodolaz095/purser/internal/service"
	"github.com/vodolaz095/purser/internal/transport/grpc/proto"
	"github.com/vodolaz095/purser/model"
)

func convertModelToDto(secret model.Secret) *proto.Secret {
	meta := make([]*proto.Meta, len(secret.Meta))
	for k := range secret.Meta {
		meta = append(meta, &proto.Meta{
			Key:   k,
			Value: secret.Meta[k],
		})
	}
	return &proto.Secret{
		Id:        secret.ID,
		Body:      secret.Body,
		Meta:      meta,
		CreatedAt: timestamppb.New(secret.CreatedAt),
		ExpiresAt: timestamppb.New(secret.ExpireAt),
	}
}

func convertMetaDTO(meta []*proto.Meta) (ret map[string]string) {
	ret = make(map[string]string, len(meta))
	for k := range meta {
		ret[meta[k].Key] = meta[k].Value
	}
	return
}

type PurserGrpcServer struct {
	proto.UnimplementedPurserServer
	Service service.SecretService
}

func (pgs *PurserGrpcServer) extractJwtSubject(ctx context.Context) (string, error) {
	raw := ctx.Value(TokenSubjectKey)
	switch raw.(type) {
	case string:
		return raw.(string), nil
	default:
		return "", fmt.Errorf("wrong type for %s: %v", TokenSubjectKey, raw)
	}
}

func (pgs *PurserGrpcServer) GetSecretByID(ctx context.Context, request *proto.SecretByIDRequest) (*proto.Secret, error) {
	ctx2, span := pgs.Service.Tracer.Start(ctx, "transport/grpc/GetSecretByID")
	defer span.End()

	subject, err := pgs.extractJwtSubject(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	span.AddEvent("JWT token validated")
	span.SetAttributes(attribute.String("subject", subject))

	secret, err := pgs.Service.FindByID(ctx2, request.GetId())
	if err != nil {
		if errors.Is(err, model.SecretNotFoundError) {
			log.Debug().
				Str("trace_id", span.SpanContext().TraceID().String()).
				Str("secret_id", request.GetId()).
				Str("subject", ctx2.Value(TokenSubjectKey).(string)).
				Msgf("Пользователь %s не нашёл секрет %s",
					ctx.Value(TokenSubjectKey).(string), secret.ID,
				)
			return nil, status.Errorf(codes.NotFound, "secret %s is not found", request.GetId())
		}
		log.Error().Err(err).
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("secret_id", request.GetId()).
			Str("subject", ctx2.Value(TokenSubjectKey).(string)).
			Msgf("Ошибка при поиске секрета %s: %s", request.GetId(), err)
		return nil, err
	}
	log.Info().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("secret_id", request.GetId()).
		Str("subject", ctx2.Value(TokenSubjectKey).(string)).
		Msgf("Пользователь %s нашёл секрет %s",
			ctx.Value(TokenSubjectKey).(string), secret.ID,
		)
	return convertModelToDto(secret), nil
}

func (pgs *PurserGrpcServer) DeleteSecretByID(ctx context.Context, request *proto.SecretByIDRequest) (*proto.Nothing, error) {
	ctx2, span := pgs.Service.Tracer.Start(ctx, "transport/grpc/DeleteSecretByID")
	defer span.End()

	subject, err := pgs.extractJwtSubject(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	span.AddEvent("JWT token validated")
	span.SetAttributes(attribute.String("subject", subject))

	err = pgs.Service.DeleteByID(ctx2, request.GetId())
	if err != nil {
		if errors.Is(err, model.SecretNotFoundError) {
			return nil, status.Errorf(codes.NotFound, "secret %s is not found", request.GetId())
		}
		log.Error().Err(err).
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("secret_id", request.GetId()).
			Str("subject", ctx2.Value(TokenSubjectKey).(string)).
			Msgf("Ошибка при удалении секрета %s : %s", request.GetId(), err)
		return nil, err
	}
	log.Info().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("secret_id", request.GetId()).
		Str("subject", ctx2.Value(TokenSubjectKey).(string)).
		Msgf("Пользователь %s удалил секрет %s",
			ctx.Value(TokenSubjectKey).(string), request.GetId(),
		)
	return nil, nil
}

func (pgs *PurserGrpcServer) CreateSecret(ctx context.Context, request *proto.NewSecretRequest) (*proto.Secret, error) {
	ctx2, span := pgs.Service.Tracer.Start(ctx, "transport/grpc/CreateSecret")
	defer span.End()

	subject, err := pgs.extractJwtSubject(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	span.AddEvent("JWT token validated")
	span.SetAttributes(attribute.String("subject", subject))

	secret, err := pgs.Service.Create(ctx2, request.Body, convertMetaDTO(request.Meta))
	if err != nil {
		log.Error().Err(err).
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("subject", ctx2.Value(TokenSubjectKey).(string)).
			Msgf("Ошибка при создании секрета : %s", err)
		return nil, err
	}
	log.Info().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("secret_id", secret.ID).
		Str("subject", ctx2.Value(TokenSubjectKey).(string)).
		Msgf("Пользователь %s создал секрет %s",
			ctx.Value(TokenSubjectKey).(string), secret.ID,
		)
	return convertModelToDto(secret), nil
}
