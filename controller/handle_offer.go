package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/karalakrepp/Golang/freelancer-project/models"
)

func (s *ApiService) addOffer(w http.ResponseWriter, r *http.Request) error {
	var customer_id = idToken
	project_id_str := chi.URLParam(r, "projectID")
	customerId, err := strconv.Atoi(customer_id)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid customerID"})
	}

	project_id, err := strconv.Atoi(project_id_str)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid projectID"})
	}

	// Convert userIDStr to an integer

	var req = new(models.Offer)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return WriteJSON(w, 400, err)
	}
	project, err := s.store.GetProjectByID(project_id)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid project"})
	}

	fmt.Printf("project :%v", project)
	req.ProjectID = project_id
	req.Status = "Beklemede"
	req.Price = project.Price
	req.ProjectOwnerID = project.Owner.ID
	req.CustomerID = customerId

	id, err := s.store.CreateOffer(req)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return WriteJSON(w, 200, map[string]interface{}{
		"message": "Offer is on",
		"offerID": id,
	})
}

// for projectowners
func (s *ApiService) getOfferByOwnerId(w http.ResponseWriter, r *http.Request) error {

	var ownerIdstr = idToken

	ownerID, err := strconv.Atoi(ownerIdstr)

	if err != nil {
		return err
	}

	offers, err := s.store.GetAllOfferByOwnerId(ownerID)

	if err != nil {
		return err
	}

	return WriteJSON(w, 200, offers)
}

// for customer orders

func (s *ApiService) getOfferByCustomerID(w http.ResponseWriter, r *http.Request) error {
	var customerIdstr = idToken

	customerId, err := strconv.Atoi(customerIdstr)

	if err != nil {
		return err
	}

	offers, err := s.store.GetAllOfferByCustomerId(customerId)

	if err != nil {
		return err
	}

	return WriteJSON(w, 200, offers)
}

func (s *ApiService) offerIsDone(w http.ResponseWriter, r *http.Request) error {
	offerIDStr := chi.URLParam(r, "offerID")
	ownerID := idToken

	offerID, err := strconv.Atoi(offerIDStr)
	if err != nil {
		return err
	}

	ownerIDInt, err := strconv.Atoi(ownerID)
	if err != nil {
		return err
	}

	fmt.Println(ownerID)

	// Check if this offer is available on this account
	isThis := s.store.IsThisYourOffer(offerID, ownerIDInt)

	fmt.Printf("owner is : %d", ownerIDInt)

	if !isThis {
		return WriteJSON(w, http.StatusNotFound, map[string]interface{}{
			"message": "It is not your offer",
			"check":   isThis,
			"code":    http.StatusBadRequest,
		})
	}
	fmt.Printf(" is : %t", isThis)

	// Update offer status
	err = s.store.OfferIsDone(offerID)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": "Internal Server Error"})
	}

	return WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "Offer status updated to 'Done'"})
}
