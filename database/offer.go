package database

import (
	"errors"
	"fmt"

	"github.com/karalakrepp/Golang/freelancer-project/models"
)

func (s *PostgresStore) CreateOffer(offer *models.Offer) (int64, error) {
	var id int64

	// Prepare the SQL statement
	stmt, err := s.DB.Prepare(createOffer)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Execute the SQL statement and scan the result into id
	err = stmt.QueryRow(
		offer.CustomerID,
		offer.CustomerNote,
		offer.ProjectOwnerID,
		offer.ProjectID,
		offer.Price,
		offer.Status,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	fmt.Printf("Inserted a single record with ID: %v\n", id)

	return id, nil
}

func (s *PostgresStore) GetAllOfferByOwnerId(ownerId int) (*[]models.Offer, error) {
	var allData []models.Offer
	query := "Select * from offers where owner_id =$1"

	rows, err := s.DB.Query(query, ownerId)

	if err != nil {
		return &[]models.Offer{}, nil
	}
	defer rows.Close()

	for rows.Next() {
		var dbData models.Offer
		if err := rows.Scan(&dbData.ID, &dbData.CustomerID, &dbData.CustomerNote, &dbData.ProjectOwnerID, &dbData.ProjectID, &dbData.Status, &dbData.Price); err != nil {
			return &[]models.Offer{}, err
		}
		allData = append(allData, dbData)
	}

	// rows.Err() ile olası bir tarama hatasını kontrol et
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// rows.Next() içindeki tarama hatasını kontrol et
	if len(allData) == 0 {
		return nil, errors.New("no projects found for the given category ID")
	}

	return &allData, nil
}

func (s *PostgresStore) GetAllOfferByCustomerId(customerId int) (*[]models.Offer, error) {
	var allData []models.Offer
	query := "Select * from offers where customer_id =$1"

	rows, err := s.DB.Query(query, customerId)

	if err != nil {
		return &[]models.Offer{}, nil
	}
	defer rows.Close()

	for rows.Next() {
		var dbData models.Offer
		if err := rows.Scan(&dbData.ID, &dbData.CustomerID, &dbData.CustomerNote, &dbData.ProjectOwnerID, &dbData.ProjectID, &dbData.Status, &dbData.Price); err != nil {
			return &[]models.Offer{}, err
		}
		allData = append(allData, dbData)
	}

	// rows.Err() ile olası bir tarama hatasını kontrol et
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// rows.Next() içindeki tarama hatasını kontrol et
	if len(allData) == 0 {
		return nil, errors.New("no projects found for the given category ID")
	}

	return &allData, nil
}
func (s *PostgresStore) OfferIsDone(offerid int) error {

	query := "UPDATE offers SET status = 'Done' WHERE id =$1"

	_, err := s.DB.Exec(query, offerid)

	if err != nil {
		return err

	}

	return nil

}

func (s *PostgresStore) IsThisYourOffer(offerID int, ownerID int) bool {
	row := s.DB.QueryRow("SELECT owner_id FROM offers WHERE id = $1 AND owner_id = $2", offerID, ownerID)

	var foundOwnerID int
	err := row.Scan(&foundOwnerID)
	if err != nil {
		return false
	}

	// If a row is found, and ownerID matches foundOwnerID, it is the user's offer
	return foundOwnerID == ownerID
}

func (s *PostgresStore) GetUserCompletedProject(id int) (int, error) {
	var count int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM offers WHERE owner_id=$1 AND status='Done'", id).Scan(&count)

	if err != nil {
		return -1, errors.New("server unable to execute query to database")
	}

	return count, nil
}
