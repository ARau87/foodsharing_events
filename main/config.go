package main

type Config struct {

	JwtKey []byte
	RegisteredUserTopic string
	DSN string
	AWSAvailabilityZone string
	AWSRegisteredUserSNSTopic string
	AWSProfile string
}
