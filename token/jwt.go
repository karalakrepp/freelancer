package token

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/karalakrepp/Golang/freelancer-project/models"
)

const (
	minSecretKeySize = 32
)

// this has secret key for the symmetric algorithm
type jwtMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &jwtMaker{secretKey}, nil
}

func (j *jwtMaker) CreateToken(user *models.User) (string, error) {

	mySigningKey := os.Getenv("SECRET_KEY")
	// Create the Claims
	claims := &jwt.MapClaims{
		"ExpiresAt": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"email":     user.Email,
		"ID":        user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(mySigningKey))

}
func (j *jwtMaker) ValidateJWT(tokenString string) (*jwt.Token, error) {

	secretkey := os.Getenv("SECRET_KEY")
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secretkey), nil
	})

}
