package api

import (
	DB "cars/internal/db"
	"database/sql"
	"log"
	"net/url"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

type RegNumsInput struct {
	RegNums []string `json:"regNums" validate:"required,dive,alphanum"`
}

func (r *RegNumsInput) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)

	if err != nil {
		return err
	}
	return nil
}

type Owner struct {
	Name       string  `json:"name" db:"name" validate:"required"`
	Surname    string  `json:"surname" db:"surname" validate:"required"`
	Patronymic *string `json:"patronymic,omitempty" db:"patronymic"`
}

type Car struct {
	RegNum string `json:"regnum" db:"regnum" validate:"required"`
	Mark   string `json:"mark" db:"mark" validate:"required"`
	Model  string `json:"model" db:"model" validate:"required"`
	Year   int    `json:"year,omitempty" db:"year"`
	Owner  Owner  `json:"owner"`
}

type CarDB struct {
	RegNum string `json:"regnum" db:"regnum" validate:"required"`
	Mark   string `json:"mark" db:"mark" validate:"required"`
	Model  string `json:"model" db:"model" validate:"required"`
	Year   int    `json:"year,omitempty" db:"year"`
	Owner
}

func (c *Car) Validate() error {
	validate := validator.New()
	err := validate.Struct(c)

	if err != nil {
		return err
	}
	return nil
}
func (c *Car) Create(db *sqlx.DB) (err error) {

	_, err = db.Exec("INSERT INTO cars (regNum, mark, model, year, name, surname, patronymic) VALUES ($1, $2, $3, $4, $5, $6, $7)", c.RegNum, c.Mark, c.Model, c.Year, c.Owner.Name, c.Owner.Surname, c.Owner.Patronymic)
	if err != nil {
		return err
	}
	return nil
}

func (c *Car) GetCarById(id int, db *sqlx.DB) (Car, error) {
	var ID int
	err := db.Get(&ID, "SELECT id FROM cars WHERE id = $1", id)

	if err != nil {
		log.Printf("error while querying car: %s", err)
		return *c, err
	} else {
		return *c, nil
	}
}
func (c *Car) Delete(id int, db *sqlx.DB) error {
	_, err := db.Exec("DELETE FROM cars WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return err
	}
	return nil
}

func (c *Car) Update(id int, db *sqlx.DB) error {

	_, err := db.Exec(`UPDATE cars SET regnum = COALESCE($1, regnum), mark = COALESCE($2, mark), model = COALESCE($3, model),  year = COALESCE($4, year), name = COALESCE($5, name), 
					   surname = COALESCE($6, surname), patronymic = COALESCE($7, patronymic) WHERE id = $8 RETURNING id`,
		c.RegNum, c.Mark, c.Model, c.Year, c.Owner.Name, c.Owner.Surname, c.Owner.Patronymic, id)

	if err != nil {
		log.Println("err while car updating", err)
		return err
	}
	return nil
}

type Cars struct {
	Page     int   `json:"page"`
	NextPage int   `json:"next_page,omitempty"`
	PrevPage int   `json:"prev_page,omitempty"`
	Cars     []Car `json:"cars"`
}

func (c *Cars) GetPaginationParams(params url.Values) (p, rp int) {
	pageparam := params.Get("page")
	resppageparam := params.Get("resperpage")

	if len(pageparam) > 0 && len(resppageparam) > 0 {
		page, err := strconv.Atoi(pageparam)
		if err != nil {
			return 0, 5
		}
		resppage, err := strconv.Atoi(resppageparam)
		if err != nil {
			return page, 5
		}
		return page, resppage
	}

	return 0, 5 // 5 результатов на странце по умолчанию
}

func (c *Cars) GetValuesFromDB(carsDB []CarDB) {
	for _, car := range carsDB {
		var normalCar Car
		normalCar.RegNum = car.RegNum
		normalCar.Mark = car.Mark
		normalCar.Model = car.Model
		normalCar.Year = car.Year
		normalCar.Owner.Name = car.Owner.Name
		normalCar.Owner.Surname = car.Owner.Surname
		normalCar.Owner.Patronymic = car.Owner.Patronymic

		c.Cars = append(c.Cars, normalCar)
	}
}

func (c *Cars) GetAllCars(db *sqlx.DB, params url.Values) (Cars, error) {
	carsDB := []CarDB{}
	res := Cars{}

	var rowsNum int

	page, resppage := c.GetPaginationParams(params)
	offset := page*resppage - resppage
	if offset < 0 {
		offset = 0
	}

	regNum := params.Get("regNum")
	mark := params.Get("mark")
	model := params.Get("model")
	owner_name := params.Get("owner_name")
	owner_surname := params.Get("owner_surname")

	if len(params.Get("regNum")) > 0 {

		err := db.Select(&carsDB, DB.GetListCarsRegNumStmt, regNum)
		if err != nil {
			return res, err
		}
		rowsNum = 1
		res.GetValuesFromDB(carsDB)
	}

	if len(params.Get("mark")) > 0 {

		err := db.Select(&carsDB, DB.GetListCarsMarkStmt, mark)
		if err != nil {
			return res, err
		}
		err = db.Get(&rowsNum, "SELECT COUNT(id) from cars WHERE mark = $1", mark)
		if err != nil {
			return res, err
		}
		res.GetValuesFromDB(carsDB)
	}

	if len(params.Get("model")) > 0 {

		err := db.Select(&carsDB, DB.GetListCarsModelStmt, model)
		if err != nil {
			return res, err
		}
		err = db.Get(&rowsNum, "SELECT COUNT(id) from cars WHERE model = $1", model)
		if err != nil {
			return res, err
		}
		res.GetValuesFromDB(carsDB)
	}

	if len(params.Get("owner_name")) > 0 {

		err := db.Select(&carsDB, DB.GetListCarsOwnerNameStmt, owner_name)
		if err != nil {
			return res, err
		}
		err = db.Get(&rowsNum, "SELECT COUNT(id) from cars WHERE name = $1", owner_name)
		if err != nil {
			return res, err
		}
		res.GetValuesFromDB(carsDB)
	}

	if len(params.Get("owner_surname")) > 0 {

		err := db.Select(&carsDB, DB.GetListCarsOwnerSurNameStmt, owner_surname)
		if err != nil {
			return res, err
		}
		err = db.Get(&rowsNum, "SELECT COUNT(id) from cars WHERE surname = $1", owner_surname)
		if err != nil {
			return res, err
		}
		res.GetValuesFromDB(carsDB)
	}

	// if no params

	if regNum == "" && mark == "" && model == "" && owner_name == "" && owner_surname == "" {
		err := db.Select(&carsDB, DB.GetListCarsStmt, offset, resppage)
		if err != nil {
			return res, err
		}

		err = db.Get(&rowsNum, "SELECT COUNT(id) from cars")
		if err != nil {
			return res, err
		}
		res.GetValuesFromDB(carsDB)
	}

	if page != 0 && resppage != 0 {
		if page*resppage-resppage < rowsNum {
			if page > 1 {
				res.PrevPage = page - 1
				res.Page = page
			}
			if rowsNum/(page*resppage) > 0 && rowsNum/(page*resppage) > 1 {
				res.NextPage = page + 1
				res.Page = page
			}
			res.Page = 1
			return res, nil
		} else {
			res.Page = 1
			return res, nil
		}
	}
	res.Page = 1
	return res, nil
}
