package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/karalakrepp/Golang/freelancer-project/database"
	"github.com/karalakrepp/Golang/freelancer-project/models"
	"github.com/karalakrepp/Golang/freelancer-project/util"
)

func (s *ApiService) Register(w http.ResponseWriter, r *http.Request) error {

	var req = new(models.CreateUserRequest)

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Fatal(err)
	}
	password, err := util.HashPassword(req.Password)
	if err != nil {
		log.Fatal(err)
	}

	user := database.CreateUserParams{
		Username:       req.Username,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		HashedPassword: password,
		Location:       req.Location,
		PhoneNumber:    req.PhoneNumber,
	}

	res, err := s.store.CreateAccount(user)

	if err != nil {
		log.Fatal(err)
	}

	return WriteJSON(w, 200, res)
}

func (s *ApiService) Login(w http.ResponseWriter, r *http.Request) error {

	var req = new(models.LoginUserRequest)
	fmt.Println("dsadasasdasd")

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		fmt.Println("hata")
		log.Fatal(err)
	}
	user, err := s.store.GetUserByEmail(req.Email)
	if err != nil {
		log.Fatal("email does not match")
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		log.Fatal("pass does not match")
	}

	token, err := s.maker.CreateToken(user)
	if err != nil {
		return WriteJSON(w, 500, err)
	}

	cookieValue := fmt.Sprintf("%s|%d|%s", token, user.ID, user.Email)

	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    cookieValue,
		Path:     "",
		Domain:   "",
		MaxAge:   3600 * 24 * 30,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	fmt.Println("login kullandi,token", token)
	return json.NewEncoder(w).Encode(token)

}

func (s *ApiService) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.GetUserProfile(w, r)

	}
	if r.Method == "POST" {
		return s.CreateUserProfile(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)

}
