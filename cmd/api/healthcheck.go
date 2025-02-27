package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter,r *http.Request){

	data := map[string]interface{} {
		"status": "available",
		"system_info": map[string]string{
			"environment": app.cfg.env,
			"version": version,
		},
	}

	
	err := app.writeJSON(w, payload{"data": data}, nil, http.StatusOK) // Tu peux changer l'en-tete plus tard
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
	
	return
}

