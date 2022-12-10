package main

type Statistic struct {
	nbrCoC       int
	nbrRezTech   int
	nbrGeneral   int
	nbrDiceRolls int
	nbrPlayer    int
}

func updateStatistics(rooms map[int]Room, players map[int][]int) (Statistic, bool) {
	ret := Statistic{
		nbrCoC:       0,
		nbrRezTech:   0,
		nbrGeneral:   0,
		nbrDiceRolls: 0,
		nbrPlayer:    0}

	for _, val := range rooms {
		switch val.Game {
		case CoC:
			ret.nbrCoC++
		case RezTech:
			ret.nbrRezTech++
		case General:
			ret.nbrGeneral++
		}
		ret.nbrDiceRolls += len(val.DiceRolls)
	}
	ret.nbrPlayer = len(players)

	return ret, true
}
