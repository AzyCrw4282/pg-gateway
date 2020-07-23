package db

import (
	"fmt"
	"github.com/deanishe/go-env"
	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
	"os"
)

/*
TODO: //system issue - my pgsql does not accept any incodming connections despite trying everything.
This should not be the case in for other uses. It currently connects to the local
system defined user;So in order to make connecting easier, Postgres.app creates a user with the
same name as your system user when it starts the first time.

See here->https://github.com/PostgresApp/PostgresApp/issues/313#issuecomment-191119317

*/
var DatabaseHost = os.Getenv("pghost")
var DatabaseUser = os.Getenv("pguser")
var DatabasePassword = os.Getenv("pgpassword")
var DatabaseDatabase = os.Getenv("pgdb")
var DatabasePort = env.GetInt("pgport")

func connectOrDie() (conn *pgx.Conn) {
	var err error
	godotenv.Load()
	//var DatabaseHost = os.Getenv("pghost")
	//var DatabaseUser = os.Getenv("pguser")
	//var DatabasePassword = os.Getenv("pgpassword")
	//var DatabaseDatabase = os.Getenv("pgdb")
	//var DatabasePort = 56418

	conn, err = pgx.Connect(pgx.ConnConfig{
		Host:     DatabaseHost,
		Port:     uint16(DatabasePort),
		User:     DatabaseUser,
		Password: DatabasePassword,
		Database: DatabaseDatabase,
	})
	fmt.Println("-> ", DatabaseDatabase, DatabasePort, DatabaseUser)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
	return

}
