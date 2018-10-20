package token

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	envHmacSecret string = "HMACSECRET"
	envPrivKey    string = "JWTPRIVATEKEY"
	envPubKey     string = "JWTPUBLICKEY"
	envLifetime   string = "JWTLIFETIME"
)

func NewJwt(user string, attr string) (string, error) {
	// Choose type of JWT, depending on which env var is set
	if os.Getenv(envPrivKey) != "" {
		tok, err := newRsaJwt(user, attr)
		return tok, err
	}
	if os.Getenv(envHmacSecret) != "" {
		tok, err := newHmacJwt(user, attr)
		return tok, err
	}
	return "", errors.New("No env vars set for any of the signing algorithms.")
}

func newHmacJwt(username string, attr string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims(username, attr))
	// Sign and get the complete encoded jwt (header, payload, signature)
	// as a string using the secret
	tokenString, _ := token.SignedString([]byte(os.Getenv(envHmacSecret)))
	return tokenString, nil
}

func newRsaJwt(username string, attr string) (string, error) {
	return "", nil
}

func claims(username string, attr string) jwt.MapClaims {
	// How long should the JWT be valid?
	d := time.Hour // Default value
	// JWT lifetime could be overridden
	if lifetime := os.Getenv(envLifetime); lifetime != "" {
		if parseD, err := time.ParseDuration(lifetime); err == nil {
			d = parseD
		}
	}
	return jwt.MapClaims{
		"iat":   time.Now().UTC().Unix(),
		"exp":   time.Now().UTC().Add(d).Unix(),
		"email": username,
		"ext":   attr,
	}
}

func ValidJwt(jwtstring string) bool {
	// Override jwt lib time with UTC timezone as we use UTC when creating new JWT
	jwt.TimeFunc = time.Now().UTC
	// Parse JWT
	token, err := jwt.Parse(jwtstring, validatorCallback)
	// Was parsing the JWT successful?
	if err != nil {
		return false
	}
	if token.Valid {
		return true
	}
	// Token is not valid
	return false
}

// Returns the secret/pubKey to verifiy the jwt signature
func validatorCallback(t *jwt.Token) (interface{}, error) {
	// If the RSA Public Key Env variable is set, return pub key
	if key := os.Getenv(envPubKey); key != "" {
		// Validate if 'alg' in header is RSA
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf(`RSA signing method expected in `+
				`jwt header but found method: %v`, t.Header["alg"])
		}
		// 'alg' == RSA
		return key, nil
	}
	// If the HMAC Secret Env variable is set, return secret
	if secret := os.Getenv(envHmacSecret); secret != "" {
		// Validate if 'alg' in header is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(`HMAC signing method expected in `+
				`jwt header but found method: %v`, t.Header["alg"])
		}
		// 'alg' == HMAC
		return []byte(secret), nil
	}
	return nil, errors.New("No env vars set for JWT validation.")
}
