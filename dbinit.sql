create table room (
    roomID serial PRIMARY KEY,
    name text,
    teacher text,
    roomNum text
);

/* what else do we want in family? */
create table family (
    familyID serial PRIMARY KEY,
    
);

create table user (
    userID serial PRIMARY KEY,
    familyID serial REFERENCES family (familyID),
    firstName text,
    lastName text,
    phoneNumber text, /* text for now, probably a better format available */
)

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

