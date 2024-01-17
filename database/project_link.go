package database

import "github.com/karalakrepp/Golang/freelancer-project/models"

func (s *PostgresStore) AddProjectLinks(offer *models.Offer, isOwner bool) error {
	query := "INSERT INTO project_link (project_id, offer_id, isownerok) VALUES ($1, $2, $3)"
	_, err := s.DB.Exec(query, offer.ProjectID, offer.ID, isOwner)
	if err != nil {
		return err
	}
	return nil

}

func (s *PostgresStore) IsCustomerOk(id int) error {

	query := "UPDATE project_link SET iscustomerok = true WHERE id =$1"
	_, err := s.DB.Exec(query, id)

	if err != nil {
		return err

	}

	return nil

}

func (s *PostgresStore) GetProjectLink(id int) (*models.ProjectLink, error) {

	query := "Select * from project_link WHERE id =$1"
	rows := s.DB.QueryRow(query, id)

	var dbData models.ProjectLink

	if err := rows.Scan(&dbData.ID, &dbData.ProjectID, &dbData.OfferID, &dbData.IsCustomerOk, &dbData.IsOwnerOk); err != nil {
		return &models.ProjectLink{}, err
	}

	return &dbData, nil

}
