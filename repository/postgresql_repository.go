package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/withoutsecondd/kamibooking/internal"
	"github.com/withoutsecondd/kamibooking/model"
	"net/http"
)

type PostgresqlRepository struct {
	Conn *pgxpool.Pool
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
	tx, err := repo.Conn.Begin(context.Background())
	if err != nil {
		return err
	}

	query := `SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;`

	_, err = tx.Exec(context.Background(), query)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	query = `
		SELECT * FROM reservations
        WHERE room_id = $1 AND ((start_time >= $2 AND start_time <= $3) OR (end_time >= $2 AND end_time <= $3));
	`

	rows, err := tx.Query(context.Background(), query, reservation.RoomID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	conflicts := make([]model.Reservation, 0)

	for rows.Next() {
		var reservation model.Reservation

		if err := rows.Scan(&reservation.ID, &reservation.RoomID, &reservation.StartTime, &reservation.EndTime); err != nil {
			return err
		}

		conflicts = append(conflicts, reservation)
	}

	if len(conflicts) > 0 {
		return internal.HttpError{
			Err:  errors.New("reservations conflict: some reservations have conflicting timestamps"),
			Code: http.StatusConflict,
		}
	}

	query = `
		INSERT INTO reservations(room_id, start_time, end_time)
		VALUES($1, $2, $3);
	`

	_, err = tx.Exec(context.Background(), query, reservation.RoomID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	return nil
}
