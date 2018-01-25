create database caraway;

\c caraway;

create table room (
    RoomID serial PRIMARY KEY,
    Name text,
    Teacher text,
    RoomNum text
);

/* what else do we want in family? */
create table family (
    FamilyID serial PRIMARY KEY,
    SurName text
);

create table users (
    UserID serial PRIMARY KEY,
    Role int DEFAULT 1,
    FamilyID serial REFERENCES family (familyID),
    UserName text UNIQUE,
    Password text, /* CURRENTLY PLAINTEXT!!! */
    FirstName text,
    LastName text,
    PhoneNumber text /* text for now, probably a better format available */

);

create table block (
    BlockID serial PRIMARY KEY,
    BlockStart TIMESTAMP,
    BlockEnd TIMESTAMP,
    Room serial REFERENCES room(roomID),
    Modifier int DEFAULT 1,
    Note text
);

create table booking (
    BookingID serial PRIMARY KEY,
    BlockID serial REFERENCES block (blockID),
    BamilyID serial REFERENCES family (familyID),
    BookingStart TIMESTAMP,
    BookingEnd TIMESTAMP
);

create table clocking (
    BookingID serial REFERENCES booking (BookingID),
    ClockIn TIMESTAMP,
    ClockOut TIMESTAMP
);

