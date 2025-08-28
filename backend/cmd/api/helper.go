package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
)

var (
	ErrUserIDNotFloat64 = errors.New("User ID is not a float64")
)

type envelope map[string]any

func (app *app) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
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

func (app *app) getUserIDFromContext(ctx context.Context) (int64, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, err
	}

	floatUserID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, ErrUserIDNotFloat64
	}

	userID := int64(floatUserID)

	return userID, nil
}

func (app *app) readString(queryString url.Values, key string, defaultValue string) string {
	s := queryString.Get(key)
	if s == "" {
		return defaultValue
	}

	return s
}

func (app *app) readInt(queryString url.Values, key string, defaultValue int64) int64 {
	s := queryString.Get(key)		
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return int64(i)
}
