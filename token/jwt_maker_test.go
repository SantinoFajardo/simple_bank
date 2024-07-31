package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/santinofajardo/simpleBank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	newJwtMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	userRole := util.DepositorRole
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := newJwtMaker.CreateToken(username, userRole, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = newJwtMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotZero(t, payload.ID)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.Equal(t, username, payload.UserName)
}

func TestExpirtedJWTToken(t *testing.T) {
	newJwtMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := -time.Minute

	token, payload, err := newJwtMaker.CreateToken(username, util.DepositorRole, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = newJwtMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpirationMessage.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTToken(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), util.DepositorRole, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorInvalidToken.Error())
	require.Nil(t, payload)

}

func TestInvalidLengthJWTToken(t *testing.T) {
	newPasetoMaker, err := NewJWTMaker(util.RandomString(31))
	require.Error(t, err)
	require.EqualError(t, err, "invalid key size")
	require.Empty(t, newPasetoMaker)
}
