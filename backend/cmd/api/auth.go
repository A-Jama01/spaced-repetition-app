package main

import (
	"github.com/A-Jama01/spaced-repetition-app/internal/store"
	"net/http"
	"time"
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

	expiration := time.Now().Add(3 * time.Hour).Unix()
	claims := map[string]any{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      expiration,
	}
	_, tokenString, err := app.jwtAuth.Encode(claims)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   int(expiration),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	response := "Successful registration."
	err = app.writeJSON(w, http.StatusCreated, envelope{"auth": response}, nil)
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

	expiration := time.Now().Add(3 * time.Hour).Unix()
	claims := map[string]any{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      expiration,
	}
	_, tokenString, err := app.jwtAuth.Encode(claims)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   int(expiration),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	response := "Successful login."
	err = app.writeJSON(w, http.StatusOK, envelope{"auth": response}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *app) meHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user_id, err := app.getUserIDFromContext(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user_id}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
