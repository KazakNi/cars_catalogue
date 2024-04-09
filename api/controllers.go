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
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func load_secret() string {

	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../example.env")
		if err != nil {
			panic("error while loading .env file")
		}
	}
	return os.Getenv("TOKEN_SECRET")
}

func ReDoc(w http.ResponseWriter, r *http.Request) {
	static_path := utils.GetStaticPath()
	tmpl, err := template.ParseFiles(static_path)
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

/*
	func DeleteCar(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/actors/"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Undefined id request Err: %s", err)
			return
		}
		var actor Actor
		actor, err = actor.GetActorById(id, db.DBConnection)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print(err)
			return
		}
		err = actor.Delete(id, db.DBConnection)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Error while delete actor operation: %s", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

	func UpdateActor(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/actors/"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Undefined id request Err: %s", err)
			return
		}

		actor := &Actor{}
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		err = d.Decode(actor)

		if err != nil {
			log.Printf("Error while %s endpoint response body parsing: %s", r.URL, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = actor.GetActorById(id, db.DBConnection)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print(err)
			return
		}

		err = actor.Update(id, db.DBConnection)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Error while update actor operation: %s", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

	func CreateFilm(w http.ResponseWriter, r *http.Request) {
		film := &PostFilm{}
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		err := d.Decode(film)

		if err != nil {
			log.Printf("Error while %s endpoint response body parsing: %s", r.URL, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		var f = Film{Name: film.Name,
			Description:  film.Description,
			Release_date: film.Release_date,
			Rating:       film.Rating}

		film_id, err := f.Create(db.DBConnection)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Error while creating an film: %s", err)
			return
		}

		err = f.InsertCast(film_id, film.Actors_list, db.DBConnection)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Error while inserting actors: %s", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		resp, err := json.Marshal(CreatedId{Id: film_id})

		if err != nil {
			log.Printf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(resp)

}

	func DeleteFilm(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/films/"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Undefined id request Err: %s", err)
			return
		}
		var film Film
		film, err = film.GetFilmById(id, db.DBConnection)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print(err)
			return
		}
		err = film.Delete(id, db.DBConnection)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Error while delete actor operation: %s", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

	func UpdateFilm(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/films/"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Undefined id request Err: %s", err)
			return
		}

		film := &PostFilm{}
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		err = d.Decode(film)

		if err != nil {
			log.Printf("Error while %s endpoint response body parsing: %s", r.URL, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var f = Film{Name: film.Name,
			Description:  film.Description,
			Release_date: film.Release_date,
			Rating:       film.Rating}

		_, err = f.GetFilmById(id, db.DBConnection)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print(err)
			return
		}

		err = f.Update(id, film.Actors_list, db.DBConnection)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Error while update actor operation: %s", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
*/
func GetListCars(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	var cars Cars
	cars, err := cars.GetAllCars(db.DBConnection, params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error while extracting cars: %s", err)
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
		log.Print(err)
		return
	}
	if err = regnums.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Неверный формат данных!"))
		log.Print(err)
		return
	}

	for _, number := range regnums.RegNums {
		var car Car
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/info/%s", number))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&car)

		err = car.Validate()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Неверный формат данных")
			return
		}

		err = car.Create(db.DBConnection)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
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
		log.Println("DB Query error: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
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
		log.Printf("Undefined id request Err: %s", err)
		return
	}

	var car Car
	car, err = car.GetCarById(id, db.DBConnection)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		log.Print("car's id is not found: ", id)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	err = car.Delete(id, db.DBConnection)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error while delete car operation: %s", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateCar(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/cars/"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Undefined id request Err: %s", err)
		return
	}

	var car Car

	d := json.NewDecoder(r.Body)
	err = d.Decode(&car)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}
	if err = car.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Неверный формат данных!"))
		log.Print(err)
		return
	}

	err = car.Update(id, db.DBConnection)

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		log.Print("car's id is not found: ", id)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error while update car operation: %s", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
