package pkg

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Platform string `json:"platform"`
	RoleID   string `json:"role_id"`
	jwt.RegisteredClaims
	Login_time      string `json:"login_time"`
	Login_time_unix int64  `json:"login_time_unix"`
}

func NewToken(id int, username, platform, role_id, loginTime string, loginTimeUnix int64) *claims {
	return &claims{
		ID:              id,
		Username:        username,
		Platform:        platform,
		Login_time:      loginTime,
		RoleID:          role_id,
		Login_time_unix: loginTimeUnix,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "CIS Mobile Backend",
		},
	}
}

func (c *claims) Generate() (string, error) {
	secrets := os.Getenv("JWTSECRET")
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return tokens.SignedString([]byte(secrets))
}

func VerifyToken(token string) (*claims, error) {
	secrets := os.Getenv("JWTSECRET")
	data, err := jwt.ParseWithClaims(token, &claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secrets), nil
	})

	if err != nil {
		return nil, err
	}

	claimData := data.Claims.(*claims)
	return claimData, nil
}
