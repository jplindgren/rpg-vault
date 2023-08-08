package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	common "github.com/jplindgren/rpg-vault/internal"
	"github.com/jplindgren/rpg-vault/internal/validator"
	"github.com/jplindgren/rpg-vault/internal/worlds"
)

func (app application) createNewWorldHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name   string
		Genres []string
		Cover  string
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)
	world := &worlds.World{
		UserId:     user.Email,
		Name:       input.Name,
		Genres:     input.Genres,
		CoverImage: input.Cover,
	}

	v := validator.New()
	if worlds.ValidateWorld(v, world); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.services.Worlds.Insert(world)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/worlds/%s", world.Id))
	err = app.writeJSON(w, http.StatusCreated, envelope{"world": world}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app application) getWorldHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user := app.contextGetUser(r)

	world, err := app.services.Worlds.Get(user.Email, id)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrorRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"world": world}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app application) listMyWorldsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	worlds, err := app.services.Worlds.List(user.Email)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	app.writeJSON(w, http.StatusOK, envelope{"worlds": worlds}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app application) updateWorldHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := app.contextGetUser(r)
	id := vars["id"]

	world, err := app.services.Worlds.Get(user.Email, id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if world == nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Name       *string  `json:"name" dynamodbav:"name"`
		Intro      *string  `json:"intro" dynamodbav:"intro"`
		Genres     []string `json:"genres" dynamodbav:"genres,stringset,omitempty"`
		CoverImage *string  `json:"coverImage" dynamodbav:"coverImage"`
	}

	if input.Name != nil {
		world.Name = *input.Name
	}

	if input.Intro != nil {
		world.Intro = *input.Intro
	}

	if input.Genres != nil {
		world.Genres = input.Genres
	}

	imgUpdated := false
	if input.CoverImage != nil {
		imgUpdated = *input.CoverImage != ""
		world.CoverImage = *input.CoverImage
	}

	v := validator.New()
	if worlds.ValidateWorld(v, world); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.services.Worlds.Update(user.Email, id, world, imgUpdated)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/worlds/%s", world.Id))
	app.writeJSON(w, http.StatusOK, envelope{"world": world}, headers)
}

// DeleteWorld deletes a world and all its content (characters, maps, etc)
// swagger:route DELETE /worlds/{id} deleteWorldHandler
// Delete a world.
//
// responses:
//
//	200:
//	400: ErrorResponse
//	500: ErrorResponse
func (app application) deleteWorldHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	worldId := vars["id"]
	user := app.contextGetUser(r)

	cKeys, err := app.services.Characters.ListKeys(worldId)
	if err != nil {
		app.deleteItemResponse(w, r, "Character")
	}

	err = app.services.Characters.DeleteByKeys(cKeys)
	if err != nil {
		app.deleteItemResponse(w, r, "Character")
	}

	err = app.services.Worlds.Delete(user.Email, worldId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, nil, nil)
}
