package database

import (
	"errors"

	"github.com/karalakrepp/Golang/freelancer-project/helper"
)

func (s *PostgresStore) IsThisCategoryIDExist(id int) error {
	var allCategories []interface{}

	rows, err := s.DB.Query("SELECT id from category")
	if err != nil {
		return err
	}
	defer rows.Close()

	if err != nil {
		return errors.New("server is unable to execute query to the database")
	}

	var categoryID int
	for rows.Next() {
		err := rows.Scan(&categoryID)

		if err != nil {
			return err
		}

		allCategories = append(allCategories, categoryID)
	}

	exist := helper.Contains(allCategories, id)

	if !exist {
		return errors.New("not exist")
	}

	return nil

}
