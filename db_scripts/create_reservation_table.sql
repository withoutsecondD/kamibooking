CREATE TABLE reservations(
    id BIGSERIAL PRIMARY KEY,
    room_id BIGINT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL
)