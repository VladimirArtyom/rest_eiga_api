package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, err_message interface{}) {

	payload_data := payload{"error": err_message}

	err := app.writeJSON(w, payload_data, nil, status)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(status)
		return
	}
	return
}

// specific error response for each code

// FailedValidationResponse message is depend on
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "The server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
	return
}

func (app *application) badRequestErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
	return

}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {

	message := "The resource cannot be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
	return
}
