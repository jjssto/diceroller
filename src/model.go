// This code is licensed under the MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
	General
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

type Char struct {
	Name  string
	Color string
}

type CharList struct {
	Id     []int
	Detail map[int]Char
}

type Room struct {
	Id           int
	Game         Game
	Players      CharList
	Color        string
	OwnerId      int
	IsOwner      bool
	Name         string
	DiceRolls    []DiceRoll
	Created      time.Time
	LastActivity time.Time
	GameStr      string
}

func (d *Die) json() string {
	return fmt.Sprintf("{\"E\":\"%d\", \"R\":\"%d\"}", d.Eyes, d.Result)
}

//func (d *DiceRoll) json(room int, loc *time.Location) string {
//	player := globRooms[room].Players.Detail[d.Player]
//	timeStamp := d.Time.Local().In(loc).Format("15:04:05")
//	info := fmt.Sprintf("\"P\":\"%s\", \"C\":\"%s\",\"A\":\"%s\", \"T\":\"%s\", \"R\":\"%d\"",
//		player.Name, player.Color, d.Action, timeStamp, d.Result)
//	dice := ""
//	first := true
//	for _, val := range d.Dice {
//		if first {
//			first = false
//		} else {
//			dice += ","
//		}
//		dice += val.json()
//	}
//	return fmt.Sprintf("{%s, \"D\":[%s]}", info, dice)
//}

func (r *Room) roll(dice []int8, mod int, player int, action string) (DiceRoll, error) {

	var diceRoll DiceRoll
	diceRoll.Time = time.Now()
	diceRoll.Player = player
	diceRoll.Action = action
	switch r.Game {
	case CoC:
		if diceRoll.rollCoC(dice, mod) != nil {
			return diceRoll, errors.New("no dice")
		}
	case RezTech:
		if diceRoll.rollRezTech(dice, mod) != nil {
			return diceRoll, errors.New("no dice")
		}
	case General:
		if diceRoll.rollGeneral(dice, mod) != nil {
			return diceRoll, errors.New("no dice")
		}
	}
	if r.DiceRolls == nil {
		r.DiceRolls = make([]DiceRoll, 0, 1000)
	}
	r.DiceRolls = append(r.DiceRolls, diceRoll)
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
		return errors.New("invalid modifier")
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
	for indx, val := range r.Dice {
		if val.Eyes == 0 {
			if val.Result == 10 {
				r.Dice[indx].Result = 0
				basis10Arr = append(basis10Arr, 0)
			} else {
				basis10Arr = append(basis10Arr, int(val.Result))
			}
		} else {
			if val.Result == 10 {
				r.Dice[indx].Result = 0
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

func (r *DiceRoll) evaluateRezTech() {
	result := 0
	for _, val := range r.Dice {
		if val.Result >= 11 {
			result += 2
		} else if val.Result >= 5 {
			result += 1
		}
	}
	r.Result = result
}

func (r *DiceRoll) rollRezTech(dice []int8, mod int) error {
	if len(dice) > 0 {
		r.roll(dice)
		r.evaluateRezTech()
		return nil
	} else {
		return errors.New("dice invalid")
	}
}

func (r *DiceRoll) evaluateGeneral() {
	result := 0
	for _, val := range r.Dice {
		result += int(val.Result)
	}
	r.Result = result
}
func (r *DiceRoll) rollGeneral(dice []int8, mod int) error {
	if len(dice) > 0 {
		r.roll(dice)
		r.evaluateGeneral()
		return nil
	} else {
		return errors.New("dice invalid")
	}
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
		max = 10
	} else if eyes > 0 {
		max = int(eyes)
	} else {
		return Die{}, errors.New("invalid number of eyes")
	}
	return Die{Eyes: eyes, Result: int8(rand.Intn(max) + 1)}, nil
}

func (r *Room) addPlayer(id int, name string, col string) error {
	if r.Players.Detail == nil {
		r.Players.Detail = make(map[int]Char)
	}
	player, ok := r.Players.Detail[id]
	if !ok {
		p := Char{Name: name, Color: col}
		r.Players.Detail[id] = p
		r.Players.Id = append(r.Players.Id, id)
	} else {
		if name != "" {
			player.Name = name
		}
		if col != "" {
			player.Color = col
		}
		r.Players.Detail[id] = player
	}
	return nil
}
