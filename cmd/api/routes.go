package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandleFunc("/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	router.HandleFunc("/v1/worlds", app.requirePermission("worlds:write", app.createNewWorldHandler)).Methods("POST")
	router.HandleFunc("/v1/worlds/{id}", app.requirePermission("worlds:read", app.getWorldHandler)).Methods("GET")
	router.HandleFunc("/v1/worlds", app.requirePermission("worlds:read", app.listMyWorldsHandler)).Methods("GET")
	router.HandleFunc("/v1/worlds/{id}", app.requirePermission("worlds:write", app.updateWorldHandler)).Methods("PATCH")
	router.HandleFunc("/v1/worlds/{id}", app.requirePermission("worlds:write", app.deleteWorldHandler)).Methods("DELETE")

	router.HandleFunc("/v1/worlds/{worldId}/characters", app.requirePermission("characters:write", app.createCharacterHandler)).Methods("POST")
	router.HandleFunc("/v1/worlds/{worldId}/characters/{id}", app.requirePermission("characters:read", app.getCharacterHandler)).Methods("GET")
	router.HandleFunc("/v1/worlds/{worldId}/characters", app.requirePermission("characters:read", app.listCharacterHandler)).Methods("GET")
	router.HandleFunc("/v1/worlds/{worldId}/characters/{id}", app.requirePermission("characters:write", app.updateCharacterHandler)).Methods("PATCH")
	router.HandleFunc("/v1/worlds/{worldId}/characters/{id}", app.requirePermission("characters:write", app.deleteCharacterHandler)).Methods("DELETE")

	router.HandleFunc("/v1/users", app.registerUserHandler).Methods("POST")
	//router.HandleFunc("/v1/users/activated", app.activateUserHandler).Me	thods("PUT")

	router.HandleFunc("/v1/tokens/authentication", app.generateAuthenticationTokenHandler).Methods("POST")

	return router
}
