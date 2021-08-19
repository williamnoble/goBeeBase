package main

import (
	"fmt"
	"net/http"
)

func (app *application) generateErrorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{} ) {
	wrapped := wrap{"error": message}

	err := app.WriteJSONResponse(w, status, wrapped)
	if err != nil {
		app.logError(r.Method + " " +  err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func (app *application) genericNotFoundResponse(w http.ResponseWriter, r *http.Request) {
	friendlyMsg := "Sorry, the requested resource wasn't found"
	app.generateErrorResponse(w, r, http.StatusNotFound, friendlyMsg)
}

func (app *application) genericIncorrectMethodResponse(w http.ResponseWriter, r *http.Request) {
		friendlyMsg := fmt.Sprintf("Sorry, the %s method is not supported for this resource", r.Method)
		app.generateErrorResponse(w, r, http.StatusNotFound, friendlyMsg)
}

func (app *application) genericBadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	//friendlyMsg := "Sorry, that was a bad request"
	app.generateErrorResponse(w,r, http.StatusBadRequest, err.Error() )
}

func (app *application) genericIncorrectCredentialResponse(w http.ResponseWriter, r *http.Request){
	friendlyMsg := "Sorry, the credentials provided were incorrect"
	app.generateErrorResponse(w,r, http.StatusUnauthorized, friendlyMsg )
}

func (app *application) genericInternalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error){
	app.logError("internal server error " + err.Error())
	friendlyMsg := "Sorry, we experience an internal server error, please try again"
	app.generateErrorResponse(w,r, http.StatusInternalServerError, friendlyMsg )
}