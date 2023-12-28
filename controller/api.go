package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/karalakrepp/Golang/freelancer-project/database"
	"github.com/karalakrepp/Golang/freelancer-project/token"
)

const (
	RemoteAddress string = ":9000"
)

type ApiService struct {
	listenAddr string
	store      database.Storage
	maker      token.Maker
}

func NewApiService(listenAddr string, storer database.Storage, maker token.Maker) *ApiService {
	return &ApiService{
		listenAddr: listenAddr,
		store:      storer,
		maker:      maker,
	}
}

type ApiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func (s *ApiService) authRouter(router chi.Router) {
	router.Post("/api/v1/signup", makeHTTPHandleFunc(s.Register))
	router.Post("/login", makeHTTPHandleFunc(s.Login))
}

func (s *ApiService) userRouter(router chi.Router) {

	router.HandleFunc("/api/v1/profile/{id}", WithJWTAuth(makeHTTPHandleFunc(s.handleAccount), s.store, s.maker))
	router.HandleFunc("/api/v1/project", WithJWTAuth(makeHTTPHandleFunc(s.handleProject), s.store, s.maker))
	router.HandleFunc("/api/v1/project/{categoryID}", WithJWTAuth(makeHTTPHandleFunc(s.GetProjectByCategoryID), s.store, s.maker))

}

func (s *ApiService) Routes() {
	router := chi.NewRouter()

	// Middleware'ları kullanarak rotaları tanımla
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "x-jwt-token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	//auth
	s.authRouter(router)

	//user
	s.userRouter(router)

	// Genel rotalar

	// Admin rotaları sadece yetkilendirilmiş kullanıcılara açık olmalı
	router.Group(func(r chi.Router) {

		// Admin rotaları
		r.Post("/api/v1/admin/category", makeHTTPHandleFunc(s.CreateCategory))
	})

	// Tüm kullanıcılara açık kategori rotaları

	fmt.Println("API listening on port", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
