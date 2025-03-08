package models

type Users struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
}
