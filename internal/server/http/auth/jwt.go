package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
	"todoNote/internal/model"
)

const (
	issuer = "todo notes service"
)

var _ IAuth = &JwtAuth{}

type JwtToken = string

type Claims struct {
	UserId model.Id
	jwt.StandardClaims
}

type IAuth interface {
	CreateToken(user model.UserInReq) (JwtToken, error)
	ValidateToken(t JwtToken) (model.UserInReq, error)
}


type JwtAuth struct {
	Lifetime time.Duration
	verifyKey *rsa.PublicKey
	signKey *rsa.PrivateKey
}

func NewJwtAuth(lifetime time.Duration, privateKey, publicKey string) (IAuth, error) {
	a := JwtAuth{
		Lifetime: lifetime,
	}
	err := a.initAuthKeys(privateKey, publicKey)
	return &a, err
}

func (auth *JwtAuth) initAuthKeys(privateKey, publicKey string) error {
	pk, _ := base64.StdEncoding.DecodeString(privateKey)
	sk, err := jwt.ParseRSAPrivateKeyFromPEM(pk)
	if err != nil {
		return fmt.Errorf("initKeys read private key: %w", err)
	}
	auth.signKey = sk

	pk, _ = base64.StdEncoding.DecodeString(publicKey)
	vk, err := jwt.ParseRSAPublicKeyFromPEM(pk)
	if err != nil {
		return fmt.Errorf("initKeys read public key: %w", err)
	}
	auth.verifyKey = vk

	return nil
}

func(auth *JwtAuth) CreateToken(user model.UserInReq) (JwtToken, error) {
	c := &Claims{
		user.Id,
		jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
			ExpiresAt: time.Now().Add(auth.Lifetime).Unix(),
			Issuer:    issuer,
		},
	}

	tk := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	t, err := tk.SignedString(auth.signKey)
	if err != nil {
		return "", fmt.Errorf("sign token for user.go %v : %w", user.Id, err)
	}

	return t, err
}

func(auth *JwtAuth) ValidateToken(t JwtToken) (model.UserInReq, error) {
	tk, err := jwt.ParseWithClaims(t, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return auth.verifyKey, nil
	})

	if err != nil {
		return model.UserInReq{}, fmt.Errorf("validate token: %w", err)
	}

	if !tk.Valid {
		return model.UserInReq{}, fmt.Errorf("validate token: not valid token: %w", err)
	}

	c, ok := tk.Claims.(*Claims)
	if !ok {
		return model.UserInReq{}, fmt.Errorf("validate token: invalid token structure: %w", err)
	}

	if c.Issuer != issuer {
		return model.UserInReq{}, fmt.Errorf("validate token: not valid token: %w", err)
	}

	u := model.UserInReq{
		Id: c.UserId,
	}

	return u, nil
}

