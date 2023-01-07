-- create database diceroller_test_db;
-- create user db_user@localhost identified WITH mysql_native_password by '123';
-- grant all on diceroller_test_db.* to db_user@localhost;

use diceroller_test_db;

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
drop procedure if exists clean_up;

DELIMITER $$

create function get_rolls_json(
    room_id integer,
    token integer,
    last_roll integer
)
returns text deterministic
begin    
    if last_roll >= (
        select last_nr
        from last_roll_nbr where last_roll_nbr.room_id = room_id
    ) then
        return '';
    end if;

	select
		group_concat(
			concat(
                '\"', sub.nr,'\":[',
                    '{',
						'\"P\":\"', ifnull(sub.na, ''), '\",',
                        '\"C\":\"', ifnull(sub.co, ''), '\",',
                        '\"A\":\"', ifnull(sub.ac, ''), '\",',
                        '\"T\":\"', ifnull(sub.cr, ''), '\",',
                        '\"R\":\"', ifnull(sub.re, ''), '\",',
                        '\"D\":[', ifnull(sub.dice, ''), ']',
                    '}, ',
                '\"', sub.ow, '\"]')
			order by sub.nr
			separator ', '
		)
	into @json_str
	from (
		select
			roll.roll_nbr as nr,
			roll.chr_action as ac,
			roll.result as re,
			unix_timestamp(roll.created) as cr,
			group_concat( 
				concat('{\"E\": \"', die.eyes, '\", \"R": \"', die.result, '\"}')
				order by die.eyes desc 
				separator ', '
			) as dice,
			chr.chr_name as na,
			chr.chr_color as co,
			if(user_token.token = token, 1, 0) as ow
				
			from 
				room 
				join game on room.game_id = game.id
				join chr on room.id = chr.room_id  
				join roll on chr.id = roll.chr_id
				join die on roll.id = die.roll_id
				left join user_token on chr.user_token_id = user_token.id
			where 
				room.id = room_id and ifnull(roll.roll_nbr, 1) > ifnull(last_roll, 0)
			group by roll.id
			order by roll.roll_nbr
	) as sub;
    
    return @json_str;
end;

$$

DELIMITER $$
create procedure get_room(room_id int, token int)
begin
    select id into @user_id from user_token where user_token.token = token;
    select 
        game.game,
        room.id,
        room.owner_id,
        room.room_name,
        room.color
    into @game, @room_id, @owner_id, @name, @color
    from
		room
        join game on room.game_id = game.id
    where room.id = room_id;
    if @owner_id is null then
        update room
        set owner_id = @user_id
        where id = @room_id;
    end if; 
    select 
        ifnull(@game, ''), 
        if(@user_id = ifnull(@owner_id, @user_id), 1, 0),
        ifnull(@name, ''), 
        ifnull(@color, '');
end;
$$

DELIMITER $$
create function get_roll_nbr (arg_room_id int)
returns integer deterministic
begin
	select last_nr into @roll_nbr from last_roll_nbr where room_id = arg_room_id;
    if ifnull(@roll_nbr, 0) = 0 then
		insert into last_roll_nbr(room_id, last_nr) values (arg_room_id, 1);
        return 1;
    else
		update last_roll_nbr set last_nr = @roll_nbr + 1 where room_id = arg_room_id;
		return @roll_nbr + 1;
    end if;
end;
$$


DELIMITER $$
create procedure insert_roll (
    in arg_chr_id int, 
    in arg_chr_action varchar(16), 
    in arg_result int
) 
begin 
    insert into roll(roll_nbr, chr_id, chr_action, result)
    select
        get_roll_nbr(chr.room_id), 
        arg_chr_id, 
        arg_chr_action,
        arg_result
    from chr 
    where chr.id = arg_chr_id;
    select last_insert_id();
end;
$$

DELIMITER $$
create procedure create_user_token (
    in arg_token int
) 
begin 
    declare ret_token int;
    declare ret_id int;
    if (select count(*) from user_token where token = arg_token) > 0 then
        select id, token into ret_id, ret_token from user_token where token = arg_token;
    else
        while_loop: while 1 = 1 do
            set ret_token = floor( 100000 + rand() * 999999);
            if (select count(*) from user_token where token = ret_token) = 0 then
                leave while_loop;
            end if;
        end while while_loop;
        insert into user_token(token) values (ret_token);
        set ret_id = last_insert_id();
    end if;
    select ret_token, ret_id;
end;
$$


DELIMITER $$
create procedure create_room (in game int) 
begin
    declare id_found int;
    declare new_id int;
    set id_found = 0;
    while id_found = 0 do 
        select floor(1 + rand() * 999999) into new_id;
        if 0 = (
            select count(*) from room where id = new_id
        ) then
            set id_found = 1;
        end if;
    end while;
    insert into room (id, game_id) values (new_id, game);
    select new_id;
end;
$$

DELIMITER $$
create procedure change_room_settings (
    in arg_room_id int,
    in arg_token int,
    in arg_name char(64),
    in arg_col char(12)
) 
begin
    declare ret int;
    if (
        select count(*)
        from room join user_token on room.owner_id = user_token.id
        where room.id = arg_room_id and user_token.token = arg_token
    ) > 0 then
        update room
            set room.room_name = if(trim(arg_name) <> '', arg_name, room.room_name),
            room.color = if(trim(arg_col) <> '', arg_col, room.color)
        where room.id = arg_room_id;
        set ret = 1;
    else
        set ret = -1;
    end if;
    select ret;
end;
$$




DELIMITER $$
create procedure remove_old_rooms (
    in nbr_of_days int
) 
begin
    declare compar timestamp;
    declare zero timestamp;
    set compar = curdate() - interval nbr_of_days day;
    set zero = date('2000-01-01');
    delete room 
    from
        room
        left join (
			select
				chr.room_id as room_id,
                max(roll.created) as last_roll
            from
				chr
				left join roll on chr.id = roll.chr_id
			group by room_id
		) as sub on room.id = sub.room_id
    where
       room.created < compar and (sub.last_roll is null or sub.last_roll < compar);
end;
$$

DELIMITER $$
create procedure remove_old_tokens () 
begin
    delete user_token
    from
        user_token
        left join chr on user_token.id = chr.user_token_id
        left join room on user_token.id = room.owner_id
    where
		chr.id is null and room.id is null;
end;
$$

DELIMITER $$
create procedure clean_up( in nbr_of_days int)
begin
	call remove_old_rooms(nbr_of_days);
    call remove_old_tokens();
end;
$$