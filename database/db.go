package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/karalakrepp/Golang/freelancer-project/models"
	_ "github.com/lib/pq"
)

const DATABASE_URL = "postgres://postgres:12345678@localhost:5433/freelancer?sslmode=disable"

var DB *sql.DB

type PostgresStore struct {
	DB *sql.DB
}

type Storage interface {
	CreateAccount(CreateUserParams) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)
	GetUserByID(int) (*models.User, error)

	CreateProfile(models.CreateUserProfileReq, string) (*models.UserProfile, error)
	GetProfile(int) (*models.QueryUserProfile, error)
	UpdateProfile(int, *models.UpdateProfile) (int, error)
	DoesUserProfileExist(int) (bool, error)

	UpdateBalance(userID int, newBalance float64) error

	CreateCategory(models.Category) (int64, error)
	GetCategoryByParentId(int) (*models.Category, error)
	GetCategoryRow() (*[]models.Category, error)

	IsThisCategoryIDExist(id int) error

	CreateProject(*models.CreateProject, string, string) (int, error)
	GetAllProject() ([]models.FilterNeededData, error)
	GetProjectByCategoryID(int) (*[]models.FilterNeededData, error)
	GetProjectByOwnerID(int) (*[]models.FilterNeededData, error)
	GetProjectByID(int) (*models.FilterNeededData, error)
	DeleteProject(int) error

	GetUserCompletedProject(id int) (int, error)

	//offer
	CreateOffer(*models.Offer) (int64, error)
	GetAllOfferByOwnerId(int) (*[]models.Offer, error)
	GetAllOfferByCustomerId(int) (*[]models.Offer, error)
	OfferIsDone(int) error
	IsThisYourOffer(int, int) bool
	GetCustomersOfferDone(int) (*[]models.Offer, error)

	GetOfferById(int) (*models.Offer, error)

	//Project Links
	AddProjectLinks(*models.Offer, bool) error
	IsCustomerOk(id int) error
	GetProjectLink(id int) (*models.ProjectLink, error)
	//Skills
	GetAllSkills() ([]models.UserSkills, error)
	UserSkills(id string) ([]models.UserSkills, error)
}

func NewPostgresStore() (*PostgresStore, error) {
	var err error
	DB, err = sql.Open("postgres", DATABASE_URL)

	if err != nil {
		return nil, err
	}

	if err := DB.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		DB: DB,
	}, nil

}
func (s *PostgresStore) createAccountTable() error {
	query := `
	CREATE TABLE  IF NOT EXISTS Laccounts (
		"id" bigserial PRIMARY KEY NOT NULL,
		"first_name" varchar NOT NULL,
		"last_name" varchar NOT NULL,
		"username" varchar UNIQUE NOT NULL,
		"email" varchar UNIQUE NOT NULL,
		"password" varchar NOT NULL,
		"phone" serial NOT NULL,
		"balance" float DEFAULT 0.0,
		"location" varchar DEFAULT 'turkey',
		"is_admin" boolean NOT NULL DEFAULT false,
		"usertype" varchar NOT NULL DEFAULT 'customer',
		"is_active" boolean NOT NULL DEFAULT true,
		"deactivated_at" timestamptz NOT NULL DEFAULT ('0001-01-01 00:00:00Z'),
		"is_deleted" boolean NOT NULL DEFAULT false,
		"deleted_at" timestamptz NOT NULL DEFAULT ('0001-01-01 00:00:00Z'),
		"created_at" timestamptz NOT NULL DEFAULT (now())
	);
	
	CREATE TABLE IF NOT EXISTS user_profiles (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		description TEXT,
		title VARCHAR(255),
		skill VARCHAR(255),
		picture VARCHAR(255),
		project_completed INTEGER NOT NULL DEFAULT 0,
	
		FOREIGN KEY (user_id) REFERENCES Laccounts(id) 
	);

	CREATE TABLE IF NOT EXISTS category (
		id SERIAL PRIMARY KEY,
		parentid INTEGER NOT NULL,
	
		name varchar(100)
	
	);

	CREATE TABLE IF NOT EXISTS project_link (
		id SERIAL PRIMARY KEY,


		project_id INTEGER NOT NULL,
		offer_id   INTEGER NOT NULL,

		isCustomerOk bool default false,
		isOwnerOk    bool default false,

		FOREIGN KEY (offer_id) REFERENCES offers(id)
		
	
	);
		
	CREATE TABLE IF NOT EXISTS projects (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		owner_id    varchar(100),
		skills_id varchar(100),
		price NUMERIC(10, 2),
		attachment TEXT[],
		status varchar(50) NOT NULL DEFAULT 'listed',
		category_id INTEGER,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

		FOREIGN KEY (category_id) REFERENCES category(id) 
	);
	
	CREATE TABLE IF NOT EXISTS offers (
		id SERIAL PRIMARY KEY,
		customer_id INT NOT NULL,
		customer_note TEXT,
		owner_id INT NOT NULL,
		project_id INT NOT NULL,
		price INT NOT NULL,
		status VARCHAR(255) NOT NULL,
		FOREIGN KEY (project_id) REFERENCES projects(id)
	);
	
	CREATE TABLE IF NOT EXISTS skills (
		ID SERIAL PRIMARY KEY,
		Name VARCHAR(255) NOT NULL,
		Created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		Updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
		
		
		
		
	`

	_, err := s.DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("migration başarılı")
	return nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}
