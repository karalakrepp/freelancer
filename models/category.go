package models

type Category struct {
	ID       int    `json:"id"`
	Name     string `json:"category_name"`
	ParentID int    `json:"parentid"`
}

// duzelicek yanlis  yazildi
type GetCategoryByParentIdReq struct {
	ParentID int `json:"parentid"`
}
