package triki

import (
	"errors"
)

type Board struct {
	Tboard    [3][3]string
	Available int
}

func NewBoard() *Board {
	b := &Board{Tboard: [3][3]string{},
		Available: 9}
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			b.Tboard[x][y] = E
		}
	}
	return b
}

func (b *Board) CheckWin(p *Player) bool {
	return (((b.Tboard[0][0] == p.Symbol) && (b.Tboard[0][1] == p.Symbol) && (b.Tboard[0][2] == p.Symbol)) ||
		((b.Tboard[0][0] == p.Symbol) && (b.Tboard[1][1] == p.Symbol) && (b.Tboard[2][2] == p.Symbol)) ||
		((b.Tboard[0][0] == p.Symbol) && (b.Tboard[1][0] == p.Symbol) && (b.Tboard[2][0] == p.Symbol)) ||
		((b.Tboard[0][1] == p.Symbol) && (b.Tboard[1][1] == p.Symbol) && (b.Tboard[2][1] == p.Symbol)) ||
		((b.Tboard[0][2] == p.Symbol) && (b.Tboard[1][1] == p.Symbol) && (b.Tboard[2][0] == p.Symbol)) ||
		((b.Tboard[0][2] == p.Symbol) && (b.Tboard[1][2] == p.Symbol) && (b.Tboard[2][2] == p.Symbol)) ||
		((b.Tboard[1][0] == p.Symbol) && (b.Tboard[1][1] == p.Symbol) && (b.Tboard[1][2] == p.Symbol)) ||
		((b.Tboard[2][0] == p.Symbol) && (b.Tboard[2][1] == p.Symbol) && (b.Tboard[2][2] == p.Symbol)))
}

func (b *Board) Play(player *Player, posx, posy int) error {
	if b.Tboard[posx][posy] == E {
		b.Tboard[posx][posy] = player.Symbol
		b.Available--
		return nil
	} else {
		return errors.New("Invalid Position")
	}
}

func (b *Board) Print() {
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			print(b.Tboard[x][y])
		}
		println()
	}
}
