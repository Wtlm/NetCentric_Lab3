package utils

import (
	"user-management/models"

	"github.com/go-resty/resty/v2"
)

func GetRandomUser() (models.User, error) {
	client := resty.New()
	resp, err := client.R().SetResult(map[string]interface{}{}).
		Get("https://random-data-api.com/api/v2/users")
	if err != nil {
		return models.User{}, err
	}

	data := resp.Result().(map[string]interface{})
	address := data["address"].(map[string]interface{})

	return models.User{
		Username:      data["username"].(string),
		FirstName:     data["first_name"].(string),
		LastName:      data["last_name"].(string),
		Email:         data["email"].(string),
		Avatar:        data["avatar"].(string),
		Phone:         data["phone_number"].(string),
		DOB:           data["date_of_birth"].(string),
		Country:       address["country"].(string),
		City:          address["city"].(string),
		StreetName:    address["street_name"].(string),
		StreetAddress: address["street_address"].(string),
	}, nil
}
