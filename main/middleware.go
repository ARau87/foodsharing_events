package main

import (
	"github.com/ARau87/foodsharing_events/lib"
	"github.com/gorilla/context"
	"net/http"
)

func (app *application) AuthRequired(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		accessKey := &lib.AccessKey{Token:token}

		user, err := app.CurrentUser(accessKey)
		if err != nil {
			app.Logger.Error(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		context.Set(r, "user", user)
		next.ServeHTTP(w, r)

	})
}

func (app *application) AllowCors(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)

	})

}
