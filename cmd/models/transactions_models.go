package models

import "time"

type Trasactions struct {
	ToUserID  int       `json:"to_userid"`
	AtUserID  int       `json:"at_userid"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
