package db

import "github.com/maxlandon/wiregost/db/models"

// Start - Starts one or more components of the Data Service
func Start() error {

	// Load config for:
	// - PostgreSQL credentials
	// - gRPC service options
	// - gRPC-REST gateway options

	// Connect to DB
	db, _ := models.ConnectDatabase("", "", "")

	// Load certificates/key pairs (stored in DB)

	// Migrate Schema
	MigrateShema(db)

	// Register gRPC services

	// Start listening components (gRPC and/or REST)

	return nil
}

// StartRESTGateway - Start listening for REST requests
func StartRESTGateway() error {
	return nil
}

// StopRESTGateway - Stop the REST server
func StopRESTGateway() error {
	return nil
}
