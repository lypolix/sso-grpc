package tests

import (
	"sso-auth/tests/suite"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ssov1 "github.com/lypolix/protos_sso/gen/go/sso"
)

const(
	emptyAppID = 0
	appID = 1
	appSecret = "test-secret"

	passDefaultLen = 10
)

func TestAuthRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email: email,
		Password: pass,
	
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())
	
	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email: email,
		Password: pass,
		AppId: appID,
	})

	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, float64(appID), claims["app_id"].(float64))

	const deltaSeconds = 1
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestAuthRegisterLogin_DuplicatedRegisration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email: email,
		Password: pass,
	
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())
	
	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email: email,
		Password: pass,
	
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user with this email already exists")
}
func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}