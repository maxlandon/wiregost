package db

import (
	"github.com/maxlandon/wiregost/db/models"
	"github.com/maxlandon/wiregost/server/assets"
)

// CheckPostgreSQLAccess - Verifies PostgreSQL installation and access level
func CheckPostgreSQLAccess() (err error) {
	return
}

// InitDatabase - Create Database and sets all needed privileges
func InitDatabase() (err error) {

	// Configuration
	conf := assets.ServerConfiguration

	// Test connection
	_, err = models.ConnectDatabase(conf.DBName, conf.DBUser, conf.DBPassword)
	if err != nil {
		// Switch between various edge cases

		// Handle them, most of the time we just create a new DB because it does not exit yet
	}

	// Mock queries on certificates, and default user wiregost, to check everything works fine.

	return
}
