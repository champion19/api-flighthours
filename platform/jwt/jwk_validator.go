package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenExpired     = errors.New("token has expired")
	ErrTokenNotValidYet = errors.New("token is not valid yet")
	ErrInvalidSignature = errors.New("token signature is invalid")
	ErrInvalidIssuer    = errors.New("token issuer is invalid")
	ErrInvalidClaims    = errors.New("token claims are invalid")
	ErrJWKSUnavailable  = errors.New("JWKS endpoint is unavailable")
	ErrTokenMalformed   = errors.New("token is malformed")
)

type JWKSValidator struct {
	jwks           *keyfunc.JWKS
	expectedIssuer string
}
type JWKSConfig struct {
	JWKSURL         string
	Issuer          string
	RefreshInterval time.Duration
}


func NewJWKSValidator(ctx context.Context, config JWKSConfig) (*JWKSValidator, error) {
	if config.JWKSURL == "" {
		return nil, errors.New("JWKS URL cannot be empty")
	}

	if config.RefreshInterval == 0 {
		config.RefreshInterval = time.Hour
	}


	options := keyfunc.Options{
		Ctx: ctx,
		RefreshErrorHandler: func(err error) {

			fmt.Printf("JWKS refresh error: %v\n", err)
		},
		RefreshInterval:   config.RefreshInterval,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	}

	jwks, err := keyfunc.Get(config.JWKSURL, options)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJWKSUnavailable, err)
	}

	return &JWKSValidator{
		jwks:           jwks,
		expectedIssuer: config.Issuer,
	}, nil
}


func (v *JWKSValidator) ValidateToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, v.jwks.Keyfunc)
	if err != nil {
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) {
			switch {
			case validationErr.Errors&jwt.ValidationErrorExpired != 0:
				return nil, ErrTokenExpired
			case validationErr.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, ErrTokenNotValidYet
			case validationErr.Errors&jwt.ValidationErrorSignatureInvalid != 0:
				return nil, ErrInvalidSignature
			case validationErr.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, ErrTokenMalformed
			}
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidClaims, err)
	}

	if !token.Valid {
		return nil, ErrInvalidSignature
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}
	if v.expectedIssuer != "" {
		issuer, ok := claims["iss"].(string)
		if !ok || issuer != v.expectedIssuer {
			return nil, ErrInvalidIssuer
		}
	}

	result := make(map[string]interface{})
	for k, val := range claims {
		result[k] = val
	}

	return result, nil
}

func (v *JWKSValidator) Close() {
	if v.jwks != nil {
		v.jwks.EndBackground()
	}
}
