package main

import (
	"net/http"
	"time"
	"github.com/A-Jama01/spaced-repetition-app/internal/store"
)


func (app *app) registerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input struct {
		Username string `json:"username" validate:"required,min=8,max=40"`
		Password string `json:"password" validate:"required,min=8,max=50"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.validate.Struct(input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: input.Username,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.store.Users.Create(ctx, user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	
	claims := map[string]interface{}{
		"user_id": user.ID,
		"username": user.Username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	_, tokenString, err := app.jwtAuth.Encode(claims)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envolope{"auth_token": tokenString}, nil)	
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *app) loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input struct {
		Username string `json:"username" validate:"required,min=8,max=40"`
		Password string `json:"password" validate:"required,min=8,max=50"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.validate.Struct(input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.store.Users.GetByUsername(ctx, input.Username)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = user.Password.Matches(input.Password)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	claims := map[string]interface{}{
		"user_id": user.ID,
		"username": user.Username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	_, tokenString, err := app.jwtAuth.Encode(claims)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envolope{"auth_token": tokenString}, nil)	
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
