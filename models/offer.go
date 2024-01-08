package models

type Offer struct {
	ID             int     `json:"id"`
	CustomerID     int     `json:"customer_id"`
	CustomerNote   string  `json:"customer_note"`
	ProjectOwnerID int     `json:"owner_id"`
	ProjectID      int     `json:"project_id"`
	Price          float64 `json:"price"`
	Status         string  `json:"status"`
}
