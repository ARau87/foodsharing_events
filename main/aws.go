package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func (app *application) setupAwsSession(){

	app.AwsSession = session.Must(session.NewSessionWithOptions(
		session.Options{
			Profile: "default",
			Config: aws.Config{Region: aws.String("eu-central-1")},
		}))

}

func (app *application) setupSNSService(){

	app.SNSService = sns.New(app.AwsSession)

}