package repository

import (
	"github.com/stretchr/testify/mock"
	"github.com/withoutsecondd/kamibooking/model"
)

type Repository interface {
	GetReservationsByRoomId(roomId int64) ([]model.Reservation, error)
	CreateReservation(*model.Reservation) error
	GetConflictingReservationsCount(reservation *model.Reservation) (int, error)
}

type MockRepository struct {
	Mock mock.Mock
}

func (m *MockRepository) GetReservationsByRoomId(roomId int64) ([]model.Reservation, error) {
	args := m.Mock.Called(roomId)
	return args.Get(0).([]model.Reservation), args.Error(1)
}

func (m *MockRepository) CreateReservation(reservation *model.Reservation) error {
	args := m.Mock.Called(reservation)
	return args.Error(0)
}

func (m *MockRepository) GetConflictingReservationsCount(reservation *model.Reservation) (int, error) {
	args := m.Mock.Called(reservation)
	return args.Int(0), args.Error(1)
}
