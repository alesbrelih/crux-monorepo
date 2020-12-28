package my_jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type TokenPair struct {
	Refresh string
	Access  string
}

var InvalidTokenError = errors.New("Invalid jwt token")

func NewJwtService(secret string, accessExp, refreshExp time.Duration) JwtService {
	return &jwtService{
		secret:     []byte(secret),
		accessExp:  accessExp,
		refreshExp: refreshExp,
	}
}

type JwtService interface {
	GenerateJwtPair(id string) (*TokenPair, error)
	RefreshToken(refreshToken string) (*TokenPair, error)
	GetClaims(token string) (*jwt.StandardClaims, error)
}

type jwtService struct {
	secret     []byte
	accessExp  time.Duration
	refreshExp time.Duration
}

// Generates access and refresh token
// id is id of user which will be included as subject inside jwt body
func (a *jwtService) GenerateJwtPair(id string) (*TokenPair, error) {
	accessClaims := &jwt.StandardClaims{
		Subject:   id,
		ExpiresAt: time.Now().Add(a.accessExp).Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := access.SignedString(a.secret)
	if err != nil {
		return nil, errors.Wrap(err, "Error signing access token")
	}

	refreshClaims := &jwt.StandardClaims{
		Subject:   id,
		ExpiresAt: time.Now().Add(a.refreshExp).Unix(),
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refresh.SignedString(a.secret)
	if err != nil {
		return nil, errors.Wrap(err, "Error signing refresh token")
	}

	tokenPair := &TokenPair{
		Refresh: refreshToken,
		Access:  accessToken,
	}
	return tokenPair, nil
}

// refreshes token in parameter if valid and returns access and refresh token
func (a *jwtService) RefreshToken(refreshToken string) (*TokenPair, error) {

	claims, err := a.GetClaims(refreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting claims in refresh token")
	}
	return a.GenerateJwtPair(claims.Subject)
}

// gets claims from jwt token
func (a *jwtService) GetClaims(jwtString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(jwtString, &jwt.StandardClaims{}, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Wrap(InvalidTokenError, "Missmatch signing method")
		}
		return a.secret, nil
	})

	if err != nil {
		return nil, errors.Wrap(InvalidTokenError, "Invalid token")
	}

	if !token.Valid {
		return nil, errors.Wrap(InvalidTokenError, "Token is not valid")
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, errors.Wrap(InvalidTokenError, "Claims cant be parsed")
	}
	return claims, nil
}
