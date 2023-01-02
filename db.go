package main

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type DB struct {
	Db         *sql.DB
	Cfg        mysql.Config
	Configured bool
}

func (db *DB) config() {
	db.Cfg = mysql.Config{
		User:   "",
		Passwd: "",
		Net:    "tcp",
		Addr:   "",
		DBName: "",
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
	userToken string,
	charName string,
	charColor string,
	mod int,
	dice []int8,
) error {
	row := db.Db.QueryRow(
		`select 
			game.game,
			isnull(player.id, -1),
			isnull(player.char_name, ''),
			isnull(player.char_color, '')
		from
			room
			left join player
			left join user
			left join game
		where room.id = ? and user.token = ?`,
		roomId, userToken,
	)
	var gameStr string
	var charId int
	var dbCharName string
	var dbCharColor string
	err := row.Scan(&gameStr, &charId, &dbCharName, &dbCharColor)
	if err != nil {
		return err
	}
	if charId == -1 {
		charId, err = db.addChar(userToken, charName, charColor, roomId)
		if err != nil {
			return err
		}
	} else {
		db.updateChar(charId, charName, charColor, dbCharName, dbCharColor)
	}
	var roll DiceRoll
	switch gameStr {
	case "CoC":
		err = roll.rollCoC(dice, mod)
	case "RezTech":
		err = roll.rollRezTech(dice, mod)
	case "General":
		err = roll.rollGeneral(dice, mod)
	}
	if err != nil {
		return err
	}
	if err = db.saveRoll(&roll, charId); err != nil {
		return err
	}
	return nil
}

func (db *DB) saveRoll(roll *DiceRoll, charId int) error {
	result, err := db.Db.Exec(
		`insert into roll(chr_id, chr_action, result)
		values (?, ?, ?)
		`, charId, roll.Action, roll.Result,
	)
	if err != nil {
		return err
	}
	rollId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	for eyes, result := range roll.Dice {
		_, err = db.Db.Exec(
			`insert into die(roll_id, eyes, result)
			values(?, ?, ?)`,
			rollId, eyes, result,
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
) {
	if charName == dbCharName && charColor == dbCharColor {
		return
	}
	if charName == "" {
		charName = dbCharName
	}
	if charColor == "" {
		charColor = dbCharColor
	}
	db.Db.Exec(
		`update chr 
			set chr_name = ?, chr_color = ?
		where chr_id = ?`, charName, charColor, charId,
	)
}

func (db *DB) createRoom(game Game) int {
	return 0
}

func (db *DB) addUser(token string) (int, error) {
	return -1, nil
}

func (db *DB) addChar(
	token string,
	charName string,
	charColor string,
	roomId int,
) (int, error) {
	row := db.Db.QueryRow(
		`select 
			user.id,
			isnull(chr.id, -1)
		from
			user
			left join chr on 
				user.id = chr.user_id
				and chr.room_id = ?
		where user.token = ?`,
		roomId, token,
	)
	var userId int
	var charId int
	var err error
	err = row.Scan(&userId, &charId)
	if err != nil {
		if err == sql.ErrNoRows {
			userId, err = db.addUser(token)
		} else {
			return -1, err
		}
	}
	if charId == -1 {
		// insert char
		result, err := db.Db.Exec(
			`insert into chr(room_id, user_id, chr_name, chr_color)
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

func (db *DB) getRolls(roomId int, token string, lastRoll int) string {
	rows, err := db.Db.Query(
		`select
			game.game	
			if user.token = ? then 1 else 0 end if
		from
			roll
			join die on roll.id = die.roll_id
			join chr on roll.chr_id = chr.id
			join user on chr.user_id = user.id
			join room on roll.room_id = room.id
			join game on room.game_id = game.id
		where roll.room_id = ? and roll.id > ?		
		order by roll.id, die.eyes`, token, roomId, lastRoll,
	)
	if err != nil {
		return ""
	}
	for rows.Next() {
		err = rows.Scan()
	}
	return ""
}
