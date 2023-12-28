package helper

import (
	"database/sql"
	"errors"
)

func SkillList(s []int, db *sql.DB) error {
	var allSkills []interface{}

	rows, err := db.Query("SELECT id from skills")
	defer rows.Close()

	if err != nil {
		return errors.New("Server is unable to execute query to the database")
	}

	var skillid int
	for rows.Next() {
		err := rows.Scan(&skillid)

		if err != nil {
			return err
		}

		allSkills = append(allSkills, skillid)
	}

	for i := 0; i < len(s); i++ {
		exist := Contains(allSkills, s[i])

		if !exist {
			return errors.New("not exist")
		}
	}

	return nil
}
