package common

import (
	"context"

	"github.com/go-oauth2/oauth2/v4"
	"google.golang.org/grpc/metadata"
)

type tokenKey struct{}

// WithToken returns a copy of parent in which the value associated with key is val.
func WithToken(parent context.Context, token oauth2.TokenInfo) context.Context {
	return context.WithValue(parent, tokenKey{}, token)
}

func WithOutoginRole(ctx context.Context, role string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.Pairs("Role", role)
	} else {
		md = metadata.Join(md, metadata.Pairs("Role", role))
	}

	return metadata.NewOutgoingContext(ctx, md)
}

// FromToken returns the value associated with key from parent.
func FromToken(parent context.Context) (oauth2.TokenInfo, bool) {
	token, ok := parent.Value(tokenKey{}).(oauth2.TokenInfo)
	return token, ok
}

// GetMetadata returns the value associated with key from parent.
func GetMetadata(parent context.Context, key string) (string, bool) {
	md, ok := metadata.FromIncomingContext(parent)
	if !ok {
		return "", false
	}

	vals := md.Get(key)
	if len(vals) == 0 {
		return "", false
	}

	return vals[0], true
}
