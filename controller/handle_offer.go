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

	//add project_links for  customers approval

	offer, err := s.store.GetOfferById(offerID)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": "Internal Server Error"})
	}
	if err := s.store.AddProjectLinks(offer, true); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "Offer status updated to 'Done'"})
}

//get all done offer

func (s *ApiService) getDoneOfferByCustomer(w http.ResponseWriter, r *http.Request) error {
	var customerIdstr = idToken

	customerId, err := strconv.Atoi(customerIdstr)

	if err != nil {
		return err
	}

	offers, err := s.store.GetCustomersOfferDone(customerId)

	if err != nil {
		return err
	}

	return WriteJSON(w, 200, offers)
}

func (s *ApiService) customerIsOK(w http.ResponseWriter, r *http.Request) error {
	project_id_str := chi.URLParam(r, "id")

	id, err := strconv.Atoi(project_id_str)
	if err != nil {
		return WriteJSON(w, 400, err)
	}

	if err := s.store.IsCustomerOk(id); err != nil {
		return WriteJSON(w, 500, err)
	}

	//get offer
	link, err := s.store.GetProjectLink(id)

	if link.IsCustomerOk {
		return WriteJSON(w, 400, map[string]interface{}{
			"message": "Your order already confirmed",
		})
	}

	if err != nil {
		return WriteJSON(w, 500, err)
	}

	offer, err := s.store.GetOfferById(link.OfferID)
	if err != nil {
		return WriteJSON(w, 500, err)
	}

	//get offers profile

	offer_profile, err := s.store.GetUserByID(offer.CustomerID)
	if err != nil {
		return WriteJSON(w, 500, err)
	}
	fmt.Println(offer_profile)

	project_owner_profile, err := s.store.GetUserByID(offer.ProjectOwnerID)
	if err != nil {
		return WriteJSON(w, 500, err)
	}

	//payment it is simple paymnet with db
	if err := s.store.UpdateBalance(offer.CustomerID, offer_profile.Balance-offer.Price); err != nil {
		return WriteJSON(w, 500, err)
	}

	// Corrected the parameter here to use offer.ProjectOwnerID
	errs := s.store.UpdateBalance(offer.ProjectOwnerID, project_owner_profile.Balance+offer.Price)
	if errs != nil {
		return WriteJSON(w, 500, errs)
	}

	return WriteJSON(w, 200, map[string]interface{}{
		"message": "proccess succesfull",
	})

}
