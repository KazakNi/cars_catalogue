package api

import (
	"cars/internal/utils"
	"net/http"
)

func SetRoutes() {

	mux := http.NewServeMux()

	GetCarsHandler := http.HandlerFunc(GetListCars)
	CreateCarHandler := http.HandlerFunc(CreateCar)
	deleteCarHandler := http.HandlerFunc(DeleteCar)
	patchActorHandler := http.HandlerFunc(UpdateCar)

	// Swagger specification
	fileDir := utils.GetStaticRoot()
	mux.HandleFunc("GET /redoc", ReDoc)
	mux.Handle("/films.yaml", http.FileServer(http.Dir(fileDir)))

	mux.Handle("GET /cars", LogRequest(GetCarsHandler))
	mux.Handle("POST /cars", LogRequest(CreateCarHandler))
	mux.Handle("DELETE /cars/{id}", LogRequest(deleteCarHandler))
	mux.Handle("PATCH /cars/{id}", LogRequest(patchActorHandler))

	// test api
	mux.HandleFunc("GET /info/{regNum}", GetCarInfo)

	utils.Logger.Info("Starting server")

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		utils.Logger.Error("Error while listening", "error:", err)
	}
}
