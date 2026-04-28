package main

import (
	"database/sql"
	"net/http"
	"sugdio/api"
	"sugdio/internal/handlers"
	repository "sugdio/internal/repository/postgres"
	"sugdio/internal/service"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	dbConnString := "user=postgres password=postgres dbname=sugdio sslmode=disable"
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		panic(err)
	}

	repo := repository.NewPostgresRepo(db)

	empService := service.NewEmployeeService(repo, repo, repo, repo, repo, repo)
	shiftService := service.NewShiftService(repo, repo)

	h := handlers.NewHandler(empService, shiftService)

	strictHandler := api.NewStrictHandler(h, nil)

	r := chi.NewRouter()

	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}

	r.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		data, _ := swagger.MarshalJSON()
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/openapi.json"),
	))

	api.HandlerFromMux(strictHandler, r)

	http.ListenAndServe(":8080", r)
}
