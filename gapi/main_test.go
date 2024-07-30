package gapi

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/santinofajardo/simpleBank/db/sqlc"
	"github.com/santinofajardo/simpleBank/util"
	"github.com/santinofajardo/simpleBank/workers"
	"github.com/stretchr/testify/require"
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
