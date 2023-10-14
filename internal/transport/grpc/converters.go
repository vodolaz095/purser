package grpc

import (
	"github.com/vodolaz095/purser/internal/transport/grpc/proto"
	"github.com/vodolaz095/purser/model"
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
