package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.genericNotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.genericIncorrectMethodResponse)
	router.HandlerFunc(http.MethodGet, "/health", app.health)
	router.HandlerFunc(http.MethodPost, "/keeper/login", app.AuthenticateKeeper)
	router.HandlerFunc(http.MethodPost, "/keeper/new", app.CreateAccountHandler)
	router.HandlerFunc(http.MethodPost, "/bees/new", app.CreateBeeHandler)
	router.HandlerFunc(http.MethodGet, "/bees", app.GetBeesHandler)
	//router.HandlerFunc(http.MethodGet, "", app.CreateChildrenHandler)
	return app.AuthenticationMiddleware(router)
	//return router
}

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	wrapped := wrap{
		"status": "online",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version": version,
		},

	}
	err := app.WriteJSONResponse(w, http.StatusOK, wrapped)
	if err != nil {
		//
		return
	}
}