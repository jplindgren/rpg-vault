package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jplindgren/rpg-vault/internal/characters"
)

func (app application) createCharacterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	worldId := vars["worldId"]

	var input struct {
		Name       string
		Intro      string
		OwnerId    string
		CoverImage string
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	character := &characters.Character{
		WorldId:    worldId,
		Name:       input.Name,
		Intro:      input.Intro,
		OwnerId:    input.OwnerId,
		CoverImage: input.CoverImage,
	}

	err = app.services.Characters.Insert(character)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/worlds/%s/characters/%s", worldId, character.Id))
	app.writeJSON(w, http.StatusCreated, envelope{"character": character}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app application) updateCharacterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	worldId := vars["worldId"]
	id := vars["id"]

	character, err := app.services.Characters.Get(worldId, id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if character == nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Name       *string
		Intro      *string
		Attributes *string
		CoverImage *string
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		character.Name = *input.Name
	}
	if input.Intro != nil {
		character.Intro = *input.Intro
	}
	if input.Attributes != nil {
		character.AttributesJSON = *input.Attributes
		err = json.Unmarshal([]byte(character.AttributesJSON), &character.Attributes)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
	}
	if input.CoverImage != nil {
		character.CoverImage = *input.CoverImage
	}

	err = app.services.Characters.Update(worldId, id, character)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/worlds/%s/characters/%s", worldId, character.Id))
	app.writeJSON(w, http.StatusOK, envelope{"character": character}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCharacterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	worldId := vars["worldId"]
	id := vars["id"]

	result, err := app.services.Characters.Get(worldId, id)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"character": result}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app application) listCharacterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	worldId := vars["worldId"]

	characters, err := app.services.Characters.List(worldId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"characters": characters}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app application) deleteCharacterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	worldId := vars["worldId"]
	id := vars["id"]

	err := app.services.Characters.Delete(worldId, id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
