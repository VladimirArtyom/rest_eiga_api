package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/VladimirArtyom/rest_eiga_api/internal/data"
)

func (app* application) createMovieHandler(w http.ResponseWriter,r *http.Request){
 // Wait a moment
	var inputData struct {
		Title string `json:"title"`
		Year int32 `json:"year"`
		Runtime int32 `json:"runtime"`
		Genres []string `json:"genres"`
	}

	err := app.readJSON(w, r, inputData)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", inputData)
	fmt.Println("Movies created successfully")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r*http.Request){

	id, err := app.readIDParameter(r)
	app.notFoundResponse(w, r)

	var movieDummy = data.Movie {
		ID: id,
		Title: "Primitif",
		Runtime: 120,
		Genres: []string{"Comedy", "Drama", "俺"},
		Version: 1,
		CreatedAt: time.Now(),
	} 


	err = app.writeJSON(w, payload{"movie": movieDummy}, nil, http.StatusOK)

	app.serverErrorResponse(w, r, err)

	return
}
