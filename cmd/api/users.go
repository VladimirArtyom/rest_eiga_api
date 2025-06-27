package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/VladimirArtyom/rest_eiga_api/internal/data"
	"github.com/VladimirArtyom/rest_eiga_api/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var inputData struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type email_data struct {
		ActivationToken string 
		Email string
		ID int64
	}

	err := app.readJSON(w, r, &inputData)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	// creer une nouvelle utilisateur
	var user *data.User = &data.User{
		Name: inputData.Name,
		Email: inputData.Email,
		Activated: false,
	}

	//Hasher le mot de passe 
	err = user.Password.Set(inputData.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Valider les données
	v := validator.New()
	data.ValidateUser(v, user)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// データをデータベースに保存する
	err = app.models.Users.Insert(user)

	if err != nil {
		switch {
			case errors.Is(err, data.ErrDuplicateEmail):
				v.AddError("email", "a user with this email address already exists")
				app.failedValidationResponse(w, r, v.Errors)
			default:
				app.serverErrorResponse(w, r, err)
		}
			return
	}
	
	token , err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// メールを送信する
		out := &email_data {
				Email: user.Email,
				ActivationToken: token.Plaintext,
				ID: user.ID,
		}

		app.background( func(v interface{}) {
		user, _ := v.(email_data)
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	}, *out)


	//JSONレスポンスを書く
	err = app.writeJSON(w, payload{"user": user}, nil, http.StatusCreated )
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	return 
}


func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		TokenPlainText string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	var v *validator.Validator = validator.New()
	
	data.ValidateToken(v, input.TokenPlainText)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlainText)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrRecordNotFound):
				v.AddError("token", "invalid or expired activation token")
				app.failedValidationResponse(w, r, v.Errors)
			default:
				app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true
	

	// Mise a jour un nouveau utilisateur
	err = app.models.Users.Update(user)
	if err != nil {
			switch {
				case errors.Is(err, data.ErrEditConflict):
					app.editConflictResponse(w, r)
				default:
					app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Supprimer tous les tokens utilisateur
	err = app.models.Tokens.DeleteAllTokensForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//Envoyer l'utilisateur mis a jour au client dans une response JSON
	err = app.writeJSON(w, payload{"user": user}, nil, http.StatusOK )
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
