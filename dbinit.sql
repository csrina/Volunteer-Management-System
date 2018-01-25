create database caraway;

\c caraway;

create table room (
    room_id serial PRIMARY KEY,
    name text,
    teacher text,
    room_num text
);

/* what else do we want in family? */
create table family (
    family_id serial PRIMARY KEY,
    surname text
);

create table users (
    user_id serial PRIMARY KEY,
    role int DEFAULT 1,
    family_id serial REFERENCES family (family_id),
    username text UNIQUE,
    password text, /* CURRENTLY PLAINTEXT!!! */
    firstname text,
    lastname text,
    phonenumber text /* text for now, probably a better format available */

);

create table block (
    block_id serial PRIMARY KEY,
    block_start TIMESTAMP,
    block_end TIMESTAMP,
    room serial REFERENCES room(room_id),
    modifier int DEFAULT 1,
    note text
);

create table booking (
    booking_id serial PRIMARY KEY,
    block_id serial REFERENCES block (block_id),
    family_id serial REFERENCES family (family_id),
    booking_start TIMESTAMP,
    booking_end TIMESTAMP
);

create table clocking (
    booking_id serial REFERENCES booking (booking_id),
    clock_in TIMESTAMP,
    clock_out TIMESTAMP
);

