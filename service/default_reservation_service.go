package service

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/withoutsecondd/kamibooking/internal"
	"github.com/withoutsecondd/kamibooking/model"
	"github.com/withoutsecondd/kamibooking/repository"
	"net/http"
	"time"
)

type DefaultReservationService struct {
	Repository repository.Repository
}

func (s *DefaultReservationService) GetReservations(roomId int64) ([]model.Reservation, error) {
	reservations, err := s.Repository.GetReservationsByRoomId(roomId)
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func (s *DefaultReservationService) PostReservation(res *model.Reservation) error {
	if res.StartTime.Compare(res.EndTime) == 1 {
		return internal.HttpError{
			Err:  errors.New("reservation is invalid: start time of a reservation can't be later than its end time"),
			Code: http.StatusBadRequest,
		}
	}

	var err error
	for i := 0; i < 3; i++ {
		err = s.Repository.CreateReservation(res)
		if err == nil {
			return nil
		}

		var sqlErr *pgconn.PgError
		if errors.As(err, &sqlErr) && sqlErr.Code == "40001" {
			time.Sleep(time.Millisecond * 200)
			continue
		}

		return internal.HttpError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return internal.HttpError{
		Err:  errors.New("transaction failed after maximum retries"),
		Code: http.StatusInternalServerError,
	}
}
