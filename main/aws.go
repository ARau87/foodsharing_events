package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func (app *application) setupAwsSession(){

	app.AwsSession = session.Must(session.NewSessionWithOptions(
		session.Options{
			Profile: app.Config.AWSProfile,
			Config: aws.Config{Region: aws.String(app.Config.AWSAvailabilityZone)},
		}))

}

func (app *application) setupSNSService(){

	app.SNSService = sns.New(app.AwsSession)

}