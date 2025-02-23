package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter,r *http.Request){
	fmt.Fprintln(w, "status: available")
	fmt.Fprintln(w, "environment: ", app.cfg.env)
	fmt.Fprintln(w, "version:", version)
}

