package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/withoutsecondd/kamibooking/handler"
	"github.com/withoutsecondd/kamibooking/repository"
	"github.com/withoutsecondd/kamibooking/service"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("error reading .env file, reading variables from OS environment")
	}

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRESQL_DB_USER"),
		os.Getenv("POSTGRESQL_DB_PASSWORD"),
		os.Getenv("POSTGRESQL_DB_HOST"),
		os.Getenv("POSTGRESQL_DB_PORT"),
		os.Getenv("POSTGRESQL_DB_NAME"),
	)

	var conn *pgx.Conn
	var err error

	for i := 0; i < 5; i++ {
		conn, err = pgx.Connect(context.Background(), connString)
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
		} else {
			fmt.Println("successfully connected to the database")
			break
		}
	}

	repo := &repository.PostgresqlRepository{Conn: conn}
	reservationS := &service.DefaultReservationService{Repository: repo}
	h := handler.Handler{ReservationS: reservationS}

	r := chi.NewRouter()
	h.SetupRoutes(r)

	http.ListenAndServe(":3000", r)
}
