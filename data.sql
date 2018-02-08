
--Adds users

INSERT INTO users (user_role, username, password, first_name, last_name, email, phone_number)
VALUES (1, 'parent', 'pass', 'parentfirst', 'parentlast', 'email', '123-4567');
INSERT INTO users (user_role, username, password, first_name, last_name, email, phone_number)
VALUES (2, 'teacher', 'pass', 'teacherfirst', 'teacherlast', 'email', '123-4567');
INSERT INTO users (user_role, username, password, first_name, last_name, email, phone_number)
VALUES (3, 'admin', 'pass', 'adminfirst', 'adminlast', 'email', '123-4567');

--Adds rooms

INSERT INTO room(room_name, teacher_id, children, room_num)
VALUES('blue', 2, 10, '5-212');
INSERT INTO room(room_name, teacher_id, children, room_num)
VALUES('red', 2, 10, '5-213');
INSERT INTO room(room_name, teacher_id, children, room_num)
VALUES('green', 2, 10, '5-214');

--Adds time blocks room 1 (january 1)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-01 08:00:00', '2018-01-01 11:00:00', 1, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-01 12:00:00', '2018-01-01 14:00:00', 1, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-01 15:00:00', '2018-01-01 17:00:00', 1, 1, 'afternoon block!');

--Adds time blocks room 1 (january 2)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-02 08:00:00', '2018-01-02 11:00:00', 1, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-02 12:00:00', '2018-01-02 14:00:00', 1, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-02 15:00:00', '2018-01-02 17:00:00', 1, 1, 'afternoon block!');

--Adds time blocks room 1 (january 3)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-03 08:00:00', '2018-01-03 11:00:00', 1, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-03 12:00:00', '2018-01-03 14:00:00', 1, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-03 15:00:00', '2018-01-03 17:00:00', 1, 1, 'afternoon block!');

--Adds time blocks room 1 (january 4)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-04 08:00:00', '2018-01-04 11:00:00', 1, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-04 12:00:00', '2018-01-04 14:00:00', 1, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-04 15:00:00', '2018-01-04 17:00:00', 1, 1, 'afternoon block!');


--Adds time blocks room 1 (january 5)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-05 08:00:00', '2018-01-05 11:00:00', 1, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-05 12:00:00', '2018-01-05 14:00:00', 1, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-05 15:00:00', '2018-01-05 17:00:00', 1, 1, 'afternoon block!');


-----------------------------------------------------------------------------
-----------------------------------------------------------------------------

--Adds time blocks room 2 (january 1)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-01 08:00:00', '2018-01-01 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-01 12:00:00', '2018-01-01 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-01 15:00:00', '2018-01-01 17:00:00', 2, 1, 'afternoon block!');

--Adds time blocks room 2 (january 2)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-02 08:00:00', '2018-01-02 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-02 12:00:00', '2018-01-02 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-02 15:00:00', '2018-01-02 17:00:00', 2, 1, 'afternoon block!');

--Adds time blocks room 2 (january 3)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-03 08:00:00', '2018-01-03 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-03 12:00:00', '2018-01-03 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-03 15:00:00', '2018-01-03 17:00:00', 2, 1, 'afternoon block!');

--Adds time blocks room 2 (january 4)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-04 08:00:00', '2018-01-04 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-04 12:00:00', '2018-01-04 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-04 15:00:00', '2018-01-04 17:00:00', 2, 1, 'afternoon block!');

--Adds time blocks room 2 (january 5)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-05 08:00:00', '2018-01-05 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-05 12:00:00', '2018-01-05 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-01-05 15:00:00', '2018-01-05 17:00:00', 2, 1, 'afternoon block!');

-------------------------------------------------------------------------------
-------------------------------------------------------------------------------
