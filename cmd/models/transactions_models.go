package models

type Trasactions struct {
	ID        int     `json:"id"`
	To_userid int     `json:"to_userid"`
	Amount    float64 `json:"amount"`
}
