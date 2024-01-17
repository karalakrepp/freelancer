package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/karalakrepp/Golang/freelancer-project/models"
)

func (s *ApiService) CreateCategory(w http.ResponseWriter, r *http.Request) error {

	id, err := strconv.Atoi(idToken)
	if err != nil {
		permissionDenied(w)

	}
	user, err := s.store.GetUserByID(id)
	if err != nil {
		permissionDenied(w)

	}

	if !user.IsAdmin {
		permissionDenied(w)
	}

	var req = new(models.Category)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	category, err := s.store.CreateCategory(*req)

	if err != nil {
		return WriteJSON(w, 500, map[string]string{
			"message": "Server is unable to execute query to the database",
		})

	}

	return WriteJSON(w, 200, category)
}

// duzelicek yanlis  yazildi

func (s *ApiService) GetCategoryByParentId(w http.ResponseWriter, r *http.Request) error {

	idstr := chi.URLParam(r, "categoryID")

	id, err := strconv.Atoi(idstr)

	if err != nil {
		WriteJSON(w, 500, err)
	}

	category, err := s.store.GetCategoryByParentId(id)

	if err != nil {
		return WriteJSON(w, 500, map[string]string{
			"message": "parent not found",
		})

	}

	return WriteJSON(w, 200, category)
}

func (s *ApiService) GetAllCategories(w http.ResponseWriter, r *http.Request) error {
	data, err := s.store.GetCategoryRow()

	if err != nil {
		return WriteJSON(w, 500, map[string]string{
			"message": "Server is unable to execute query to the database",
		})
	}
	return WriteJSON(w, 200, data)
}
