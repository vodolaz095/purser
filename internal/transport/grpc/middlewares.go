package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/rs/zerolog/log"

	"github.com/vodolaz095/purser/pkg"
)

const TokenKey = "jwt_token"
const TokenSubjectKey = "jwt_token_subject"

type ValidateJWTInterceptor struct {
	HmacSecret string
}

// ServerInterceptor работает с унарными запросами
func (ji *ValidateJWTInterceptor) ServerInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "no meta found")
	}
	authHeader, ok := md["authorization"]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authorization header not found")
	}
	if len(authHeader) != 1 {
		return nil, status.Error(codes.InvalidArgument, "multiple authorization headers found")
	}
	token := authHeader[0]
	if !strings.HasPrefix(token, "Bearer ") { // https://www.rfc-editor.org/rfc/rfc6750
		return nil, status.Error(codes.InvalidArgument, "wrong authorization strategy")
	}
	token = strings.TrimPrefix(token, "Bearer ")
	subject, err := pkg.ValidateJwtAndExtractSubject(token, ji.HmacSecret)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	log.Debug().Msgf("JWT token subject = %s", subject)
	ctxWithToken := context.WithValue(ctx, TokenSubjectKey, subject)
	return handler(ctxWithToken, req)
}
