package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/VladimirArtyom/rest_eiga_api/internal/data"
	"github.com/VladimirArtyom/rest_eiga_api/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Wait a moment, this
	var inputData struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	// After reading the file, check if it fulfilled the bare minimum
	err := app.readJSON(w, r, &inputData)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}
	var v *validator.Validator = validator.New()
	var movie *data.Movie = &data.Movie{
		Title:   inputData.Title,
		Year:    inputData.Year,
		Runtime: inputData.Runtime,
		Genres:  inputData.Genres,
	}
	data.ValidateMovie(v, movie)

	if !v.Valid() {
		// If it has no content
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Make a header ? Pourqoui ? Donc tu ne changes pas directment le w.Header
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, payload{"movie": movie}, headers, http.StatusCreated)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, payload{"movie": movie}, nil, http.StatusOK)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// get the movie, check it
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// faire un lecteur
	var inputData struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	err = app.readJSON(w, r, &inputData)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
	}

	if inputData.Title != nil {
		movie.Title = *inputData.Title
	}

	if inputData.Runtime != nil {
		movie.Runtime = *inputData.Runtime
	}

	if inputData.Year != nil {
		movie.Year = *inputData.Year
	}

	if inputData.Genres != nil {
		movie.Genres = inputData.Genres
	}
	//validate
	var v *validator.Validator = validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	fmt.Println("BERAPA KALI")

	// Ecrire le fichier JSON
	err = app.writeJSON(w, payload{"movies": movie}, nil, http.StatusOK)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, payload{"message": "Movie is sucessfully deleted"}, nil, http.StatusOK)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
