package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"math"
	"time"

	"github.com/VladimirArtyom/rest_eiga_api/internal/validator"
)

const (
	ScopeActivation = "activation"
	ScopeAuthentication = "authentication"
)

const(
	tokenLength = 16
)

// Parce que cette token est utilisee dans JSON
type Token struct {
	Plaintext string `json:"token"`
	Hash      []byte `json:"-"`
	UserID    int64	 `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string `json:"-"`
}

type TokenModel struct {
	DB *sql.DB
}

func ValidateToken(v *validator.Validator, tokenPlainText string) {
	v.Check(tokenPlainText != "", "token", "has to have a value")
	v.Check(len(tokenPlainText) == int(math.Ceil(float64((tokenLength * 8)) / float64(5))) , "token", "has to have 26 bytes long")
}

func (t *TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}	
	
	err = t.Insert(token)
	if err != nil {
		return nil, err
	}

	return token, nil

}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	var token *Token = &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	var randomBytes []byte = make([]byte, tokenLength)
	_, err := rand.Read(randomBytes)

	if err != nil {
		return nil, err
	}
	
	token.Plaintext =  base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	var hash [32]byte = sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	
	return token, nil
}


func (t* TokenModel) Insert(token *Token) error{
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)`

	args := []interface{}{
		token.Hash,
		token.UserID,
		token.Expiry,
		token.Scope,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err :=	t.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (t *TokenModel) DeleteAllTokensForUser(scope string, userID int64) error {
	var query string = `DELETE FROM tokens
		WHERE scope=$1 AND user_id=$2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{
		scope,
		userID,
	};
	_, err := t.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}




