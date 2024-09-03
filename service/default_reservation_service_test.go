package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/withoutsecondd/kamibooking/model"
	"github.com/withoutsecondd/kamibooking/repository"
	"os"
	"sync"
	"testing"
	"time"
)

func TestDefaultReservationService_PostReservation(t *testing.T) {
	// ARRANGE

	testTable := []struct {
		name    string
		args    []*model.Reservation
		wantErr bool
	}{
		{
			name: "testing calls with valid reservations (without intersections)",
			args: []*model.Reservation{
				{
					ID:        0,
					RoomID:    1,
					StartTime: time.Date(2024, 8, 29, 10, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2024, 8, 29, 11, 0, 0, 0, time.UTC),
				},
				{
					ID:        0,
					RoomID:    1,
					StartTime: time.Date(2024, 8, 29, 11, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2024, 8, 29, 12, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name: "testing concurrent calls with conflict reservations (time intersections)",
			args: []*model.Reservation{
				{
					ID:        0,
					RoomID:    1,
					StartTime: time.Date(2024, 8, 28, 10, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2024, 8, 28, 11, 0, 0, 0, time.UTC),
				},
				{
					ID:        0,
					RoomID:    1,
					StartTime: time.Date(2024, 8, 28, 10, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2024, 8, 28, 11, 30, 0, 0, time.UTC),
				},
			},
			wantErr: true,
		},
		{
			name: "testing call with invalid reservations (EndTime is before StartTime)",
			args: []*model.Reservation{
				{
					ID:        0,
					RoomID:    1,
					StartTime: time.Date(2024, 8, 29, 11, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2024, 8, 29, 10, 0, 0, 0, time.UTC),
				},
			},
			wantErr: true,
		},
	}

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRESQL_DB_USER"),
		os.Getenv("POSTGRESQL_DB_PASSWORD"),
		os.Getenv("POSTGRESQL_DB_HOST"),
		os.Getenv("POSTGRESQL_DB_PORT"),
		os.Getenv("POSTGRESQL_DB_NAME"),
	)

	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		fmt.Println("error connecting to the database")
		return
	}
	repo := &repository.PostgresqlRepository{Conn: conn}
	s := &DefaultReservationService{Repository: repo}

	wG := sync.WaitGroup{}

	for _, tC := range testTable {
		t.Run(tC.name, func(t *testing.T) {
			// ACT
			actualErr := false

			for _, arg := range tC.args {
				wG.Add(1)

				go func(arg *model.Reservation) {
					err := s.PostReservation(arg)
					if err != nil {
						actualErr = true
						fmt.Println(err)
					}

					wG.Done()
				}(arg)
			}

			wG.Wait()

			if actualErr != tC.wantErr {
				t.Logf("expected error: %v, got: %v", tC.wantErr, actualErr)
				t.Fail()
			}
		})
	}
}
