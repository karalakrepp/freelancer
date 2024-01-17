package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/karalakrepp/Golang/freelancer-project/database"
	"github.com/karalakrepp/Golang/freelancer-project/token"
)

type ApiError struct {
	Error string `json:"error"`
}

type ctxKeyType string

const (
	ctxKey ctxKeyType = "myContextKey"
)

var idToken string

func WithJWTAuth(handlerFunc http.HandlerFunc, st database.Storage, mk token.Maker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")
		cookie, err := r.Cookie("Authorization")
		if err != nil {
			permissionDenied(w)
			return
		}
		cookieValues := cookie.Value
		values := strings.Split(cookieValues, "|")

		userId, _ := strconv.Atoi(values[1])

		tokenString := values[0]

		token, err := mk.ValidateJWT(tokenString)

		if err != nil {
			w.WriteHeader(400)
			permissionDenied(w)

			return
		}

		if !token.Valid {
			w.WriteHeader(400)
			permissionDenied(w)
			return
		}

		acc, err := st.GetUserByID(userId)
		if err != nil {
			fmt.Println("cant find id")
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if acc.Email != claims["email"] {
			w.WriteHeader(400)
			permissionDenied(w)
			log.Println("İts not your account")
			return
		}

		// Diğer talepleri kontrol et

		// Örnek olarak, kullanıcı adını bir context'e ekleyebilirsiniz
		// ctx := context.WithValue(r.Context(), "username", claims.Username)
		// r = r.WithContext(ctx)
		idstr := strconv.Itoa(userId)
		idToken = idstr
		ctx := context.WithValue(r.Context(), ctxKey, idToken)
		r = r.WithContext(ctx)

		fmt.Println(r.Context().Value(ctxKey))

		handlerFunc(w, r)
	}
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}
