package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/withoutsecondd/kamibooking/model"
	"github.com/withoutsecondd/kamibooking/repository"
	"sync"
	"testing"
	"time"
)

func TestDefaultReservationService_PostReservation(t *testing.T) {
	// ARRANGE

	testTable := []struct {
		name          string
		args          []*model.Reservation
		setupRepoMock func(m *repository.MockRepository, args []*model.Reservation)
		wantErr       []bool
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
				{
					ID:        0,
					RoomID:    0,
					StartTime: time.Date(2024, 8, 29, 10, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2024, 8, 29, 11, 0, 0, 0, time.UTC),
				},
			},
			setupRepoMock: func(m *repository.MockRepository, args []*model.Reservation) {
				m.Mock.On("GetConflictingReservationsCount", mock.Anything).Return(0, nil)
				m.Mock.On("CreateReservation", mock.Anything).Return(nil)
			},
			wantErr: []bool{false, false, false},
		},
		{
			name: "testing calls with conflict reservations (time intersections)",
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
					StartTime: time.Date(2024, 8, 29, 10, 50, 0, 0, time.UTC),
					EndTime:   time.Date(2024, 8, 29, 11, 30, 0, 0, time.UTC),
				},
			},
			setupRepoMock: func(m *repository.MockRepository, args []*model.Reservation) {
				m.Mock.On("GetConflictingReservationsCount", args[0]).Return(0, nil)
				m.Mock.On("GetConflictingReservationsCount", args[1]).Return(1, nil)
				m.Mock.On("CreateReservation", args[0]).Return(nil)
			},
			wantErr: []bool{false, true},
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
			setupRepoMock: func(m *repository.MockRepository, args []*model.Reservation) {},
			wantErr:       []bool{true},
		},
	}

	for _, tC := range testTable {
		t.Run(tC.name, func(t *testing.T) {
			mockRepo := &repository.MockRepository{}
			tC.setupRepoMock(mockRepo, tC.args)
			reservationS := &DefaultReservationService{Repository: mockRepo}
			a := assert.New(t)

			wG := sync.WaitGroup{}
			for i, arg := range tC.args {
				wG.Add(1)

				arg := arg
				i := i
				go func() {
					err := reservationS.PostReservation(arg)

					if tC.wantErr[i] {
						a.Error(err)
					} else {
						a.NoError(err)
					}

					wG.Done()
				}()
			}

			wG.Wait()
		})
	}
}
