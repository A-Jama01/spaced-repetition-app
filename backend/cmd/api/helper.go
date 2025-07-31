package main

import (
	"encoding/json"
	"net/http"
)

type envolope map[string]any

func (app *app) writeJSON(w http.ResponseWriter, status int, data envolope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, val := range headers {
		w.Header()[key] = val
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *app) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	//1MB limit on JSON responses
	maxBytes := 1_048_576 
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	return dec.Decode(data)
}
