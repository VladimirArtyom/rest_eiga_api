package main

import (
	"fmt"
	"net/http"
	"time"

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
	app.notFoundResponse(w, r)

	var movieDummy = data.Movie{
		ID:        id,
		Title:     "Primitif",
		Runtime:   120,
		Genres:    []string{"Comedy", "Drama", "俺"},
		Version:   1,
		CreatedAt: time.Now(),
	}

	err = app.writeJSON(w, payload{"movie": movieDummy}, nil, http.StatusOK)

	app.serverErrorResponse(w, r, err)

	return
}
