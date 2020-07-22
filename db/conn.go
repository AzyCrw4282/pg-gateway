package db

import (
	"fmt"
	"github.com/deanishe/go-env"
	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
	"os"
)

var DatabaseHost = os.Getenv("pghost")
var DatabaseUser = os.Getenv("pguser")
var DatabasePassword = os.Getenv("pgpassword")
var DatabaseDatabase = os.Getenv("pgdb")
var DatabasePort = env.GetInt("pgport")

func connectOrDie() (conn *pgx.Conn) {
	var err error
	godotenv.Load()
	conn, err = pgx.Connect(pgx.ConnConfig{
		Host:     DatabaseHost,
		Port:     uint16(DatabasePort),
		User:     DatabaseUser,
		Password: DatabasePassword,
		Database: DatabaseDatabase,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
	return

}
