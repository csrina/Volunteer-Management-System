CREATE DATABASE caraway;

\c caraway;

CREATE TABLE users (
    user_id         SERIAL      PRIMARY KEY,
    user_role       INT         DEFAULT 1,
    username        TEXT        UNIQUE,
    password        TEXT,
    first_name      TEXT,
    last_name       TEXT,
    email           TEXT,
    phone_number    TEXT 
);

CREATE TABLE room (
    room_id         SERIAL      PRIMARY KEY,
    room_name       TEXT,
    teacher_id     	INT         REFERENCES users (user_id),
	children		INT,
    room_num        TEXT
);

CREATE TABLE family (
    family_id       SERIAL      PRIMARY KEY,
    family_name     TEXT        UNIQUE,
    parent_one      INT         REFERENCES users (user_id),
    parent_two      INT         REFERENCES users (user_id),
    children		INT
);

CREATE TABLE time_block (
    block_id        SERIAL      PRIMARY KEY,
    block_start     TIMESTAMP,
    block_end       TIMESTAMP,
    room_id         INT			REFERENCES room(room_id),
    modifier        INT			DEFAULT 1,
    note            TEXT
);

CREATE TABLE booking (
    booking_id      SERIAL      PRIMARY KEY,
    block_id        INT         REFERENCES time_block (block_id),
    family_id       INT         REFERENCES family (family_id),
    user_id         INT         REFERENCES users (user_id),
    booking_start   TIMESTAMP,
    booking_end     TIMESTAMP,
    CONSTRAINT unq_booking UNIQUE(block_id, family_id, user_id)
);

CREATE TABLE clocking (
    booking_id      SERIAL      REFERENCES booking (booking_id),
    clock_in        TIMESTAMP,
    clock_out       TIMESTAMP,
    CONSTRAINT 	    unq_clocking UNIQUE(booking_id, clock_in, clock_out)
);

