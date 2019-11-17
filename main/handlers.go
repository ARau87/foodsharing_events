package main

import (
	"encoding/json"
	"github.com/ARau87/foodsharing_events/database"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

func (app *application) createEvent(w http.ResponseWriter, r *http.Request){

	user := context.Get(r, "user").(*database.User)

	var event database.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		app.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event.CreatorId = user.Id
	created , err := event.Save(app.Database)
	if err != nil {
		app.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, err := created.ToJson()
	if err != nil {
		app.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)

}

func (app *application) joinEvent(w http.ResponseWriter, r *http.Request){

	eventId, err := strconv.Atoi(mux.Vars(r)["eventId"])
	if err != nil {
		app.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := context.Get(r, "user").(*database.User)
	event := &database.Event{Id: int(eventId)}
	event, err = event.GetById(app.Database)
	if err != nil {
		app.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = event.AddParticipant(app.Database, user)
	if err != nil {
		app.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (app *application) getEventById(w http.ResponseWriter, r *http.Request){

	eventId, err := strconv.Atoi(mux.Vars(r)["eventId"])
	if err != nil {
		app.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	event := &database.Event{Id: int(eventId)}
	event, err = event.GetById(app.Database)
	if err != nil {
		app.Logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonData, err := event.ToJson()
	if err != nil {
		app.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)


}

