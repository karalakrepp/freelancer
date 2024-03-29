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
	"github.com/karalakrepp/Golang/freelancer-project/helper"
	"github.com/karalakrepp/Golang/freelancer-project/models"
)

var svc, _ = database.NewPostgresStore()

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

	categorys, err := s.store.GetProjectByCategoryID(categoryID)

	if err != nil {
		return WriteJSON(w, 500, map[string]string{
			"message": err.Error(),
		})

	}

	return WriteJSON(w, 200, categorys)
}

func (s *ApiService) GetProjectsByOwnerID(w http.ResponseWriter, r *http.Request) error {
	owner_id_str := chi.URLParam(r, "ownerID")

	// Convert userIDStr to an integer
	owner_id, err := strconv.Atoi(owner_id_str)
	if err != nil {
		// Handle the error (e.g., invalid integer in the URL parameter)
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return err // Return the error to stop further execution
	}

	categorys, err := s.store.GetProjectByOwnerID(owner_id)

	if err != nil {
		return WriteJSON(w, 500, map[string]string{
			"message": err.Error(),
		})
	}

	return WriteJSON(w, 200, categorys)
}
func (s *ApiService) GetProjectsById(w http.ResponseWriter, r *http.Request) error {
	project_id_str := chi.URLParam(r, "project_id")

	// Convert userIDStr to an integer
	project_id, err := strconv.Atoi(project_id_str)
	if err != nil {
		// Handle the error (e.g., invalid integer in the URL parameter)
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return err // Return the error to stop further execution
	}

	project, err := s.store.GetProjectByID(project_id)

	if err != nil {
		return WriteJSON(w, 500, map[string]string{
			"message": err.Error(),
		})
	}

	return WriteJSON(w, 200, project)
}
func (s *ApiService) DeleteProject(w http.ResponseWriter, r *http.Request) error {
	userID := idToken
	if userID == "" {
		return WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Invalid user ID"})
	}

	userId, err := strconv.Atoi(userID)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Invalid user ID format"})
	}

	projectIDstr := chi.URLParam(r, "projecId")

	id, err := strconv.Atoi(projectIDstr)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Invalid project ID format"})
	}
	// check if the portfolio id exist
	// get the owner id and status for verification
	project, err := s.store.GetProjectByID(id)
	if err != nil {
		return err
	}

	// check if the project owner is the one executing the delete
	if project.Owner.ID != userId {
		return WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "This project doesn't belong to the user"})

	}

	// only project with status listed can be deleted
	if project.Status != "listed" {
		return WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Only Listed project can be deleted"})

	}

	err = s.store.DeleteProject(id)

	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Server is unable to delete the data in the database"})

	}

	return WriteJSON(w, http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Successfully deleted user project"})

}

func (s *ApiService) EditProject(w http.ResponseWriter, r *http.Request) error {
	idstr := chi.URLParam(r, "projectId")

	if idstr == "" {
		return fmt.Errorf("projectId is empty")
	}

	id, err := strconv.Atoi(idstr)
	if err != nil {
		return err
	}
	data := models.EditProject{}

	err = json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "Data format is invalid"})

	}

	// skill list
	err = helper.SkillList(data.SkillsID, svc.DB)

	if err != nil {
		if err.Error() == "not exist" {
			return WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "There is skill id does not exist in the database id"})

		}

		return WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})

	}

	// helper for category value
	err = s.store.IsThisCategoryIDExist(data.CategoryID)

	if err != nil {
		if err.Error() == "not exist" {
			return WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "This category id does not exist in the database id"})

		}

		return WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})

	}

	skillDataQuery := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(data.SkillsID)), ","), "[]")

	query := `
	UPDATE projects
	SET title = $1,
		description = $2,
		skills_id = $3,
		price = $4,
		category_id = $5
	WHERE id = $6
`

	// Execute the query with the provided parameters
	_, err = svc.DB.Exec(
		query,
		data.Title,
		data.Description,
		skillDataQuery,
		data.Price,
		data.CategoryID,
		id,
	)

	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err})

	}

	// remove existing attachment db that is not in updated attachment anymore
	return WriteJSON(w, http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Successfully update project"})
}
