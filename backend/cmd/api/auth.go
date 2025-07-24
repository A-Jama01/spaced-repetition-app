package main

import (
	"net/http"
)


func (app *app) registerHandler(w http.ResponseWriter, r *http.Request) {
	//Get user data from form request
	//Check if username already exists
	//Create user in database
	//Generate JWT for that user
}

func (app *app) loginHandler(w http.ResponseWriter, r *http.Request) {
	//Get user data from form request
	//Check if username exists and password is correct
	//Generate JWT for that user
}

func (app *app) logoutHandler(w http.ResponseWriter, r *http.Request) {
	//Remove JWT or SEND a JWT that expires immediately?
}
