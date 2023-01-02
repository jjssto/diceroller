create table game (
    id tinyint primary key,
    game varchar(16)
);

insert into game
values 
(0, "General"),
(1, "CoC"),
(2, "RezTech");

create table room (
    id int primary key,
    game tinyint not null,
    created timestamp, 
    foreign key fk_room_game (game)
        references game (id);
);

create table user {
    id int primary key auto_increment,
    token long varchar
};

create table chr {
    id int primary key auto_increment,
    room_id int not null,
    user_id int not null,
    chr_name varchar(16),
    chr_color varchar(7),
    foreign kye fk_chr_room (room_id) references room (id)
        on delete cascade,
    foreign kye fk_chr_user (user_id) references user (id)
        
};
create index index_character on chr (room_id, user_id);


create table roll (
    id int primary key auto_increment,
    chr_id int not null,
    result int,
    chr_action varchar(16),
    foreign key fk_roll_chr (chr_id) references chr (id)
        on delete cascade;
);


create table die {
    id int primary key auto_increment,
    roll_id int not null,
    result tinyint,
    eyes tinyint not null,
    foreign key fk_die_roll (roll_id) references roll (id)
        on delete cascade;
}