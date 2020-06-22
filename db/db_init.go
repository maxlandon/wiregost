package db

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evilsocket/islazy/tui"
	"github.com/maxlandon/wiregost/db/server"
	"github.com/maxlandon/wiregost/server/assets"
)

// ErrDatabaseDoesNotExist - Used to notify need to create DB
var (
	ErrDatabaseDoesNotExist     = errors.New("Database wiregost_db does not exist")
	ErrWiregostRoleDoesNotExist = errors.New("user 'wiregost' does not exist in PostgreSQL")
)

// CheckPostgreSQLAccess - Verifies PostgreSQL installation and access level
func CheckPostgreSQLAccess() (err error) {

	// Configuration
	conf := assets.ServerConfiguration

	// Test connection
	_, err = server.ConnectPostgreSQL(conf.DBName, conf.DBUser, conf.DBPassword)
	fmt.Println(err)
	if err != nil {
		// Switch between various edge cases
		switch err.Error() {
		case "dial tcp [::1]:5432: connect: connection refused": // Change with Host:Port combination of the server
			fmt.Println(tui.Red("DB: ") + "PostgreSQL service is either not running, or not listening on this port")
		case "could not connect":
			// Ping the database
		case "wrong credentials":
			// Notify user
		case "database does not exist":
			// Create database
			return ErrDatabaseDoesNotExist
		case "pq: role \"wiregost\" does not exist":
			return ErrWiregostRoleDoesNotExist
		case "user postgres does not exist/not available":
			// Notify user we need access to postgres user for psql
		}

		// Handle them, most of the time we just create a new DB because it does not exit yet
	}

	return
}

// InitDatabase - Create Database and sets all needed privileges
func InitDatabase() (err error) {

	// Configuration
	conf := assets.ServerConfiguration

	// Notify we are creating the database
	fmt.Println(tui.Blue("DB:") + " Initializing Wiregost database")
	fmt.Println("     DB name:     wiregost_db")
	fmt.Println("     DB user:     wiregost")
	fmt.Println("     DB password: wiregost")

	// Create temporary SQL file (os/exec is just a mess when passing psql commands)
	saveTo := assets.GetDatabaseDir()

	if _, err := os.Stat(saveTo); os.IsNotExist(err) {
		err = os.MkdirAll(saveTo, os.ModePerm)
		if err != nil {
			return errors.New(fmt.Sprintf("Cannot write to wiregost root directory %s", err))
		}
	}

	fi, err := os.Stat(saveTo)
	if fi.IsDir() {
		filename := "default_db.sql"
		saveTo = filepath.Join(saveTo, filename)
	}

	f, err := os.OpenFile(saveTo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	if _, err = f.WriteString(sqlQueries); err != nil {
		panic(err)
	}

	// Create Database in .wiregost/db directory
	cmd1 := exec.Command("initdb", "-D", assets.GetDatabaseDir())
	res, err := cmd1.Output()
	fmt.Println(tui.Yellow("\nPostgreSQL:"))
	fmt.Println(string(res))

	// Run SQL script
	cmd := exec.Command("psql", "-U", "postgres", "-f", saveTo)
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	fmt.Printf("%sDB:%s Successfully created Database %s with user %s and password %s \n",
		tui.BLUE, tui.RESET, tui.Blue(conf.DBName), tui.Blue(conf.DBUser), tui.Blue(conf.DBPassword))

	// Delete SQL queries file
	f.Close()
	err = os.Remove(saveTo)
	if err != nil {
		fmt.Println(tui.Red("Failed to delete temporary SQL file 'default_db.sql'"))
	}

	return
}

var sqlQueries = fmt.Sprintf(`CREATE DATABASE wiregost_db WITH LOCATION = '%s';
CREATE USER wiregost;
ALTER ROLE wiregost WITH PASSWORD 'wiregost';
GRANT ALL ON DATABASE wiregost_db TO wiregost;
`, assets.GetDatabaseDir())

// var sqlQueries = fmt.Sprintf(`CREATE DATABASE wiregost_db WITH LOCATION = '%s';
// CREATE USER wiregost;
// ALTER ROLE wiregost WITH PASSWORD 'wiregost';
// GRANT ALL ON DATABASE wiregost_db TO wiregost;
// `, assets.GetDatabaseDir())
