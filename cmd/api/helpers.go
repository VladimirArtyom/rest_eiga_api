package main

import (
	"encoding/json"
	"errors"
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

func (app* application) writeJSON(w http.ResponseWriter, payload payload, headers http.Header, statusCode int) error {

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

