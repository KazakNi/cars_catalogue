package utils

import (
	"flag"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var Local_dev *bool // false for docker and true for local stage
var Logger = GetLogger()

func DevMode() *bool {
	local_dev := flag.Bool("devmode", false, "true triggers for local and false for docker env")
	flag.Parse()
	return local_dev
}
func GetLogger() *slog.Logger {

	var lvl = new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))
	return logger
}
func GetMigrPath() string {
	if *Local_dev {
		return "../internal/migrations/migrations.sql"
	} else {
		return "./internal/migrations/migrations.sql"
	}
}

func GetStaticPath() string {
	if *Local_dev {
		return "../api/static/redoc.html"
	} else {
		return "./api/static/redoc.html"
	}
}

func GetStaticRoot() string {

	if *Local_dev {
		return "../api/static"
	} else {
		return "./api/static"
	}
}

func GetEnv() (hostDB, portDB, userDB, passwordDB, dbnameDB, driverDB string) {

	var host, port, user, password, dbname, driver string
	//docker mode
	err := godotenv.Load("example.env")
	// local mode
	if err != nil {
		err = godotenv.Load("../example.env")
		if err != nil {
			panic("error while loading .env file")
		}

		host = os.Getenv("HOST")
	} else {
		host = os.Getenv("HOST")
	}

	port = os.Getenv("PORT")
	user = os.Getenv("USER")
	password = os.Getenv("PASSWORD")
	dbname = os.Getenv("DB_NAME")
	driver = os.Getenv("DRIVER")

	return host, port, user, password, dbname, driver
}
