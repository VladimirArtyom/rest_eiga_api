package data

import (
	"time"

	"github.com/VladimirArtyom/rest_eiga_api/internal/validator"
)

type User struct {
	ID int64 `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password password `json:"-"`
	Activated bool `json:"activated"`
	Version int `json:"-"`
}


func ValidateEmail(v *validator.Validator, user *User) {

	//ユーザのメール
	v.Check(user.Email != "", "email", "must be provided")
	v.Check(validator.Matches(user.Email, validator.EmailRX),"email", "must be a valid email address"	)
}

func ValidatePasswordPlainText(v *validator.Validator, password string) {
	//ユーザのパスワート"
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must be at most 72 bytes long")
	
}


func ValidateUser(v *validator.Validator, user *User) {
	
	//ユーザの名前
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <=500, "name", "must be lower than 500 bytes")

	ValidateEmail(v, user)

	if user.Password.plaintext != nil {
		ValidatePasswordPlainText(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("Missing hashed password")
	}

} 
