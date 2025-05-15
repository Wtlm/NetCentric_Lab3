package models

type User struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	FirstName     string `json:"firstname" binding:"required"`
	LastName      string `json:"lastname"`
	Email         string `json:"email"`
	Avatar        string `json:"avatar"`
	Phone         string `json:"phone"`
	DOB           string `json:"dob"`
	Country       string `json:"country"`
	City          string `json:"city"`
	StreetName    string `json:"street_name"`
	StreetAddress string `json:"street_address"`
}
