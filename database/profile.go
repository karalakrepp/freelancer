package database

import (
	"database/sql"
	"fmt"

	"github.com/karalakrepp/Golang/freelancer-project/models"
)

func (s *PostgresStore) CreateProfile(profile models.CreateUserProfileReq) (*models.UserProfile, error) {
	row := s.DB.QueryRow(createProfile, profile.UserID, profile.Description, profile.Title, profile.Skill, profile.Picture)
	var i models.UserProfile

	err := row.Scan(&i.ID, &i.UserID, &i.Description, &i.Title, &i.Skill, &i.Picture, &i.ProjectCompleted)

	return &i, err

}

func (s *PostgresStore) GetProfile(user_id int) (*models.UserProfile, error) {
	rows, err := s.DB.Query("SELECT * FROM user_profiles WHERE user_id = $1", user_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccountProfile(rows)
	}

	return nil, fmt.Errorf("account with number [%d] not found", user_id)

}

func (s *PostgresStore) UpdateProfile(user_id int, req *models.UpdateProfile) (int, error) {

	query := "UPDATE user_profiles SET description=$2, title=$3 WHERE user_id=$1"
	_, err := s.DB.Exec(query, user_id, req.Description, req.Title)

	if err != nil {
		return 0, err
	}
	return user_id, nil
}

func scanIntoAccountProfile(rows *sql.Rows) (*models.UserProfile, error) {

	var i models.UserProfile
	err := rows.Scan(
		&i.ID,
		&i.UserID,
		&i.Description,
		&i.Title,
		&i.Skill,
		&i.Picture,
		&i.ProjectCompleted,
	)
	if err != nil {
		return nil, err
	}

	// Handle nil slice

	return &i, err
}
