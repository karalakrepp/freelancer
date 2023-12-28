package database

import (
	"database/sql"
	"fmt"

	"github.com/karalakrepp/Golang/freelancer-project/models"
)

const (
	createProfile = `INSERT INTO user_profiles (
		user_id,description,
		title, skill, picture
	  ) VALUES (
		$1, $2, $3,$4,$5
	  )
	  RETURNING id, user_id, description, title, skill,picture, project_completed
	  
	`
	createUser = `INSERT INTO Laccounts (
		first_name,last_name,
		username, email, password
	  ) VALUES (
		$1, $2, $3,$4,$5
	  )
	  RETURNING id, username, email, password, balance,is_admin, is_active, deactivated_at, is_deleted, deleted_at, created_at
	  
	`
)

type CreateUserParams struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	PhoneNumber    int    `json:"phoneNumber"`
	Location       string `json:"location"`
}

// Firstname ve LastName daha sonra cekilecek
func (s *PostgresStore) CreateAccount(user CreateUserParams) (*models.User, error) {
	row := s.DB.QueryRow(createUser, user.FirstName, user.LastName, user.Username, user.Email, user.HashedPassword)

	var i models.User

	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.Balance,
		&i.IsAdmin,
		&i.IsActive,
		&i.DeactivatedAt,
		&i.IsDeleted,
		&i.DeletedAt,
		&i.CreatedAt,
	)

	return &i, err
}

func (s *PostgresStore) GetUserByEmail(email string) (*models.User, error) {
	rows, err := s.DB.Query("select * from Laccounts where email = $1", email)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with number [%s] not found", email)
}
func scanIntoAccount(rows *sql.Rows) (*models.User, error) {
	i := new(models.User)

	err := rows.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.PhoneNumber,
		&i.Balance,
		&i.Location.CountryName,
		&i.IsAdmin,
		&i.UserType,
		&i.IsActive,
		&i.DeactivatedAt,
		&i.IsDeleted,
		&i.DeletedAt,
		&i.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Handle nil slice

	return i, err
}

//After Auth

func (s *PostgresStore) GetUserByID(user_id int) (*models.User, error) {

	rows, err := s.DB.Query("SELECT * FROM Laccounts WHERE id = $1", user_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with number [%d] not found", user_id)

}
