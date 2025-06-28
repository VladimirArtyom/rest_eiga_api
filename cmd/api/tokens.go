package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/VladimirArtyom/rest_eiga_api/internal/data"
	"github.com/VladimirArtyom/rest_eiga_api/internal/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	
	var input struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return 
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlainText(v, input.Password)


	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidCredentialResponse(w,r)
			default:
				app.serverErrorResponse(w, r, err)
		}

		return
	}

	//check if the provided password matches the actual
	isMatched, err := user.Password.Matches(input.Password)
	
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	
	if !isMatched {
		app.invalidCredentialResponse(w, r)
		return
	}

	// We generate a token for 1 day expiry time
	token, err := app.models.Tokens.New(user.ID, 1*24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, payload{"authentication_token": token}, nil, http.StatusOK)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}


}
