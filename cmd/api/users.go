package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/VladimirArtyom/rest_eiga_api/internal/data"
	"github.com/VladimirArtyom/rest_eiga_api/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var inputData struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
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

	// メールを送信する

	go func(user data.User) {

		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			app.logger.PrintError(err, nil)
		}

	}(*user)

	//JSONレスポンスを書く
	err = app.writeJSON(w, payload{"user": user}, nil, http.StatusCreated )
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	
	return 
}
