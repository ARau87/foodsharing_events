package main

import (
	"github.com/ARau87/foodsharing_events/lib"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"os"
)

func main(){

	app := &application{
		lib.Logger{
			Err: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
			Inf: log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime),
		},
		nil,
		[]byte(""),
		Config{
			// TODO: Change this to use cli args in production!
			JwtKey: []byte("supersecretKey"),
		},
		nil,
		nil,
	}
	app.setupDatabaseConnection("mysql", "root:ttm1306A@/foodsharing?parseTime=true")
	app.setupAwsSession()
	app.setupSNSService()

	router := mux.NewRouter()

	// Cors
	router.HandleFunc("/auth/login", app.cors).Methods("OPTIONS")

	// Authentication related routes
	router.Handle("/auth/login", alice.New(app.AllowCors).Then(http.HandlerFunc(app.login))).Methods("POST")
	router.HandleFunc("/auth/register", app.register).Methods("POST")
	router.Handle("/auth", alice.New(app.AuthRequired).Then(http.HandlerFunc(app.currentUser))).Methods("GET")

	// Event related routes
	router.Handle("/event", alice.New(app.AuthRequired).Then(http.HandlerFunc(app.createEvent))).Methods("POST")
	router.Handle("/event/{eventId:[0-9]+}/user/current", alice.New(app.AuthRequired).Then(http.HandlerFunc(app.joinEvent))).Methods("PUT")
	router.Handle("/event/{eventId:[0-9]+}/user/current", alice.New(app.AuthRequired).Then(http.HandlerFunc(app.leaveEvent))).Methods("DELETE")
	router.Handle("/event/{eventId:[0-9]+}", alice.New(app.AuthRequired).Then(http.HandlerFunc(app.getEventById))).Methods("GET")
	router.Handle("/event/{eventId:[0-9]+}", alice.New(app.AuthRequired).Then(http.HandlerFunc(app.deleteEvent))).Methods("DELETE")

	app.Logger.Info("Starting server at address :8080")
	http.ListenAndServe(":8080", router)

	app.Database.Close()


}
