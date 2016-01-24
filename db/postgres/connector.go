package db

import (
	"database/sql"
	"fmt"
	"log"
	"peer2peer/config"

	// Ued for connecting to postgres server
	_ "github.com/lib/pq"
)

// GetDBConn : Used for DB connections
func GetDBConn(dbname string) (*sql.DB, error) {

	dbinfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		fmt.Println("Error in connecting to PostGres DB")
		return nil, err
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil

}

// CreateTablesIfNotExists : Creates an index table if not present
func CreateTablesIfNotExists() error {

	log.Println("Validating the presence of tables on server ...")

	// Create DB conn
	db, err := GetDBConn(config.DBName)

	if err != nil {
		log.Println("Error Connecting to DB")
		log.Println(err)
		return err
	}

	// Defer db close
	defer db.Close()

	// Creating the table
	_, err = db.Exec(
		"create table if not exists Visitor(id SERIAL PRIMARY KEY, firstname varchar(50) not null, lastname varchar(50), age int not null, gender varchar(1) not null, email varchar(100) not null, phonenumber varchar(15), university varchar(100), creationtime timestamp default current_timestamp);")

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Tables are ready as required.")

	return nil
}
