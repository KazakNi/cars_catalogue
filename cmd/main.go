package main

import (
	"cars/api"
	"cars/internal/db"
	"cars/internal/utils"
)

func main() {
	migration_path := utils.GetMigrPath()
	db.DBConnection, _ = db.NewDBConnection()
	db.ExecMigration(db.DBConnection, migration_path)
	api.SetRoutes()
}
