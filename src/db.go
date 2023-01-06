package main

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
)

type DB struct {
	Db         *sql.DB
	Cfg        mysql.Config
	Configured bool
}

func (db *DB) config() {
	db.Cfg = mysql.Config{
		User:                 globConfig.DBUser,
		Passwd:               globConfig.DBPassword,
		Net:                  globConfig.DBNet,
		Addr:                 globConfig.DBAddress,
		DBName:               globConfig.DBName,
		AllowNativePasswords: true,
	}
	db.Configured = true
}

func (db *DB) connect(close bool) error {
	var err error
	if !db.Configured {
		db.config()
	}
	db.Db, err = sql.Open("mysql", db.Cfg.FormatDSN())
	if err != nil {
		return err
	} else {
		if close {
			defer db.Db.Close()
		}
		return nil
	}
}

func (db *DB) close() {
	defer db.Db.Close()
}

func (db *DB) roll(
	roomId int,
	userToken int,
	charName string,
	charColor string,
	charAction string,
	mod int,
	dice []int8,
) error {
	row := db.Db.QueryRow(
		`select 
			game.game,
			ifnull(chr.id, -1),
			ifnull(chr.chr_name, ''),
			ifnull(chr.chr_color, '')
		from
			room
			left join chr on room.id = chr.room_id 
			left join user_token on chr.user_token_id = user_token.id 
			left join game on room.game_id = game.id
		where room.id = ? and user_token.token = ?`,
		roomId, userToken,
	)
	var gameStr string
	var charId int
	var dbCharName string
	var dbCharColor string
	err := row.Scan(&gameStr, &charId, &dbCharName, &dbCharColor)
	if err == sql.ErrNoRows {
		charId = -1
	} else if err != nil {
		return err
	}
	if charId == -1 {
		charId, err = db.addChar(userToken, charName, charColor, roomId)
		if err != nil {
			return err
		}
		row := db.Db.QueryRow(
			`select game.game 
			from room join game on room.game_id = game.id
			where room.id = ?`, roomId,
		)
		row.Scan(&gameStr)
	} else {
		db.updateChar(charId, charName, charColor, dbCharName, dbCharColor)
	}
	roll := DiceRoll{
		Action: charAction,
		Player: charId,
	}
	switch gameStr {
	case "CoC":
		err = roll.rollCoC(dice, mod)
	case "RezTech":
		err = roll.rollRezTech(dice, mod)
	case "General":
		err = roll.rollGeneral(dice, mod)
	default:
		return errors.New("unknown game")
	}
	if err != nil {
		return err
	}
	if err = db.saveRoll(&roll); err != nil {
		return err
	}
	return nil
}

func (db *DB) saveRoll(roll *DiceRoll) error {
	row := db.Db.QueryRow(
		`call insert_roll(?, ?, ?)
		`, roll.Player, roll.Action, roll.Result,
	)
	var rollId int
	if err := row.Scan(&rollId); err != nil {
		return err
	}
	for i := range roll.Dice {
		d := roll.Dice[i]
		_, err := db.Db.Exec(
			`insert into die(roll_id, eyes, result)
			values(?, ?, ?)`,
			rollId, d.Eyes, d.Result,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) updateChar(
	charId int,
	charName string,
	charColor string,
	dbCharName string,
	dbCharColor string,
) error {
	if charName == dbCharName && charColor == dbCharColor {
		return nil
	}
	if charName == "" {
		charName = dbCharName
	}
	if charColor == "" {
		charColor = dbCharColor
	}
	if _, err := db.Db.Exec(
		`update chr 
			set chr_name = ?, chr_color = ?
		where id = ?`, charName, charColor, charId,
	); err != nil {
		return err
	}
	return nil
}

func (db *DB) createRoom(game Game) (int, error) {
	var game_id int
	switch game {
	case General:
		game_id = 1
	case CoC:
		game_id = 2
	case RezTech:
		game_id = 3
	default:
		game_id = 1
	}
	result := db.Db.QueryRow(`call create_room(?)`, game_id)

	var roomId int
	err := result.Scan(&roomId)
	if err != nil {
		return -1, err
	}
	return roomId, nil
}

func (db *DB) createToken(oldToken int) (int, int, error) {
	row := db.Db.QueryRow(`call create_user_token(?)`, oldToken)
	var newToken int
	var id int
	if err := row.Scan(&newToken, &id); err != nil {
		return 0, 0, err
	}
	return newToken, id, nil
}

func (db *DB) getRoom(roomId int, token int) (Room, error) {
	row := db.Db.QueryRow(`call get_room(?, ?)`, roomId, token)
	var gameStr string
	var isOwnerInt int
	room := Room{Id: roomId}
	if err := row.Scan(&gameStr, &isOwnerInt, &room.Name, &room.Color); err != nil {
		return room, err
	}
	switch gameStr {
	case "CoC":
		room.Game = CoC
	case "RezTech":
		room.Game = RezTech
	case "":
		return room, errors.New("room does not exist")
	default:
		room.Game = General
	}
	room.IsOwner = (isOwnerInt == 1)
	return room, nil
}

func (db *DB) addChar(
	token int,
	charName string,
	charColor string,
	roomId int,
) (int, error) {
	row := db.Db.QueryRow(
		`select 
			user_token.id,
			ifnull(chr.id, 0)
		from
			user_token
			left join chr on 
				user_token.id = chr.user_token_id
				and chr.room_id = ?
		where user_token.token = ?`,
		roomId, token,
	)
	var userId int
	var charId int
	if err := row.Scan(&userId, &charId); err != nil {
		return -1, err
	}
	//if err != nil {
	// if err == sql.ErrNoRows {
	// 	_, userId, err = db.createToken(token)
	// 	if err != nil {
	// 		return -1, err
	// 	}
	// } else {
	//	return -1, err
	//}
	//}
	if charId == 0 {
		// insert char
		result, err := db.Db.Exec(
			`insert into chr(room_id, user_token_id, chr_name, chr_color)
			values (?, ?, ?, ?)`, roomId, userId, charName, charColor,
		)
		if err != nil {
			return -1, err
		}
		charId64, err := result.LastInsertId()
		if err != nil {
			return -1, err
		}
		charId = int(charId64)
	}
	return charId, nil
}

func (db *DB) getRolls(roomId int, token int, lastRoll int) string {

	row := db.Db.QueryRow(`select get_rolls_json(?,?,?)`, roomId, token, lastRoll)

	var str string
	err := row.Scan(&str)
	if err != nil {
		str = ""
	}
	return "{" + str + "}"
}

func (db *DB) changeRoomSettings(
	roomId int, userToken int, name string, color string,
) error {
	var ret int
	row := db.Db.QueryRow(
		`call change_room_settings(?,?,?,?)`,
		roomId, userToken, name, color,
	)
	if err := row.Scan(&ret); err != nil {
		return err
	}
	if ret > 0 {
		return nil
	} else {
		return errors.New("action is not allowed")
	}

}

func (db *DB) getOwnRooms(userToken int) ([]Room, error) {
	rows, err := db.Db.Query(
		`select distinct
			room.id, 
			ifnull(room.room_name, ''), 
			ifnull(room.color, ''), 
			ifnull(unix_timestamp(room.created), 0),
			ifnull(game.game, ''),
			ifnull(unix_timestamp(max(roll.created)), 0)
		from 
			room 
			join user_token on room.owner_id = user_token.id
			join game on room.game_id = game.id
			left join chr on room.id = chr.room_id
			left join roll on chr.id = roll.chr_id
		where user_token.token = ?
		group by room.id`, userToken,
	)
	if err != nil {
		return nil, err
	}
	var roomId int
	var roomName string
	var roomColor string
	var created int64
	var lastActivity int64
	var gameStr string
	ret := make([]Room, 0)
	for rows.Next() {
		err = rows.Scan(
			&roomId, &roomName, &roomColor, &created, &gameStr, &lastActivity,
		)
		if err != nil {
			return nil, err
		}
		ret = append(ret, Room{
			Id:           roomId,
			Name:         roomName,
			Color:        roomColor,
			Created:      time.Unix(created, 0),
			LastActivity: time.Unix(lastActivity, 0),
			GameStr:      gameStr,
		})
	}
	return ret, nil
}

func (db *DB) getAllRooms(userToken int) ([]Room, error) {
	rows, err := db.Db.Query(
		`select distinct
			room.id, 
			ifnull(room.room_name, ''), 
			ifnull(room.color, ''), 
			ifnull(unix_timestamp(room.created), 0),
			ifnull(game.game, ''),
			ifnull(unix_timestamp(max(roll.created)), 0)
		from 
			room 
			join chr on room.id = chr.room_id
			join user_token on chr.user_token_id = user_token.id
			join game on room.game_id = game.id
			left join roll on chr.id = roll.chr_id
		where user_token.token = ?
		group by room.id`, userToken,
	)
	if err != nil {
		return nil, err
	}
	var roomId int
	var roomName string
	var roomColor string
	var created int64
	var lastActivity int64
	var gameStr string
	ret := make([]Room, 0)
	for rows.Next() {
		err = rows.Scan(
			&roomId, &roomName, &roomColor, &created, &gameStr, &lastActivity,
		)
		if err != nil {
			return nil, err
		}
		ret = append(ret, Room{
			Id:           roomId,
			Name:         roomName,
			Color:        roomColor,
			Created:      time.Unix(created, 0),
			LastActivity: time.Unix(lastActivity, 0),
			GameStr:      gameStr,
		})
	}
	return ret, nil
}

func (db *DB) deleteRoom(userToken int, roomId int64) (int, int64, error) {
	result, err := db.Db.Exec(
		`delete room 
		from room join user_token on room.owner_id = user_token.id
		where room.id = ? and user_token.token = ?`, roomId, userToken,
	)
	if err != nil {
		return -1, -1, err
	}
	nbr, err := result.RowsAffected()
	if err != nil {
		return -1, -1, err
	} else {
		return int(nbr), roomId, nil
	}
}
