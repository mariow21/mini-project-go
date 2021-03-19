package testing

import "time"

// Testing model
type Testing struct {
	UserID     int       `db:"UserID" json:"user_id"`
	Nama       string    `db:"Nama" json:"nama"`
	Usia       int       `db:"Usia" json:"usia"`
	Kota       string    `db:"Kota" json:"kota"`
	LastUpdate time.Time `db:"LastUpdate" json:"last_update"`
}

// User ...
type User struct {
	UserID []int `json:"user_id"`
	UserInsert []Testing `json:"users"`
}
