package main

import (
	"encoding/json"
	"github.com/ARau87/foodsharing_events/database"
	"github.com/gorilla/context"
	"net/http"
)

func (app *application) login(w http.ResponseWriter, r *http.Request)  {

	var user database.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	foundUser, err := user.GetByCredentials(app.Database)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.Logger.Error(err)
		return
	}

	token, err := app.CreateJsonToken(foundUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		app.Logger.Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(token)

}

func (app *application) register(w http.ResponseWriter, r *http.Request){

	var user database.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		app.Logger.Error(err)
		return
	}

	created, err := user.Save(app.Database)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.Logger.Error(err)
		return
	}

	token, err := app.CreateJsonToken(created)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		app.Logger.Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(token)

}

func (app *application) currentUser(w http.ResponseWriter, r *http.Request){

	user := context.Get(r, "user").(*database.User)

	jsonData, err := user.ToJson()
	if err != nil {
		app.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type","application/json")
	w.Write(jsonData)

}
