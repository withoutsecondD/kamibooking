package repository

import (
	"github.com/withoutsecondd/kamibooking/model"
)

type Repository interface {
	GetReservationsByRoomId(roomId int64) ([]model.Reservation, error)
	CreateReservation(*model.Reservation) error
}
