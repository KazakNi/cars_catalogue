package db

import (
	"cars/internal/utils"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DBConnection *sqlx.DB

var GetListCarsStmt string = `Select regnum, mark, model, year, name, surname, patronymic
FROM cars
OFFSET $1
LIMIT $2
`

var GetListCarsRegNumStmt string = `Select regnum, mark, model, year, name, surname, patronymic
FROM cars
WHERE regnum = $1
`
var GetListCarsMarkStmt string = `Select regnum, mark, model, year, name, surname, patronymic
FROM cars
WHERE mark = $1`

var GetListCarsModelStmt string = `Select regnum, mark, model, year, name, surname, patronymic
FROM cars
WHERE model = $1`

var GetListCarsOwnerNameStmt string = `Select regnum, mark, model, year, name, surname, patronymic
FROM cars
WHERE name = $1`

var GetListCarsOwnerSurNameStmt string = `Select regnum, mark, model, year, name, surname, patronymic
FROM cars
WHERE surname = $1`

var GetCarFromTestApi string = `Select regnum, mark, model, year, name, surname, patronymic
FROM infoAPI
WHERE regnum = $1
`

func NewDBConnection() (*sqlx.DB, error) {
	host, port, user, password, dbname, driver := utils.GetEnv()

	connUrl := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect(driver, connUrl)

	if err != nil {
		panic(fmt.Sprintf("%s, %s", err, connUrl))
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to DB")
	return db, nil
}

func ExecMigration(db *sqlx.DB, path string) error {

	query, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}
	return nil
}
