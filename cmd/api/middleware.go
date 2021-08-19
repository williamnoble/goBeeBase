package main

import (
	"context"
	"goBeeBase/cmd/data"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
)

var userContextKey = "user"

type key struct {
	value uint
}


func (app *application) AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			notAuth := []string{
				"/keeper/new", "/keeper/login", "/health",
			}
			path := r.URL.Path

			for _, permitted := range notAuth {
				if path == permitted {
					next.ServeHTTP(w,r)
					return
				}
			}
			authorizationHeader := r.Header.Get("Authorization")
			plaintextToken, valid := app.validateToken(authorizationHeader)
			if !valid {
				app.genericIncorrectCredentialResponse(w,r)
				return
			}

			tk := &data.Token{}
			token, err := jwt.ParseWithClaims(plaintextToken, tk, func(token *jwt.Token) (interface{}, error){
				return []byte(os.Getenv("token_password")),nil
			})

			if err != nil {
				app.genericIncorrectCredentialResponse(w,r)
				return
			}

			if !token.Valid {
				app.genericIncorrectCredentialResponse(w,r)
				return
			}

			k := &key{
				value: tk.UserID,
			}


			ctx := context.WithValue(r.Context(), userContextKey, k)
			r = r.WithContext(ctx)
			app.infoLogger.Printf("Keeper %d authenticated successfully", tk.UserID)
			next.ServeHTTP(w, r)


	})
}