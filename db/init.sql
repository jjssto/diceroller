-- create database diceroller_test_db;
-- create user db_user@localhost identified WITH mysql_native_password by '123';
-- grant all on diceroller_test_db.* to db_user@localhost;

use diceroller_test_db;

alter table room
drop foreign key room_ibfk_2;
drop table if exists last_roll_nbr;
drop table if exists die;
drop table if exists roll;
drop table if exists chr;
drop table if exists user_token;
drop table if exists room;
drop table if exists game;
drop function if exists get_rolls_json;
drop procedure if exists get_room;
drop function if exists get_roll_nbr;
drop procedure if exists create_room;
drop procedure if exists insert_roll;
drop procedure if exists create_user_token;
drop procedure if exists change_room_settings; 
drop procedure if exists removeOldRooms;
drop procedure if exists remove_old_rooms;
drop procedure if exists remove_old_tokens;

create table game (
    id tinyint primary key,
    game varchar(16)
);

insert into game
values 
(1, "General"),
(2, "CoC"),
(3, "RezTech");

create table room (
    id int primary key,
    game_id tinyint not null,
    owner_id int,
    room_name varchar(64),
    color char(7),
    created timestamp default current_timestamp, 
    foreign key fk_room_game (game_id)
        references game (id)
);


create table user_token (
    id int primary key auto_increment,
    token int,
    constraint uc_user_token unique(token)
);

create table chr (
    id int primary key auto_increment,
    room_id int not null,
    user_token_id int not null,
    chr_name varchar(16),
    chr_color varchar(7),
    foreign key fk_chr_room (room_id) references room (id)
        on delete cascade,
    foreign key fk_chr_user (user_token_id) references user_token (id)
);

alter table room
add foreign key fk_room_owner (owner_id) references user_token (id);

create index index_character on chr (room_id, user_token_id);


create table roll (
    id int primary key auto_increment,
    roll_nbr int,
    chr_id int not null,
    result int,
    chr_action varchar(16),
    created timestamp default current_timestamp,
    foreign key fk_roll_chr (chr_id) references chr (id)
        on delete cascade
);
create index index_roll_nbr on roll (roll_nbr, chr_id);

create table die (
    id int primary key auto_increment,
    roll_id int not null,
    result tinyint,
    eyes tinyint not null,
    foreign key fk_die_roll (roll_id) references roll (id)
        on delete cascade
);

create table last_roll_nbr (
	room_id int primary key,
    last_nr int,
    foreign key fk_next_roll_nbr (room_id) references room(id)
		on delete cascade
);