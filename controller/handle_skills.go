package controller

import (
	"net/http"

	"github.com/karalakrepp/Golang/freelancer-project/models"
)

func (s *ApiService) getAllSkills(w http.ResponseWriter, r *http.Request) error {
	skills, err := s.store.GetAllSkills()

	if err != nil {
		return WriteJSON(w, 500, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    []models.UserSkills{}})

	}
	return WriteJSON(w, 200, map[string]interface{}{"code": http.StatusOK,
		"message": "All Skills data have been retrieved",
		"data":    skills})
}
