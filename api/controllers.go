package api

import (
	"cars/internal/db"
	"cars/internal/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func ReDoc(w http.ResponseWriter, r *http.Request) {
	static_path := utils.GetStaticPath()
	tmpl, err := template.ParseFiles(static_path)
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func GetListCars(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	var cars Cars
	cars, err := cars.GetAllCars(db.DBConnection, params)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Logger.Debug("Error while extracting cars: ", "error", err)
		return
	}

	b, _ := json.Marshal(cars)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func CreateCar(w http.ResponseWriter, r *http.Request) {
	var regnums RegNumsInput

	d := json.NewDecoder(r.Body)
	err := d.Decode(&regnums)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Logger.Debug("Error while decoding regnums", "error", err)
		return
	}
	if err = regnums.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Неверный формат данных!"))
		utils.Logger.Info("Error while validating regnums", "error", err, "regnums invalid:", regnums)
		return
	}

	for _, number := range regnums.RegNums {
		var car Car
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/info/%s", number))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.Logger.Debug("Error while ping info endpoint", "error", err)
			return
		}

		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&car)

		err = car.Validate()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			utils.Logger.Info("Error while validating car", "error", err, "car:", car)
			return
		}

		err = car.Create(db.DBConnection)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			utils.Logger.Debug("Error while creating car", "error", err)
			return
		}

	}
	w.WriteHeader(http.StatusCreated)
}

// test endpoint

func GetCarInfo(w http.ResponseWriter, r *http.Request) {
	regnum := strings.TrimPrefix(r.URL.Path, "/info/")

	var car Car
	rows, err := db.DBConnection.Query(db.GetCarFromTestApi, regnum)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Logger.Debug("DB Query error", "error", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			utils.Logger.Debug("DB Scan error", "error", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(car)
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func DeleteCar(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/cars/"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Logger.Debug("Undefined id request", "error", err)
		return
	}

	var car Car
	car, err = car.GetCarById(id, db.DBConnection)

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		utils.Logger.Info("The car's id is not found", "error", err)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Logger.Debug("CarById error while querying", "error", err)
		return
	}

	err = car.Delete(id, db.DBConnection)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Logger.Debug("Error while delete car operation", "error", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateCar(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/cars/"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Logger.Debug("Undefined id request", "error", err)
		return
	}

	var car Car

	d := json.NewDecoder(r.Body)
	err = d.Decode(&car)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Logger.Debug("Error while decoding car binary", "error", err)
		return
	}
	if err = car.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Неверный формат данных!"))
		utils.Logger.Debug("Error while validating car", "error", err)
		return
	}

	err = car.Update(id, db.DBConnection)

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		utils.Logger.Info("The car's id is not found", "error", err)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Logger.Debug("Error while update car operation", "error", err)
		return
	}
	b, _ := json.Marshal(car)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
