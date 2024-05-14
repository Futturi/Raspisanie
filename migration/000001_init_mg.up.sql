CREATE TABLE pairs
(
    id bigserial primary key,
    groups varchar(50),
    week int,
    day int,
    pairs json
);
CREATE TABLE users
(
    id bigserial primary key,
    email text not null unique ,
    password text not null ,
    name text not null unique ,
    groupa text not null
);
CREATE TABLE rooms
(
    id bigserial primary key,
    name text not null unique,
    groupa text not null
);

CREATE TABLE messages
(
    id bigserial primary key,
    text text not null,
    creator text not null,
    date int not null
);

CREATE TABLE rooms_messages
(
    id bigserial primary key,
    room_id bigint references rooms(id) on delete cascade,
    message_id bigint references messages(id) on delete cascade
);

CREATE  TABLE rooms_users
(
    id bigserial primary key,
    room_id bigint references rooms(id) on delete cascade,
    user_id bigint references users(id) on delete cascade
);