package service

import (
	"context"
	"errors"

	"github.com/vodolaz095/purser/internal/repository"
	"github.com/vodolaz095/purser/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type SecretService struct {
	Ready  bool
	Tracer trace.Tracer
	Repo   repository.SecretRepo
}

// Ping ensures repository of service is online
func (ss *SecretService) Ping(ctx context.Context) error {
	ctxWithTracing, span := ss.Tracer.Start(ctx, "service.Ping")
	defer span.End()
	err := ss.Repo.Ping(ctxWithTracing)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}
	span.AddEvent("Repository is online!")
	return nil
}

// Create creates new secret
func (ss *SecretService) Create(ctx context.Context, body string, meta map[string]string) (model.Secret, error) {
	ctxWithTracing, span := ss.Tracer.Start(ctx, "service.Create")
	defer span.End()
	span.SetAttributes(attribute.String("body", body))
	for k := range meta {
		span.SetAttributes(attribute.String("meta_"+k, meta[k]))
	}
	secret, err := ss.Repo.Create(ctxWithTracing, body, meta)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return model.Secret{}, err
	}
	span.AddEvent("Secret is created")
	span.SetAttributes(attribute.String("secret_id", secret.ID))
	return secret, err
}

// FindByID finds secret by id
func (ss *SecretService) FindByID(ctx context.Context, id string) (model.Secret, error) {
	ctxWithTracing, span := ss.Tracer.Start(ctx, "service.Find")
	defer span.End()
	span.SetAttributes(attribute.String("secret_id", id))
	span.AddEvent("Searching for secret by id...")
	secret, err := ss.Repo.FindByID(ctxWithTracing, id)
	if err != nil {
		if errors.Is(err, model.SecretNotFoundError) {
			span.AddEvent("Secret not found")
		} else { // unexpected error
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
		return model.Secret{}, err
	}
	span.AddEvent("Secret is found!")
	span.SetAttributes(attribute.String("body", secret.Body))
	for k := range secret.Meta {
		span.SetAttributes(attribute.String("meta_"+k, secret.Meta[k]))
	}
	return secret, nil
}

func (ss *SecretService) DeleteByID(ctx context.Context, id string) error {
	ctxWithTracing, span := ss.Tracer.Start(ctx, "service.Delete")
	defer span.End()
	span.SetAttributes(attribute.String("secret_id", id))
	span.AddEvent("Deleting secret by id...")
	err := ss.Repo.DeleteByID(ctxWithTracing, id)
	if err != nil {
		if errors.Is(err, model.SecretNotFoundError) {
			span.AddEvent("Secret not found")
		} else { // unexpected error
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}
	span.AddEvent("Secret is deleted")
	return err
}

func (ss *SecretService) Prune(ctx context.Context) error {
	ctxWithTracing, span := ss.Tracer.Start(ctx, "service.Prune")
	defer span.End()
	err := ss.Repo.Prune(ctxWithTracing)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	} else {
		span.AddEvent("Secrets pruned")
	}
	return err
}