package models

type Task struct {
	UserID      int    `json:"user_id"`
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
