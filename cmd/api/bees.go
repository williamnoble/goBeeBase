package main

import (
	"fmt"
	"goBeeBase/cmd/data"
	"net/http"
)

func(app *application) CreateBeeHandler(w http.ResponseWriter, r *http.Request) {
	bee := data.Bee{}

	err := app.readJSON(w,r, &bee)
	if err != nil {
		app.genericBadRequestResponse(w,r, err)
		return
	}

	userID, _ := r.Context().Value(userContextKey).(*key)
	x := userID.value
	bee.KeeperID = x


	err = app.models.Bees.Insert(&bee)
	if err != nil {
		app.genericInternalServerErrorResponse(w,r, err)
	}

	wrapped := wrap{
		"bee": bee,
	}

	var caste string
	switch bee.Caste {
	case 2:
		caste = fmt.Sprintf("Rejoice, the new %s is Born!!", bee.Caste)
	default:
		caste = fmt.Sprintf("A %s is Born", bee.Caste)
	}
	app.infoLogger.Print(caste)

	err = app.WriteJSONResponse(w, http.StatusCreated, wrapped)
	if err != nil {
		app.genericInternalServerErrorResponse(w,r, err)
	}

}


func (app *application) GetBeesHandler(w http.ResponseWriter, r *http.Request) {
	contextKeyUserID, _ := r.Context().Value(userContextKey).(*key)
	keeperID := contextKeyUserID.value

	bees, err := app.models.Bees.GetAll(keeperID)
	if err != nil {
		app.genericInternalServerErrorResponse(w,r, err)
		return

	}


	wrapped := wrap{
		"bee": bees,
	}

	err = app.WriteJSONResponse(w, http.StatusOK, wrapped)
	if err != nil {
		app.genericInternalServerErrorResponse(w,r, err)
	}

}