package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/huskyrobotdog/toolbox-go/id"
	"github.com/huskyrobotdog/toolbox-go/inner"
)

type Config struct {
	Issuer      string `json:"issuer"`
	SigningKey  string `json:"signingKey"`
	ExpiresTime int64  `json:"expiresTime"`
	BufferTime  int64  `json:"bufferTime"`
}

var (
	ErrFormatInvalid = errors.New("token format invalid")
	ErrExpired       = errors.New("token expired")
	ErrNotActiveYet  = errors.New("token not active yet")
	ErrInvalid       = errors.New("token invalid")
)

type API interface {
	NewToken(id id.ID) (string, error)
	NeedFlush(claims *jwt.RegisteredClaims) bool
	FlushToken(claims *jwt.RegisteredClaims) (string, error)
	ParseToken(tokenStr string) (*jwt.RegisteredClaims, error)
	Config() *Config
}

var Instance API

type provider struct {
	config *Config
}

func (p *provider) NewToken(id id.ID) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    p.config.Issuer,
		ExpiresAt: jwt.NewNumericDate(time.Unix(time.Now().Unix()+p.config.ExpiresTime, 0)),
		IssuedAt:  jwt.NewNumericDate(time.Unix(time.Now().Unix(), 0)),
		ID:        id.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if tokenStr, err := token.SignedString([]byte(p.config.SigningKey)); err == nil {
		return tokenStr, nil
	} else {
		return "", err
	}
}

func (p *provider) NeedFlush(claims *jwt.RegisteredClaims) bool {
	return claims.ExpiresAt.Unix()-time.Now().Unix() < p.config.BufferTime
}

func (p *provider) FlushToken(claims *jwt.RegisteredClaims) (string, error) {
	claims.ExpiresAt = jwt.NewNumericDate(time.Unix(time.Now().Unix()+p.config.ExpiresTime, 0))
	claims.IssuedAt = jwt.NewNumericDate(time.Unix(time.Now().Unix(), 0))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if tokenStr, err := token.SignedString([]byte(p.config.SigningKey)); err == nil {
		return tokenStr, nil
	} else {
		return "", err
	}
}

func (p *provider) ParseToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.config.SigningKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrFormatInvalid
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrNotActiveYet
			} else {
				return nil, ErrInvalid
			}
		}
		return nil, ErrInvalid
	}
	if token == nil {
		return nil, ErrInvalid
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if ok {
		return claims, nil
	}
	return nil, ErrInvalid
}

func (p *provider) Config() *Config {
	return p.config
}

func Initialization(config *Config) {
	Instance = &provider{config: config}
	inner.Debug("token initialization complete")
}
