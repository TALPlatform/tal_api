package auth

import (
	"testing"
	"time"

	"github.com/TALPlatform/tal_api/pkg/random"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(random.RandomString(32))
	require.NoError(t, err)

	username := random.RandomName()
	userId := int32(1)
	userSecurityLevel := int32(2)
	tenantId := int32(0)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	tokenPair, err := maker.CreateTokenPair(username, userId, userSecurityLevel, tenantId, duration, time.Minute*5)
	require.NoError(t, err)
	require.NotEmpty(t, tokenPair.AccessToken)
	require.NotEmpty(t, tokenPair.AccessPayload)

	payload, err := maker.VerifyToken(tokenPair.AccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.Equal(t, userId, payload.UserId)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(random.RandomString(32))
	require.NoError(t, err)

	username := random.RandomName()
	userId := int32(1)
	userSecurityLevel := int32(2)
	tenantId := int32(0)

	// Generate expired token
	tokenPair, err := maker.CreateTokenPair(username, userId, userSecurityLevel, tenantId, -time.Minute, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, tokenPair.AccessToken)

	payload, err := maker.VerifyToken(tokenPair.AccessToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestPasetoRefreshToken(t *testing.T) {
	maker, err := NewPasetoMaker(random.RandomString(32))
	require.NoError(t, err)

	username := random.RandomName()
	userId := int32(1)
	securityLevel := int32(1)
	tenantId := int32(0)

	tokenPair, err := maker.CreateTokenPair(username, userId, securityLevel, tenantId, time.Minute, time.Minute*5)
	require.NoError(t, err)

	refreshPayload, err := maker.VerifyRefreshToken(tokenPair.RefreshToken)
	require.NoError(t, err)
	require.Equal(t, userId, refreshPayload.UserId)
	require.NotZero(t, refreshPayload.ID)
}
