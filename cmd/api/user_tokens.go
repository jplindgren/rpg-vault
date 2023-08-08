package main

import (
	"net/http"
	"time"

	"github.com/jplindgren/rpg-vault/internal/users"
	"github.com/jplindgren/rpg-vault/internal/validator"
)

func (app application) generateAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	users.ValidateEmail(v, input.Email)
	users.ValidatePasswordPlainText(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.services.Users.GetByEmail(input.Email)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	isAuthenticated, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !isAuthenticated {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.services.Tokens.New(input.Email, time.Hour*24, users.Authentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
