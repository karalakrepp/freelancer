package token

import (
	"github.com/golang-jwt/jwt"
	"github.com/karalakrepp/Golang/freelancer-project/models"
)

type Maker interface {
	CreateToken(*models.User) (string, error)
	ValidateJWT(string) (*jwt.Token, error)
}
