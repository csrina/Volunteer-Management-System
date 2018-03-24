CREATE DATABASE caraway;

\c caraway;

CREATE TABLE family (
    family_id       SERIAL      PRIMARY KEY,
    family_name     TEXT,
    children	    INT
);

CREATE TABLE users (
    user_id         SERIAL      PRIMARY KEY,
    user_role       INT         DEFAULT 1,
    username        TEXT        UNIQUE,
    password        TEXT,
    first_name      TEXT,
    last_name       TEXT,
    email           TEXT,
    phone_number    TEXT,
    family_id       INT		REFERENCES family,
    bonus_hours     FLOAT,
    bonus_note 	    TEXT
);

CREATE TABLE room (
    room_id         SERIAL      PRIMARY KEY,
    room_name       TEXT        UNIQUE,
    teacher_id      INT         REFERENCES users (user_id),
    children	    INT,
    room_num        TEXT
);

CREATE TABLE time_block (
    block_id        SERIAL      PRIMARY KEY,
    block_start     TIMESTAMP,
    block_end       TIMESTAMP,
    room_id         INT			REFERENCES room(room_id),
    modifier        INT			DEFAULT 1,
    title           TEXT    DEFAULT 'Facilitation',
    note            TEXT,
    UNIQUE (block_start, block_end, room_id)
);

CREATE TABLE booking (
    booking_id      SERIAL      PRIMARY KEY,
    block_id        INT         REFERENCES time_block (block_id),
    family_id       INT         REFERENCES family (family_id),
    user_id         INT         REFERENCES users (user_id),
    booking_start   TIMESTAMP,
    booking_end     TIMESTAMP,
    CONSTRAINT 	    unq_booking UNIQUE(block_id, family_id, user_id)
);

CREATE TABLE clocking (
    booking_id      SERIAL	 REFERENCES booking (booking_id),
    clock_in        TIMESTAMP,
    clock_out       TIMESTAMP,
    CONSTRAINT 	    unq_clocking UNIQUE(booking_id, clock_in, clock_out)
);

CREATE TABLE donation (
    donation_id     SERIAL PRIMARY KEY,
    donor_id      INT    REFERENCES family (family_id),
    donee_id    INT    REFERENCES family (family_id),
    amount          FLOAT,
    date_sent       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
