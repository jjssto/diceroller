package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Game int

const (
	CoC Game = iota + 1
	RezTech
)

type Die struct {
	Eyes   int8
	Result int8
}

type DiceRoll struct {
	Dice   []Die
	Player int
	Action string
	Result int
	Time   time.Time
}

type Player struct {
	Name  string
	Color string
}

type PlayerList struct {
	Id     []int
	Detail map[int]Player
}

type Room struct {
	Id        int
	Game      Game
	Players   PlayerList
	DiceRolls map[int64]DiceRoll
}

func (d *Die) json() string {
	return fmt.Sprintf("{\"E\":\"%d\", \"R\":\"%d\"}", d.Eyes, d.Result)
}

func (d *DiceRoll) json(room int) string {
	player := rooms[room].Players.Detail[d.Player]
	info := fmt.Sprintf("\"P\":\"%s\", \"C\":\"%s\",\"A\":%s, \"T\":\"%d\", \"R\":\"%d\"",
		player.Name, player.Color, d.Action, d.Time.UnixMilli(), d.Result)
	dice := ""
	for _, val := range d.Dice {
		dice += fmt.Sprintf("%s,", val.json())
	}
	return fmt.Sprintf("{%s, \"D\":[%s]}", info, dice)
}

func (r *Room) roll(dice []int8, mod int, player int, action string) (DiceRoll, error) {

	var diceRoll DiceRoll
	diceRoll.Time = time.Now()
	diceRoll.Player = player
	diceRoll.Action = action
	switch r.Game {
	case CoC:
		diceRoll.rollCoC(dice, mod)
	case RezTech:
		diceRoll.rollRezTech(dice, mod)
	}
	if r.DiceRolls == nil {
		r.DiceRolls = make(map[int64]DiceRoll)
	}
	key := diceRoll.Time.UnixNano()
	r.DiceRolls[key] = diceRoll
	return diceRoll, nil
}

func initRand() {
	rand.Seed(time.Now().UnixNano())
}

func (r *DiceRoll) rollCoC(dice []int8, mod int) error {
	switch mod {
	case 2:
		dice = []int8{0, 0, 0, 10}
	case 1:
		dice = []int8{0, 0, 10}
	case 0:
		dice = []int8{0, 10}
	case -1:
		dice = []int8{0, 0, 10}
	case -2:
		dice = []int8{0, 0, 0, 10}
	default:
		return errors.New("Invalid modifier")
	}
	err := r.roll(dice)
	if err != nil {
		return err
	}
	r.evaluateCoC(mod)
	return nil
}

func (r *DiceRoll) evaluateCoC(mod int) {
	basis1 := int(0)
	basis10Arr := make([]int, 0, 3)
	for _, val := range r.Dice {
		if val.Eyes == 0 {
			if val.Result == 10 {
				basis10Arr = append(basis10Arr, 0)
			} else {
				basis10Arr = append(basis10Arr, int(val.Result))
			}
		} else {
			if val.Result == 10 {
				basis1 = 0
			} else {
				basis1 = int(val.Result)
			}
		}
	}

	var basis10 int
	if mod >= 0 {
		basis10 = 10
	} else {
		basis10 = -1
	}
	for _, val := range basis10Arr {
		if mod >= 0 {
			if basis1 == 0 && val == 0 && basis10 == 10 {
				basis10 = 0
			} else if val < basis10 {
				basis10 = val
			}
		} else {
			if basis1 == 0 && val == 0 {
				basis10 = 0
			} else if val > basis10 && (basis10 != 0) {
				basis10 = val
			}
		}
	}
	if basis10 == 0 && basis1 == 0 {
		r.Result = 100
	} else {
		r.Result = basis10*10 + basis1
	}
}

func (r *DiceRoll) rollRezTech(dice []int8, mod int) {
}

func (r *DiceRoll) roll(dice []int8) error {
	for _, val := range dice {
		result, err := roll(val)
		if err != nil {
			return err
		}
		r.Dice = append(r.Dice, result)
	}
	return nil
}

func roll(eyes int8) (Die, error) {
	var max int
	if eyes == 0 {
		max = 9
	} else if eyes > 0 {
		max = int(eyes) - 1
	} else {
		return Die{}, errors.New("Invalid number of eyes")
	}
	return Die{Eyes: eyes, Result: int8(rand.Intn(max) + 1)}, nil
}

func (r *Room) addPlayer(id int, name string, col string) error {
	if r.Players.Detail == nil {
		r.Players.Detail = make(map[int]Player)
	}
	player, ok := r.Players.Detail[id]
	if !ok {
		p := Player{Name: name, Color: col}
		r.Players.Detail[id] = p
		r.Players.Id = append(r.Players.Id, id)
	} else {
		player.Name = name
		player.Color = col
	}
	return nil
}

func addRoom(game Game) (int, error) {
	id := rand.Intn(9998) + 1
	r := Room{Id: id, Game: game}
	rooms[id] = r
	return id, nil
}
