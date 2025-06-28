package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	//クライアントに新しい情報を追加します。 (Add new information for the client)
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, err_message interface{}) {

	payload_data := payload{"error": err_message}

	err := app.writeJSON(w, payload_data, nil, status)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(status)
		return
	}
}

// specific error response for each code

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	var message string = "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}

// FailedValidationResponse message is depend on
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	var message string = "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
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

func (app *application) invalidCredentialResponse(w http.ResponseWriter, r *http.Request) {
	message := "Invalid authentication credential"
	app.errorResponse(w, r, http.StatusUnauthorized, message )
	return
}


