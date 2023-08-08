package main

import (
	"net/http"

	common "github.com/jplindgren/rpg-vault/internal"
	"github.com/jplindgren/rpg-vault/internal/users"
	"github.com/jplindgren/rpg-vault/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &users.User{
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: common.GetIsoString(),
		Activated: true,
		Version:   1,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if users.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.services.Users.Insert(user)
	if err != nil {
		switch {
		// case errors.Is(err, data.ErrorDuplicateEmail):
		// 	v.AddError("email", "a user with this email address already exists")
		// 	app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
