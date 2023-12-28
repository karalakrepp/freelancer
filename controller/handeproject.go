package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/karalakrepp/Golang/freelancer-project/database"
	"github.com/karalakrepp/Golang/freelancer-project/models"
)

var databases = database.PostgresStore{}

func (s *ApiService) handleProject(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.AddProject(w, r)
	}
	if r.Method == "GET" {
		return s.handleGetProject(w, r)
	}
	return WriteJSON(w, 400, "Permision Denied")
}

func (s *ApiService) AddProject(w http.ResponseWriter, r *http.Request) error {
	id := idToken

	var req = new(models.CreateProject)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return WriteJSON(w, 400, map[string]string{
			"message": "Data format is invalid",
		})
		// Bu satırı ekleyin
	}

	// helper for category value
	err := s.store.IsThisCategoryIDExist(req.CategoryID)

	if err != nil {
		if err.Error() == "not exist" {
			return WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "This category id does not exist in the database id"})

		}

		return WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})

	}

	skillDataQuery := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(req.SkillsID)), ","), "[]")

	req.Created_At = time.Now()
	svc, _ := database.NewPostgresStore()
	queryResult, err := svc.DB.Exec("INSERT INTO projects(title, description, skills_id, price, owner_id, category_id, created_at)VALUES($1,$2,$3,$4,$5,$6,$7)", req.Title, req.Description, skillDataQuery, req.Price, id, req.CategoryID, req.Created_At)
	if err != nil {
		return err
	}
	rowID, err := queryResult.RowsAffected()
	if err != nil {
		return err
	}

	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to retrieve the id inserted"})

	}

	// // add project link list
	// for i := 0; i < len(req.Attachment); i++ {
	// 	_, err = svc.DB.Exec("INSERT INTO project_links(project_link, project_id) VALUES(?,?)", req.Attachment[i], rowID)

	// 	if err != nil {
	// 		return WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
	// 			"code":    http.StatusInternalServerError,
	// 			"message": "Server unable to execute query to database"})

	// 	}
	// }
	return WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Successfully Added New Project",
		"id":      rowID,
		"ownerid": id})
}

func (s *ApiService) handleGetProject(w http.ResponseWriter, r *http.Request) error {

	projects, err := s.store.GetAllProject()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, projects)

}

func (s *ApiService) GetProjectByCategoryID(w http.ResponseWriter, r *http.Request) error {

	userIDStr := chi.URLParam(r, "categoryID")

	// Convert userIDStr to an integer
	categoryID, err := strconv.Atoi(userIDStr)
	if err != nil {
		// Handle the error (e.g., invalid integer in the URL parameter)
		http.Error(w, "Invalid userID", http.StatusBadRequest)

	}
	if err != nil {
		return err
	}

	categorys, err := s.store.GetProjectByCategoryID(categoryID)

	if err != nil {
		return WriteJSON(w, 500, map[string]string{
			"message": err.Error(),
		})

	}

	return WriteJSON(w, 200, categorys)
}
