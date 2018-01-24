create table room (
    roomID serial PRIMARY KEY,
    name text,
    teacher text,
    roomNum text
);

/* what else do we want in family? */
create table family (
    familyID serial PRIMARY KEY,
    name text
);

create table block (
    blockID serial PRIMARY KEY,
    blockStart TIMESTAMP,
    blockEnd TIMESTAMP,
    room serial REFERENCES room(roomID),
    modifier int DEFAULT 1,
    note text
);

create table booking (
    bookingID serial PRIMARY KEY,
    blockID serial REFERENCES block (blockID),
    familyID serial REFERENCES family (familyID),
    bookingStart TIMESTAMP,
    bookingEnd TIMESTAMP
);

