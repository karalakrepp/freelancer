package database

import (
	"errors"

	"github.com/karalakrepp/Golang/freelancer-project/helper"
	"github.com/karalakrepp/Golang/freelancer-project/models"
)

func (s *PostgresStore) UserSkills(id string) ([]models.UserSkills, error) {
	query, err := helper.SettingInQueryWithID("skills", id, "*")

	if err != nil {
		return []models.UserSkills{}, errors.New(err.Error())
	}

	resp, err := s.DB.Query(query)

	if err != nil {
		return []models.UserSkills{}, errors.New(err.Error())
	}
	defer resp.Close()

	allData := []models.UserSkills{}

	for resp.Next() {
		var databaseData models.UserSkills
		if err := resp.Scan(&databaseData.ID, &databaseData.Name, &databaseData.Created_at, &databaseData.Updated_at); err != nil {
			return []models.UserSkills{}, errors.New("Something is wrong with the database data")
		}

		allData = append(allData, databaseData)
	}

	if resp.Err() != nil {
		return []models.UserSkills{}, errors.New("Something is wrong with the data retrieved")
	}

	return allData, nil
}

func (s *PostgresStore) GetAllSkills() ([]models.UserSkills, error) {
	data, err := s.DB.Query("SELECT * FROM skills")

	var allData []models.UserSkills

	if err != nil {
		return allData, errors.New("server unable to execute query to database")
	}
	defer data.Close()
	for data.Next() {
		// Scan one customer record
		var skills models.UserSkills
		if err := data.Scan(&skills.ID, &skills.Name, &skills.Created_at, &skills.Updated_at); err != nil {
			return []models.UserSkills{}, errors.New("something is wrong with the database data")
		}
		allData = append(allData, skills)
	}
	if data.Err() != nil {
		return []models.UserSkills{}, errors.New("something is wrong with the data retrieved")
	}

	return allData, nil
}

func (s *PostgresStore) GetSkillNames(param string) ([]string, error) {
	var result []string

	if param == "" {
		return []string{}, nil
	}

	initialQuery, err := helper.SettingInQueryWithID("skills", param, "*")

	if err != nil {
		return nil, err
	}

	data, err := s.DB.Query(initialQuery)

	if err != nil {
		return result, errors.New("server unable to execute query to database")
	}
	defer data.Close()
	for data.Next() {
		// Scan one customer record
		var skills models.UserSkills
		if err := data.Scan(&skills.ID, &skills.Name, &skills.Created_at, &skills.Updated_at); err != nil {
			return []string{}, errors.New("something is wrong with the database data")
		}
		result = append(result, skills.Name)
	}
	if data.Err() != nil {
		return []string{}, errors.New("something is wrong with the data retrieved")
	}

	return result, nil
}
func (s *PostgresStore) GetSkillRaw(param string) ([]models.SkillRaw, error) {
	var skillRaw []models.SkillRaw
	initialQuery, err := helper.SettingInQueryWithID("skills", param, "id, name")

	if err != nil {
		return skillRaw, err
	}

	data, err := s.DB.Query(initialQuery)
	defer data.Close()

	if err != nil {
		return skillRaw, errors.New("Server is unable to execute query to the database")
	}

	for data.Next() {
		var skills models.SkillRaw
		if err := data.Scan(&skills.ID, &skills.Name); err != nil {
			return []models.SkillRaw{}, errors.New("Something is wrong with the database data")
		}

		skillRaw = append(skillRaw, skills)
	}

	return skillRaw, nil
}
