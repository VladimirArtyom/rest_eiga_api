package data

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintext *string 
	hash []byte 
}

// Calculates bcrypt hash and stores both the hash and the plaintext version in the struct

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plainTextPassword
	p.hash = hash

	return nil
}

// Checks whether the provided plaintext password matches the hashed passowrd stored 
// in the struct
func (p *password) Matches(plainTextPassword string) (bool, error) {
 	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
			case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
				return false, nil
			default:
				return false, err
		}
	}

	return true, nil
}

