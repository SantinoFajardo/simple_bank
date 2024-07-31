package token

import (
	"testing"
	"time"

	"github.com/o1egl/paseto"
	"github.com/santinofajardo/simpleBank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	newPasetoMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, newPasetoMaker)

	username := util.RandomOwner()
	userRole := util.DepositorRole
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := newPasetoMaker.CreateToken(username, userRole, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = newPasetoMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, payload.Role, userRole)
	require.Equal(t, username, payload.UserName)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	newPasetoMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, newPasetoMaker)

	username := util.RandomOwner()
	userRole := util.DepositorRole
	duration := -time.Minute

	token, payload, err := newPasetoMaker.CreateToken(username, userRole, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)
	require.Equal(t, payload.Role, userRole)

	payload, err = newPasetoMaker.VerifyToken(token)
	require.Error(t, err)
	require.Equal(t, payload.Role, userRole)
	require.EqualError(t, err, ErrorExpirationMessage.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {
	userRole := util.DepositorRole
	payload, err := NewPayload(util.RandomOwner(), userRole, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, payload.Role, userRole)

	token, err := paseto.NewV2().Encrypt([]byte(util.RandomString(32)), payload, nil)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.Equal(t, payload.Role, userRole)
	require.EqualError(t, err, ErrorInvalidToken.Error())
	require.Nil(t, payload)
}

func TestInvalidLengthPasetoToken(t *testing.T) {
	newPasetoMaker, err := NewPasetoMaker(util.RandomString(31))
	require.Error(t, err)
	require.EqualError(t, err, "invalid key size")
	require.Empty(t, newPasetoMaker)
}
