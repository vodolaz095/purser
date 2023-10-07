package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/vodolaz095/purser/internal/service"
	"github.com/vodolaz095/purser/internal/transport/grpc/proto"
	"github.com/vodolaz095/purser/model"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			return nil, status.Errorf(codes.NotFound, "secret %s is not found", request.GetId())
		}
		return nil, err
	}
	return convertModelToDto(secret), nil
}

func (pgs *PurserGrpcServer) DeleteSecretByID(ctx context.Context, request *proto.SecretByIDRequest) (*proto.Nothing, error) {
	ctx2, span := pgs.Service.Tracer.Start(ctx, "transport/grpc/GetSecretByID")
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
		return nil, err
	}
	return nil, nil
}

func (pgs *PurserGrpcServer) CreateSecret(ctx context.Context, request *proto.NewSecretRequest) (*proto.Secret, error) {
	ctx2, span := pgs.Service.Tracer.Start(ctx, "transport/grpc/GetSecretByID")
	defer span.End()

	subject, err := pgs.extractJwtSubject(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	span.AddEvent("JWT token validated")
	span.SetAttributes(attribute.String("subject", subject))

	secret, err := pgs.Service.Create(ctx2, request.Body, convertMetaDTO(request.Meta))
	if err != nil {
		return nil, err
	}
	return convertModelToDto(secret), nil
}
