package gapi

import (
	"context"
	"fmt"
	"strings"

	. "github.com/santinofajardo/simpleBank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader   = "authorization"
	authorizatioBearerKey = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("cannot get metadata")
	}
	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("empty authorization was received")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid header authorization header")
	}

	if strings.ToLower(fields[0]) != authorizatioBearerKey {
		return nil, fmt.Errorf("missing bearer word")
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %s", err)
	}

	return payload, nil
}
