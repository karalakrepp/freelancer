package models

import "time"

type OwnerInfo struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Location    string `json:"location"`
	PhoneNumber int    `json:"phoneNumber"`
}

type Project struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"desc"`
	SkillsID    []int     `json:"skills_ID"`
	Revizyon    int       `json:"revizyon"`
	Owner       OwnerInfo `json:"Owner"`

	Price      float64   `json:"price"`
	Attachment []string  `json:"attachment"`
	CategoryID int       `json:"category_id"`
	Created_At time.Time `json:"_"`
}

type ProjectLink struct {
	ID           int
	ProjectID    int
	OfferID      int
	IsCustomerOk bool
	IsOwnerOk    bool
}

type CreateProject struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"desc"`
	SkillsID    []int     `json:"skills_ID"`
	Revizyon    int       `json:"revizyon"`
	Price       float64   `json:"price"`
	Attachment  []string  `json:"attachment"`
	CategoryID  int       `json:"category_id"`
	Created_At  time.Time `json:"_"`
}
type EditProject struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"desc"`
	SkillsID    []int    `json:"skills_ID"`
	Price       float64  `json:"price"`
	Attachment  []string `json:"attachment"`
	CategoryID  int      `json:"category_id"`
}
type UserReviewInfo struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Picture   string `json:"picture"`
}

type FilterNeededData struct {
	ID          int
	Title       string
	Description string
	Skill       string
	Owner       OwnerInfo
	Price       float64
	Category    int
}
