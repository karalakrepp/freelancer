package models

type CreateUserProfileReq struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Description string `json:"description"`
	Title       string `json:"title"`
	Skill       string `json:"skill"`
	Picture     string `json:"picture"`
}

type UpdateProfile struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
