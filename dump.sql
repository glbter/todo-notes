insert into users
    (id, name, time_zone, password_hash)
values
    (1, 'user1', 'UTC', 'user'),
    (2, 'user2', 'UTC', 'user'),
    (3, 'user3', 'UTC', 'user'),
    (4, 'user4', 'UTC', 'user'),
    (5, 'user5', 'UTC', 'user'),
    (6, 'user6', 'UTC', 'user'),
    (7, 'user7', 'UTC', 'user'),
    (8, 'user8', 'UTC', 'user'),
    (9, 'user9', 'UTC', 'user');


insert into notes
    (id, user_id, title, text, date, is_finished)
values
    (1, 1, 'title', 'text', '2021-10-13 10:20:08.392115', false),
    (2, 1, 'title', 'text', '2021-10-12 10:20:08.392115', false),
    (3, 1, 'title', 'text', '2021-10-13 10:20:08.392115', false),
    (4, 1, 'title', 'text', '2021-10-13 10:20:08.392115', true),
    (5, 1, 'title', 'text', '2021-10-12 10:20:08.392115', true),
    (6, 1, 'title', 'text', '2021-10-13 10:20:08.392115', true),

    (7, 2, 'title', 'text', '2021-10-12 10:20:08.392115', true),
    (8, 2, 'title', 'text', '2021-10-13 10:20:08.392115', true),
    (9, 2, 'title', 'text', '2021-10-12 10:20:08.392115', false);

-- dates 13.10.21 and 12.10.21, user_ids 1 and 2
