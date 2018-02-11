
--add parent to january 1 all 3 blocks room 1
INSERT INTO booking (block_id, user_id)
VALUES (1, 1);
INSERT INTO booking (block_id, user_id)
VALUES (2, 1);
INSERT INTO booking (block_id, user_id)
VALUES (3, 1);

--add teacher to january 2 all 3 blocks room 1
INSERT INTO booking (block_id, user_id)
VALUES (4, 2);
INSERT INTO booking (block_id, user_id)
VALUES (5, 2);
INSERT INTO booking (block_id, user_id)
VALUES (6, 2);

--add admin to january 3 all 3 blocks room 1
INSERT INTO booking (block_id, user_id)
VALUES (7, 3);
INSERT INTO booking (block_id, user_id)
VALUES (8, 3);
INSERT INTO booking (block_id, user_id)
VALUES (9, 3);

--add parent to january 4 all 3 blocks room 1
INSERT INTO booking (block_id, user_id)
VALUES (10, 1);
INSERT INTO booking (block_id, user_id)
VALUES (11, 1);
INSERT INTO booking (block_id, user_id)
VALUES (12, 1);


--add teacher to january 4 all 3 blocks room 1
INSERT INTO booking (block_id, user_id)
VALUES (10, 2);
INSERT INTO booking (block_id, user_id)
VALUES (11, 2);
INSERT INTO booking (block_id, user_id)
VALUES (12, 2);

--add parent to january 5 all 3 blocks room 1
INSERT INTO booking (block_id, user_id)
VALUES (13, 1);
INSERT INTO booking (block_id, user_id)
VALUES (14, 1);
INSERT INTO booking (block_id, user_id)
VALUES (15, 1);

--try to add same parent to january 5 all 3 blocks of room 2
INSERT INTO booking (block_id, user_id)
VALUES (28, 1);
INSERT INTO booking (block_id, user_id)
VALUES (29, 1);
INSERT INTO booking (block_id, user_id)
VALUES (30, 1);
