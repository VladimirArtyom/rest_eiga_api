package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type payload map[string]interface{}

func (app* application) readIDParameter(r *http.Request) (int64, error){
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return -1, errors.New("Invalid id parameter")
	}
	return id, nil
} 



func (app* application) readJSON(w http.ResponseWriter, r* http.Request, destination interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&destination)

	if err != nil {

	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError

	switch {
	case errors.As(err, &syntaxError):
		return fmt.Errorf("json: syntax error at character / offset %d", syntaxError.Offset)
	case errors.As(err, &unmarshalTypeError):
		if unmarshalTypeError.Field != "" {
			return fmt.Errorf("json: cannot decode field %q from the given JSON", unmarshalTypeError.Field)
		}
			return fmt.Errorf("json: incorrect JSON type exists at character / offset %d", unmarshalTypeError.Offset)
	case errors.As(err, &invalidUnmarshalError):
		panic(err)
	case errors.Is(err, io.EOF):
		return errors.New("body must not be empty")
	default:
		return err
		}
	} 
	// Read Json content

	// Check errors
	return nil
}

func (app* application) writeJSON(w http.ResponseWriter, payload payload,
																	headers http.Header, statusCode int) error {

	jsonData, err := json.MarshalIndent(payload,"", "\t")
	if err != nil {
		return err
	}
	jsonData = append(jsonData, '\n')

	// The parameters from w.Header and header are different. Alors, tu peux ajouter header dan w.Header
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	
	w.Write(jsonData)
	w.WriteHeader(statusCode)
	return nil
}

