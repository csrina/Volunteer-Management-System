\c caraway
--Update families

INSERT INTO family (family_name, children)
VALUES ('test', 1);
INSERT INTO family (family_name, children)
VALUES ('Robert', 2);
INSERT INTO family (family_name, children)
VALUES ('Name', 1);
INSERT INTO family (family_name, children)
VALUES ('Samename', 10);



--Adds basic users

INSERT INTO users (user_role, username, password, first_name, last_name, email, phone_number, family_id)
VALUES (1, 'parent', '$2a$10$SqeZIWv4nkdfppU8TL7.hO2lwcrPPP7Dg01LXHqBW1NWQNf8Vcf6C', 'parentfirst', 'parentlast', 'email', '123-4567', 1);
INSERT INTO users (user_role, username, password, first_name, last_name, email, phone_number)
VALUES (2, 'teacher', '$2a$10$SqeZIWv4nkdfppU8TL7.hO2lwcrPPP7Dg01LXHqBW1NWQNf8Vcf6C', 'teacherfirst', 'teacherlast', 'email', '123-4567');
INSERT INTO users (user_role, username, password, first_name, last_name, email, phone_number)
VALUES (3, 'admin', '$2a$10$SqeZIWv4nkdfppU8TL7.hO2lwcrPPP7Dg01LXHqBW1NWQNf8Vcf6C', 'adminfirst', 'adminlast', 'email', '123-4567');

--Adds parents to be associated with families

INSERT INTO users (user_role, username, password, first_name, last_name, email, phone_number, family_id)
VALUES(1, 'Will', 'pass', 'William', 'Robert', 'billybob@gmail.com', '132-4365', 2);
INSERT INTO users (user_role, username, password, first_name, last_name, email, phone_number, family_id)
VALUES(1, 'Lope', 'pass', 'Penelope', 'Name', 'penName@gmail.com', '132-4365', 2);
INSERT INTO users (user_role, username, password, first_name, last_name, email, phone_number, family_id)
VALUES(1, 'Susie', 'pass', 'Susie', 'Samename', 'susie53@gmail.com', '132-4365', 3);
INSERT INTO users (user_role, username, password, first_name, last_name, email, phone_number, family_id)
VALUES(1, 'Stevie', 'pass', 'Stevie', 'Samename', 'stevie54@gmail.com', '132-4365', 4);
INSERT INTO users (user_role, username, password, first_name, last_name, email, phone_number, family_id)
VALUES(1, '1234', '1234', 'Bad', 'Username', 'badusername@gmail.com', '132-4365', 4);

--Adds rooms

INSERT INTO room(room_name, teacher_id, children, room_num)
VALUES('blue', 2, 10, '5-212');
INSERT INTO room(room_name, teacher_id, children, room_num)
VALUES('red', 2, 10, '5-213');
INSERT INTO room(room_name, teacher_id, children, room_num)
VALUES('green', 2, 10, '5-214');
INSERT INTO room(room_name, teacher_id, children, room_num)
VALUES('purple', 2, 10, '5-215');

--Adds time blocks room 1 (january 1)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-01 08:00:00', '2018-04-01 11:00:00', 1, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-01 12:00:00', '2018-04-01 14:00:00', 1, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-01 15:00:00', '2018-04-01 17:00:00', 1, 1, 'afternoon block!');

--Adds time blocks room 1 (january 2)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-02 08:00:00', '2018-04-02 11:00:00', 1, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-02 12:00:00', '2018-04-02 14:00:00', 1, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-02 15:00:00', '2018-04-02 17:00:00', 1, 1, 'afternoon block!');

--Adds time blocks room 1 (january 3)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-03 08:00:00', '2018-04-03 11:00:00', 1, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-03 12:00:00', '2018-04-03 14:00:00', 1, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-03 15:00:00', '2018-04-03 17:00:00', 1, 1, 'afternoon block!');

--Adds time blocks room 1 (january 4)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-04 08:00:00', '2018-04-04 11:00:00', 1, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-04 12:00:00', '2018-04-04 14:00:00', 1, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-04 15:00:00', '2018-04-04 17:00:00', 1, 1, 'afternoon block!');


--Adds time blocks room 1 (january 5)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-05 08:00:00', '2018-04-05 11:00:00', 1, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-05 12:00:00', '2018-04-05 14:00:00', 1, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-05 15:00:00', '2018-04-05 17:00:00', 1, 1, 'afternoon block!');


-----------------------------------------------------------------------------
-----------------------------------------------------------------------------

--Adds time blocks room 2 (january 1)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-01 08:00:00', '2018-04-01 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-01 12:00:00', '2018-04-01 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-01 15:00:00', '2018-04-01 17:00:00', 2, 1, 'afternoon block!');

--Adds time blocks room 2 (january 2)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-02 08:00:00', '2018-04-02 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-02 12:00:00', '2018-04-02 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-02 15:00:00', '2018-04-02 17:00:00', 2, 1, 'afternoon block!');

--Adds time blocks room 2 (january 3)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-03 08:00:00', '2018-04-03 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-03 12:00:00', '2018-04-03 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-03 15:00:00', '2018-04-03 17:00:00', 2, 1, 'afternoon block!');

--Adds time blocks room 2 (january 4)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-04 08:00:00', '2018-04-04 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-04 12:00:00', '2018-04-04 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-04 15:00:00', '2018-04-04 17:00:00', 2, 1, 'afternoon block!');

--Adds time blocks room 2 (january 5)

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-05 08:00:00', '2018-04-05 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-05 12:00:00', '2018-04-05 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-05 15:00:00', '2018-04-05 17:00:00', 2, 1, 'afternoon block!');

-------------------------------------------------------------------------------
-------------------------------------------------------------------------------

--february times

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-19 08:00:00', '2018-04-19 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-19 12:00:00', '2018-04-19 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-19 15:00:00', '2018-04-19 17:00:00', 2, 1, 'afternoon block!');


INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-20 08:00:00', '2018-04-20 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-20 12:00:00', '2018-04-20 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-20 15:00:00', '2018-04-20 17:00:00', 2, 1, 'afternoon block!');

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-21 08:00:00', '2018-04-21 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-21 12:00:00', '2018-04-21 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-21 15:00:00', '2018-04-21 17:00:00', 2, 1, 'afternoon block!');

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-22 08:00:00', '2018-04-22 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-22 12:00:00', '2018-04-22 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-22 15:00:00', '2018-04-22 17:00:00', 2, 1, 'afternoon block!');

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-23 08:00:00', '2018-04-23 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-23 12:00:00', '2018-04-23 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-23 15:00:00', '2018-04-23 17:00:00', 2, 1, 'afternoon block!');


INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-26 08:00:00', '2018-04-26 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-26 12:00:00', '2018-04-26 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-26 15:00:00', '2018-04-26 17:00:00', 2, 1, 'afternoon block!');

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-07 08:00:00', '2018-04-07 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-08 12:00:00', '2018-04-08 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-09 15:00:00', '2018-04-09 17:00:00', 2, 1, 'afternoon block!');

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-07 08:00:00', '2018-04-07 11:00:00', 1, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-08 12:00:00', '2018-04-08 14:00:00', 1, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-09 15:00:00', '2018-04-09 17:00:00', 1, 1, 'afternoon block!');

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-14 08:00:00', '2018-04-14 11:00:00', 1, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-13 12:00:00', '2018-04-13 14:00:00', 1, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-12 15:00:00', '2018-04-12 17:00:00', 1, 1, 'afternoon block!');

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-14 08:00:00', '2018-04-14 11:00:00', 2, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-13 12:00:00', '2018-04-13 14:00:00', 2, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-12 15:00:00', '2018-04-12 17:00:00', 2, 1, 'afternoon block!');

INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-14 08:00:00', '2018-04-14 11:00:00', 3, 1, 'morning block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-13 12:00:00', '2018-04-13 14:00:00', 3, 1, 'noon block!');
INSERT INTO time_block(block_start, block_end, room_id, modifier, note)
VALUES('2018-04-12 15:00:00', '2018-04-12 17:00:00', 3, 1, 'afternoon block!');
