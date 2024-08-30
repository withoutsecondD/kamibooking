package service

import (
	"errors"
	"github.com/withoutsecondd/kamibooking/internal"
	"github.com/withoutsecondd/kamibooking/model"
	"github.com/withoutsecondd/kamibooking/repository"
	"net/http"
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

	conflicts, err := s.Repository.GetConflictingReservationsCount(res)
	if err != nil {
		return err
	}

	if conflicts > 0 {
		return internal.HttpError{
			Err:  errors.New("reservations conflict: some reservations have conflicting timestamps"),
			Code: http.StatusConflict,
		}
	}

	return s.Repository.CreateReservation(res)
}
