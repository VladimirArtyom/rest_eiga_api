package main

import (
	"fmt"
	"net/http"
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

	fmt.Fprintf(w, "show movie %d", id)

}
