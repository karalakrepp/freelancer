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
