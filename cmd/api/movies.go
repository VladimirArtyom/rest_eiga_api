package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/VladimirArtyom/rest_eiga_api/internal/data"
)

func (app* application) createMovieHandler(w http.ResponseWriter,r *http.Request){
	fmt.Fprintln(w, "create movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r*http.Request){

	id, err := app.readIDParameter(r)
	if err != nil {
		
		http.NotFound(w, r)
		return
	}

	var movieDummy = data.Movie {
		ID: id,
		Title: "Primitif",
		Runtime: 120,
		Genres: []string{"Comedy", "Drama", "俺"},
		Version: 1,
		CreatedAt: time.Now(),
	} 


	err = app.writeJSON(w, payload{"movie": movieDummy}, nil, http.StatusOK)

	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return 
	}

	return
}
