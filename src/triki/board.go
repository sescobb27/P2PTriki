package triki

import (
	"errors"
)

type Board struct {
	board   []string
	player1 Player
	player2 Player
}

func (b *Board) makeBoard() {
	for x := 0; x < 9; x++ {
		b.board[x] = E
	}
}

func (b *Board) String() string {
	tmp := ""
	pos := 0
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			pos++
			tmp += b.board[pos] + " "
		}
		tmp += "\n"
	}
	return tmp
}

func (b *Board) play(player Player, pos int) error {
	if b.board[pos] == E {
		b.board[pos] = player.Symbol
		return nil
	} else {
		return errors.New("Invalid Position")
	}
}

func init() {
	board := &Board{board: make([]string, 0, 9)}
	board.makeBoard()
}
