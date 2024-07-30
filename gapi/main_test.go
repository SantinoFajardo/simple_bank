package gapi

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/santinofajardo/simpleBank/db/sqlc"
	"github.com/santinofajardo/simpleBank/token"
	"github.com/santinofajardo/simpleBank/util"
	"github.com/santinofajardo/simpleBank/workers"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor workers.TaskDistributor) *Server {
	config := util.Config{
		TokenSecret:   util.RandomString(32),
		TokenDuration: time.Minute,
	}
	newTestServer, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return newTestServer
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, timeDuration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(username, timeDuration)
	require.NoError(t, err)
	bearerToken := fmt.Sprintf("%s %s", authorizatioBearerKey, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{bearerToken},
	}
	return metadata.NewIncomingContext(context.Background(), md)
}
