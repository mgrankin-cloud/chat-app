package tests

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	ssov1 "github.com/mgrankin-cloud/messenger/contract/gen/go/sso"
	"github.com/mgrankin-cloud/messenger/internal/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID = 0
	appID     = 1
	appSecret = "test-secret"

	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()
	username := gofakeit.Username()
	phone := gofakeit.Phone()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
		Username: username,
		Phone: phone,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Authenticate(ctx, &ssov1.AuthenticateRequest{
		Identifier: email,
		Password:   pass,
	})
	require.NoError(t, err)

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, username, claims["username"].(string))
	assert.Equal(t, phone, claims["phone"].(string))
	assert.Equal(t, appID, int64(claims["app_id"].(float64)))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	username := gofakeit.Username()
	pass := randomFakePassword()
	phone := gofakeit.Phone()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Username: username,
		Password: pass,
		Phone:    phone,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Username: username,
		Password: pass,
		Phone:    phone,
	})
	require.Error(t, err)
	require.Empty(t, respReg.GetUserId())
	assert.Error(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		username    string
		password    string
		phone       string
		expectedErr string
	}{
		{
			name:        "Register with empty password",
			email:       gofakeit.Email(),
			username:    gofakeit.Username(),
			password:    "",
			phone:       gofakeit.Phone(),
			expectedErr: "password is required",
		},
		{
			name:        "Register with empty email",
			email:       "",
			username:    gofakeit.Username(),
			password:    randomFakePassword(),
			phone:       gofakeit.Phone(),
			expectedErr: "email is required",
		},
		{
			name:        "Register with empty username",
			email:       gofakeit.Email(),
			username:    "",
			password:    randomFakePassword(),
			phone:       gofakeit.Phone(),
			expectedErr: "username is required",
		},
		{
			name:        "Register with empty phone",
			email:       gofakeit.Email(),
			username:    gofakeit.Username(),
			password:    randomFakePassword(),
			phone:       "",
			expectedErr: "phone is required",
		},
		{
			name:        "Register with empty frames",
			email:       "",
			password:    "",
			username:    "",
			phone:       "",
			expectedErr: "email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Username: tt.username,
				Password: tt.password,
				Phone:    tt.phone,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		identifier  string
		password    string
		appID       int32
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			identifier:  gofakeit.Email(),
			password:    "",
			appID:       appID,
			expectedErr: "password is required",
		},
		{
			name:        "Login with Empty Identifier",
			identifier:  "",
			password:    randomFakePassword(),
			appID:       appID,
			expectedErr: "identifier is required",
		},
		{
			name:        "Login with Both Empty Identifier and Password",
			identifier:  "",
			password:    "",
			appID:       appID,
			expectedErr: "identifier is required",
		},
		{
			name:        "Login with Non-Matching Password",
			identifier:  gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       appID,
			expectedErr: "invalid identifier or password",
		},
		{
			name:        "Login without AppID",
			identifier:  gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       emptyAppID,
			expectedErr: "app_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    gofakeit.Email(),
				Username: gofakeit.Username(),
				Password: randomFakePassword(),
				Phone:    gofakeit.Phone(),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Authenticate(ctx, &ssov1.AuthenticateRequest{
				Identifier: tt.identifier,
				Password:   tt.password,
				AppId:      tt.appID,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
