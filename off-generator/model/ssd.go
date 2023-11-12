package model

type SSD struct {
	UserId int    `json:"user_id" faker:"-"`
	MAC    string `json:"mac_address" faker:"mac_address"`
}
