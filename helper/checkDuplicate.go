package helper

import (
	"database/sql"
	"errors"
)

func CheckUpdatedUsername(db *sql.DB, username string, id string) error {
	var existingID int
	row := db.QueryRow("SELECT id FROM login WHERE username=? && id !=?", username, id).Scan(&existingID)

	if row == sql.ErrNoRows {
		return nil
	} else if row == nil {
		return errors.New("Username has already been taken")
	} else {
		return errors.New(row.Error())
	}
}
