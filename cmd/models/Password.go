package models

type Password struct {
	NewPassword string `json:"new_password"`
	Password    string `json:"password"`
}
