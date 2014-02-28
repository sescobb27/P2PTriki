package triki

import (
	"errors"
	"math"
)

type Board struct {
	Tboard    [9]string
	Available int
}

func NewBoard() *Board {
	b := &Board{Tboard: [9]string{},
		Available: 9}
	for x := 0; x < 9; x++ {
		b.Tboard[x] = E
	}
	return b
}

func (b *Board) CheckWin(p *Player) bool {
	return (((b.Tboard[0] == p.Symbol) && (b.Tboard[1] == p.Symbol) && (b.Tboard[2] == p.Symbol)) ||
		((b.Tboard[0] == p.Symbol) && (b.Tboard[4] == p.Symbol) && (b.Tboard[8] == p.Symbol)) ||
		((b.Tboard[0] == p.Symbol) && (b.Tboard[3] == p.Symbol) && (b.Tboard[6] == p.Symbol)) ||
		((b.Tboard[1] == p.Symbol) && (b.Tboard[4] == p.Symbol) && (b.Tboard[7] == p.Symbol)) ||
		((b.Tboard[2] == p.Symbol) && (b.Tboard[4] == p.Symbol) && (b.Tboard[6] == p.Symbol)) ||
		((b.Tboard[2] == p.Symbol) && (b.Tboard[5] == p.Symbol) && (b.Tboard[8] == p.Symbol)) ||
		((b.Tboard[3] == p.Symbol) && (b.Tboard[4] == p.Symbol) && (b.Tboard[5] == p.Symbol)))
}

func (b *Board) Play(player *Player, posx, posy float64) error {
	pos := int(math.Abs(posx - posy))
	if b.Tboard[pos] == E {
		b.Tboard[pos] = player.Symbol
		b.Available--
		return nil
	} else {
		return errors.New("Invalid Position")
	}
}

func (b *Board) Print() {
	for i := 0; i < len(b.Tboard); i++ {
		if i%3 == 0 {
			println()
		}
		print(b.Tboard[i])
	}
	println()
}
