package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/VladimirArtyom/rest_eiga_api/internal/validator"
	"github.com/lib/pq"
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

type UserModel struct {
	DB *sql.DB
}

var ErrDuplicateEmail = errors.New("duplicate email")

func (u *UserModel) Insert(user *User) error {

	query := `
		INSERT INTO users (name, email, password_hash, activated)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version
	`
	args := []interface{}{
	  user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
	} 

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		pqErr, ok := err.(*pq.Error); if ok {
		switch {
			case pqErr.Code.Name() == "unique_violation" && pqErr.Constraint == "users_email_key":
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (u *UserModel)GetByEmail(email string) (*User, error ) {
	
	query := `
		SELECT id, created_at, name, email, password_hash, activated, version
		FROM users
		WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := u.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.CreatedAt,
		&user.Name, &user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}


func (u *UserModel)Update(user *User) error {

	query := `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
		WHERE id = $5 AND version = $6 
		RETURNING version
	`

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		pqErr, ok := err.(*pq.Error); if ok {
			switch {
				case pqErr.Code.Name() == "unique_violation" && pqErr.Constraint == "users_email_key":
					return ErrDuplicateEmail
				case errors.Is(err, sql.ErrNoRows):
					return ErrEditConflict
			}
		}

		return err
	}

	return nil
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
