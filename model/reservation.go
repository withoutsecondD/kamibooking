package model

import "time"

// Reservation отображает бронь конференц-зала. StartTime и EndTime в JSON
// должны быть формата RFC3339 для корректного декодирования
type Reservation struct {
	ID        int64     `json:"id" db:"id"`
	RoomID    int64     `json:"room_id" db:"room_id"`
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime   time.Time `json:"end_time" db:"end_time"`
}
