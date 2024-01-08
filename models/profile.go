package models

type CreateUserProfileReq struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Description string `json:"description"`
	Title       string `json:"title"`
	Skills      []int  `json:"skill"`
	Picture     string `json:"picture"`
}

type UpdateProfile struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
type QueryUserProfile struct {
	ID          int
	Owner       OwnerInfo
	UserID      int
	Description string
	Title       string
	Picture     string

	Skills string

	ProjectCompleted int
}
