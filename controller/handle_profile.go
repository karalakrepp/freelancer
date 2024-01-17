package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/karalakrepp/Golang/freelancer-project/helper"
	"github.com/karalakrepp/Golang/freelancer-project/models"
)

func (s *ApiService) CreateUserProfile(w http.ResponseWriter, r *http.Request) error {
	userIDStr := idToken

	// Convert userIDStr to an integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		// Handle the error (e.g., invalid integer in the URL parameter)
		http.Error(w, "Invalid userID", http.StatusBadRequest)

	}

	exists, err := s.store.DoesUserProfileExist(userID)
	fmt.Println(exists)
	if err != nil {
		// Handle the error (e.g., database error)
		http.Error(w, "Error checking user profile existence", http.StatusInternalServerError)

	}

	if exists {
		// Handle the case where the user profile exists
		return WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "User profile already exists",
		})
	}
	var req = new(models.CreateUserProfileReq)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Fatal(err)

	}
	req.UserID = userID

	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return json.NewEncoder(w).Encode(err)
	}

	skillDataQuery := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(req.Skills)), ","), "[]")
	fmt.Printf("skiils : %s", skillDataQuery)
	profile, err := s.store.CreateProfile(*req, skillDataQuery)
	if err != nil {
		return json.NewEncoder(w).Encode(err)
	}

	profile.Owner.FirstName = user.FirstName
	profile.Owner.LastName = user.LastName
	profile.Owner.Email = user.Email
	profile.Owner.Location = user.Location.CountryName
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode("created succs")

}
func (s *ApiService) GetUserProfile(w http.ResponseWriter, r *http.Request) error {

	userIDStr := idToken
	var data models.UserProfile
	// Convert userIDStr to an integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		// Handle the error (e.g., invalid integer in the URL parameter)
		http.Error(w, "Invalid userID", http.StatusBadRequest)

	}

	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return err
	}
	profile, err := s.store.GetProfile(userID)
	if err != nil {
		return err
	}
	projectCompleted, err := s.store.GetUserCompletedProject(userID)
	if err != nil {
		return err
	}

	skillData, err := s.store.UserSkills(profile.Skills)

	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error()})

	}

	data.Owner.FirstName = user.FirstName
	data.Owner.LastName = user.LastName
	data.Owner.Email = user.Email
	data.Owner.Location = user.Location.CountryName
	data.ProjectCompleted = projectCompleted
	data.Skill = skillData
	data.UserID = userID
	data.Description = profile.Description
	data.Title = profile.Title
	data.Picture = profile.Picture

	if err != nil {
		return err
	}
	return WriteJSON(w, 200, data)

}

// For title desc
func (s *ApiService) UpdateProfileSection(w http.ResponseWriter, r *http.Request) error {
	userIDStr := chi.URLParam(r, "id")

	// Convert userIDStr to an integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		// Handle the error (e.g., invalid integer in the URL parameter)
		http.Error(w, "Invalid userID", http.StatusBadRequest)

	}
	var req = new(models.UpdateProfile)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return WriteJSON(w, 400, map[string]string{
			"message": "Data format is invalid",
		})
		// Bu satırı ekleyin
	}

	// check empty values
	empty := helper.IsEmptyData(req)
	if empty {
		return WriteJSON(w, 400, map[string]string{
			"message": "There is empty value(s) in the data parameters",
		})

	}

	_, err = s.store.UpdateProfile(userID, req)

	if err != nil {
		return WriteJSON(w, 500, map[string]string{
			"message": "Server is unable to execute query to the database",
		})

	}

	return WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Successfully updating user",
		"id":      userIDStr,
	})
}
