package main

import (
	"flag"
	"fmt"
	"github.com/ARau87/foodsharing_events/lib"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"os"
)

func main(){

	dsn := flag.String("dsn", "root:ttm1306A@/foodsharing?parseTime=true", "The connection string to the database")
	jwtKey := flag.String("jwt", "supersecretKey", "JWT key used for creating access token")
	registerUserTopic := flag.String("registered-user-topic", "arn:aws:sns:eu-central-1:451558607227:ActivateUserNotification", "ARN of the SNS topic that is used to inform the sys admin that a new user has registered" )
	awsProfile := flag.String("aws-profile", "default", "The name of the AWS CLI profile")
	availabilityZone := flag.String("az", "eu-central-1", "The availability zone of the app")

	flag.Parse()

	app := &application{
		lib.Logger{
			Err: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
			Inf: log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime),
		},
		nil,
		Config{
			// TODO: Change this to use cli args in production!
			JwtKey: []byte(*jwtKey),
			DSN: *dsn,
			AWSAvailabilityZone: *availabilityZone,
			AWSProfile: *awsProfile,
			AWSRegisteredUserSNSTopic: *registerUserTopic,
		},
		nil,
		nil,
	}
	app.setupDatabaseConnection("mysql", "root:ttm1306A@/foodsharing?parseTime=true")
	app.setupAwsSession()
	app.setupSNSService()

	app.Logger.Info(fmt.Sprintf("DSN - %s", *dsn))
	app.Logger.Info(fmt.Sprintf("JWT-KEY - %s", *jwtKey))
	app.Logger.Info(fmt.Sprintf("AZ - %s", *availabilityZone))
	app.Logger.Info(fmt.Sprintf("PROFILE - %s", *awsProfile))
	app.Logger.Info(fmt.Sprintf("TOPIC - %s", *registerUserTopic))


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
