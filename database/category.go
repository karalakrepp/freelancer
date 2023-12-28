package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/karalakrepp/Golang/freelancer-project/models"
)

func (s *PostgresStore) CreateCategory(req models.Category) (int64, error) {

	sqlStatement := `INSERT INTO category(name,parentid) VALUES($1, $2) RETURNING id`

	var id int64
	err := s.DB.QueryRow(sqlStatement, req.Name, req.ParentID).Scan(&id)

	if err != nil {
		return 0, err
	}

	fmt.Printf("Inserted a single record %v", id)

	return id, nil
}
func (s *PostgresStore) GetCategoryByParentId(parent_id int) (*models.Category, error) {

	rows, err := s.DB.Query("SELECT * FROM category WHERE parentid = $1", parent_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoCategory(rows)
	}

	return nil, fmt.Errorf("account with number [%d] not found", parent_id)

}
func (s *PostgresStore) GetCategoryRow() (*[]models.Category, error) {
	allData := []models.Category{}
	result, err := s.DB.Query("SELECT * FROM category")
	if err != nil {
		log.Println("Error querying database:", err)
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		data := models.Category{}
		if err := result.Scan(&data.ID, &data.ParentID, &data.Name); err != nil {
			log.Println("Error scanning result:", err)
			return nil, err
		}
		allData = append(allData, data)
	}

	return &allData, nil
}

func scanIntoCategory(rows *sql.Rows) (*models.Category, error) {
	i := new(models.Category)

	err := rows.Scan(
		&i.ID,
		&i.ParentID,
		&i.Name,
	)
	if err != nil {
		return nil, err
	}

	// Handle nil slice

	return i, err
}
