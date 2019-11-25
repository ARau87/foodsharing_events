package main

import (
	"database/sql"
	"os"

	"github.com/ARau87/foodsharing_events/database"
	"github.com/ARau87/foodsharing_events/lib"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type application struct {
	Logger     lib.Logger
	Database   *sql.DB
	Config     Config
	AwsSession *session.Session
	SNSService *sns.SNS
}

func (app *application) setupDatabaseConnection(driver, dsn string) {

	db, err := sql.Open(driver, dsn)
	if err != nil {
		app.Logger.Error(err)
		os.Exit(1)
	}

	app.Database = db

}

func (app *application) CreateJsonToken(user *database.User) ([]byte, error) {
	claims := lib.Claims{
		Email: user.Email,
		Id:    user.Id,
	}

	token, err := claims.CreateToken(app.Config.JwtKey)
	if err != nil {
		return nil, err
	}

	data, err := token.ToJson()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (app *application) CurrentUser(key *lib.AccessKey) (*database.User, error) {

	claims, err := lib.ClaimsFromToken(app.Config.JwtKey, []byte(key.Token))
	if err != nil {
		return nil, err
	}

	user := &database.User{Id: claims.Id}
	user, err = user.GetById(app.Database)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (app *application) CurrentUserId(key *lib.AccessKey) (int, error) {

	claims, err := lib.ClaimsFromToken(app.Config.JwtKey, []byte(key.Token))
	if err != nil {
		return -1, err
	}

	return claims.Id, nil

}
