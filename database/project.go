package database

import (
	"database/sql"
	"errors"
	"log"

	"github.com/karalakrepp/Golang/freelancer-project/models"
)

func (s *PostgresStore) CreateProject(param *models.CreateProject, id string, skills string) (int, error) {
	queryResult, err := s.DB.Exec("INSERT INTO projects(title, description, skills, price, owner_id, category_id, created_at) VALUES(?,?,?,?,?,?,?)", param.Title, param.Description, skills, param.Price, id, param.CategoryID, param.Created_At)

	if err != nil {
		return 0, err
	}

	rowID, err := queryResult.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(rowID), nil

}

func (s *PostgresStore) GetAllProject() ([]models.FilterNeededData, error) {
	var allData []models.FilterNeededData

	rows, err := s.DB.Query("SELECT id, title, description,owner_id, skills_id, price, category_id FROM projects WHERE status='listed'")

	if err != nil {
		return []models.FilterNeededData{}, err
	}
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var dbData models.FilterNeededData
		if err := rows.Scan(&dbData.ID, &dbData.Title, &dbData.Description, &dbData.Owner.ID, &dbData.Skill, &dbData.Price, &dbData.Category); err != nil {
			return []models.FilterNeededData{}, errors.New("something is wrong with the database data")
		}

		allData = append(allData, dbData)
	}

	return allData, nil
}

func (s *PostgresStore) GetProjectByCategoryID(categoryID int) (*[]models.FilterNeededData, error) {
	var allData []models.FilterNeededData
	query := "SELECT id, title, description, skills_id, category_id, owner_id FROM projects WHERE category_id = $1"
	rows, err := s.DB.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dbData models.FilterNeededData
		if err := rows.Scan(&dbData.ID, &dbData.Title, &dbData.Description, &dbData.Skill, &dbData.Category, &dbData.Owner.ID); err != nil {
			return &[]models.FilterNeededData{}, errors.New("scan hata")
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

func (s *PostgresStore) GetProjectByOwnerID(owner_id int) (*[]models.FilterNeededData, error) {

	var allData []models.FilterNeededData

	query := "SELECT id, title, description, skills_id, category_id, owner_id from projects WHERE owner_id =$1"

	rows, err := s.DB.Query(query, owner_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var dbData models.FilterNeededData
		if err := rows.Scan(&dbData.ID, &dbData.Title, &dbData.Description, &dbData.Skill, &dbData.Category, &dbData.Owner.ID); err != nil {
			return &[]models.FilterNeededData{}, err
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

func (s *PostgresStore) GetProjectByID(project_id int) (*models.FilterNeededData, error) {

	query := "SELECT id, title, description, skills_id, category_id, owner_id,price,status from projects WHERE id =$1"

	row := s.DB.QueryRow(query, project_id)

	var dbData models.FilterNeededData

	// Scan the values from the row into the variables
	if err := row.Scan(&dbData.ID, &dbData.Title, &dbData.Description, &dbData.Skill, &dbData.Category, &dbData.Owner.ID, &dbData.Price, &dbData.Status); err != nil {
		// Handle the error, e.g., if no rows were found
		if err == sql.ErrNoRows {
			return nil, errors.New("no project found for the given project ID")
		}
		return nil, err
	}

	// Return the single row of data
	return &dbData, nil
}

// Update Project

// Delete Project
func (s *PostgresStore) DeleteProject(project_id int) error {
	query := "DELETE FROM projects WHERE id = $1"

	_, err := s.DB.Exec(query, project_id)
	if err != nil {
		// Log the error for debugging purposes
		log.Printf("Error deleting project with ID %d: %s", project_id, err)
		return err
	}

	// Return nil if the deletion was successful
	return nil
}
