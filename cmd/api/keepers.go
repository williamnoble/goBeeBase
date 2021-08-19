package main

import (
	"goBeeBase/cmd/data"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

func (app *application) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	keeper := &data.Keeper{}

	err := app.readJSON(w, r, &keeper)
	if err != nil {
		app.genericBadRequestResponse(w,r, err)
		return
	}

	k, err := app.models.Keepers.Create(keeper)
	if err != nil {
		app.genericInternalServerErrorResponse(w,r, err)
		return
	}
	wrapped := wrap{
		"user": k,
	}
	app.infoLogger.Printf("a new keeper %d - %s was created", k.KeeperID, k.KeeperEmail)
	err = app.WriteJSONResponse(w, http.StatusCreated, wrapped)
	if err != nil {
		app.genericInternalServerErrorResponse(w,r, err)
		return
	}
}

func (app *application) AuthenticateKeeper(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.genericInternalServerErrorResponse(w,r, err)
		return
	}

	k, err := app.models.Keepers.GetByEmail(input.Email)
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.genericNotFoundResponse(w, r)
			return
		}
		app.genericInternalServerErrorResponse(w,r, err)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(k.Password), []byte(input.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		app.genericIncorrectCredentialResponse(w, r)
		return
	}


	claim := &data.Token{UserID: k.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claim)
	tokenToString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	k.Token = tokenToString
	wrapped := wrap{
		"keeper": k,
	}

	err = app.WriteJSONResponse(w, http.StatusOK, wrapped)
	if err != nil {
		app.genericInternalServerErrorResponse(w,r, err)
		return
	}
}


