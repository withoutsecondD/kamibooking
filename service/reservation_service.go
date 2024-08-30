package service

import "github.com/withoutsecondd/kamibooking/model"

type ReservationService interface {
	GetReservations(roomId int64) ([]model.Reservation, error)
	PostReservation(res *model.Reservation) error
}
