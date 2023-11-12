package model

type User struct {
	FirstName string `json:"first_name" faker:"first_name"`
	LastName  string `json:"last_name" faker:"last_name"`
	Email     string `json:"email" faker:"email"`
	Password  string `json:"password" faker:"password"`
}
