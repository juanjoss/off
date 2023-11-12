package model

type SSD struct {
	Id     int    `db:"id" json:"ssd_id"`
	UserId int    `db:"user_id" json:"user_id"`
	MAC    string `db:"mac_address" json:"mac_address"`
}
