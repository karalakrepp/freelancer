package controller

import (
	"encoding/json"
	"net/http"

	"github.com/karalakrepp/Golang/freelancer-project/models"
)

func (s *ApiService) CreateCategory(w http.ResponseWriter, r *http.Request) error {

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

	var req = new(models.GetCategoryByParentIdReq)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	category, err := s.store.GetCategoryByParentId(req.ParentID)

	if err != nil {
		return WriteJSON(w, 500, map[string]string{
			"message": "Server is unable to execute query to the database",
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
