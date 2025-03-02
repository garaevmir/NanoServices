package models

type Register struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateProfile struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Birthdate   string `json:"birthdate"`
	PhoneNumber string `json:"phone_number"`
	Bio         string `json:"bio"`
}
