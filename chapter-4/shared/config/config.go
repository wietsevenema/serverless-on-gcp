package config

import (
	"cloud.google.com/go/compute/metadata"
	"log"
	"os"
)

var ProjectID string
var Port string
var Region string
var FunctionServiceAccount string

func init() {
	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}
	ProjectID = os.Getenv("GCP_PROJECT")
	if ProjectID == "" {
		var err error
		ProjectID, err = metadata.ProjectID()
		if err != nil {
			log.Fatal(err)
		}
	}
	Region = os.Getenv("GCP_REGION")
	if Region == "" {
		var err error
		// Cloud Run is a regional resource,
		// the zone is reported with the suffix
		// '-1' instead of the usual '-a', '-b',
		// or '-c'.
		// Example: europe-west1-1
		zone, err := metadata.Zone()
		if err != nil {
			log.Fatal(err)
		}
		Region = zone[:len(zone)-2]
	}
	FunctionServiceAccount = os.Getenv("GCP_SERVICE_ACCOUNT")
	if FunctionServiceAccount == "" {
		email, err := metadata.Email("default")
		if err != nil {
			log.Fatal(err)
		}
		FunctionServiceAccount = email
	}
}
