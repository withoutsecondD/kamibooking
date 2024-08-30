package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/withoutsecondd/kamibooking/internal"
	"github.com/withoutsecondd/kamibooking/model"
	"net/http"
)

type PostgresqlRepository struct {
	Conn *pgx.Conn
}

func (repo *PostgresqlRepository) GetReservationsByRoomId(roomId int64) ([]model.Reservation, error) {
	query := `SELECT * FROM reservations WHERE room_id = $1`

	rows, err := repo.Conn.Query(context.Background(), query, roomId)
	if err != nil {
		return nil, internal.HttpError{Err: err, Code: http.StatusInternalServerError}
	}
	defer rows.Close()

	reservations := make([]model.Reservation, 0)
	for rows.Next() {
		var reservation model.Reservation

		if err := rows.Scan(&reservation.ID, &reservation.RoomID, &reservation.StartTime, &reservation.EndTime); err != nil {
			return nil, internal.HttpError{Err: err, Code: http.StatusInternalServerError}
		}

		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

func (repo *PostgresqlRepository) CreateReservation(reservation *model.Reservation) error {
	query := `
		INSERT INTO reservations(room_id, start_time, end_time)
		VALUES($1, $2, $3)
	`

	_, err := repo.Conn.Exec(context.Background(), query, reservation.RoomID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		return internal.HttpError{Err: err, Code: http.StatusInternalServerError}
	}

	return nil
}

func (repo *PostgresqlRepository) GetConflictingReservationsCount(reservation *model.Reservation) (int, error) {
	query := `
		SELECT * FROM reservations
        WHERE room_id = $1 AND ((start_time >= $2 AND start_time <= $3) OR (end_time >= $2 AND end_time <= $3))
	`

	rows, err := repo.Conn.Query(context.Background(), query, reservation.RoomID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		return 0, internal.HttpError{Err: err, Code: http.StatusInternalServerError}
	}

	reservations := make([]model.Reservation, 0)

	for rows.Next() {
		var reservation model.Reservation

		if err := rows.Scan(&reservation.ID, &reservation.RoomID, &reservation.StartTime, &reservation.EndTime); err != nil {
			return 0, internal.HttpError{Err: err, Code: http.StatusInternalServerError}
		}

		reservations = append(reservations, reservation)
	}

	return len(reservations), nil
}
