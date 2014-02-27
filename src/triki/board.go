package triki

import (
	"errors"
	"fmt"
)

type Board struct {
	board   [9]string
	Player1 Player
	Player2 Player
}

func NewBoard() *Board {
	return makeBoard()
}

func makeBoard() *Board {
	b := &Board{board: [9]string{}}
	for x := 0; x < len(b.board); x++ {
		b.board[x] = E
		fmt.Println(b.board[x])
	}
	fmt.Print("Error: ")
	fmt.Println(b)
	return b
}

func (b *Board) CheckWin(p *Player) bool {
	return (((b.board[0] == p.Symbol) && (b.board[1] == p.Symbol) && (b.board[2] == p.Symbol)) ||
		((b.board[0] == p.Symbol) && (b.board[4] == p.Symbol) && (b.board[8] == p.Symbol)) ||
		((b.board[0] == p.Symbol) && (b.board[3] == p.Symbol) && (b.board[6] == p.Symbol)) ||
		((b.board[1] == p.Symbol) && (b.board[4] == p.Symbol) && (b.board[7] == p.Symbol)) ||
		((b.board[2] == p.Symbol) && (b.board[4] == p.Symbol) && (b.board[6] == p.Symbol)) ||
		((b.board[2] == p.Symbol) && (b.board[5] == p.Symbol) && (b.board[8] == p.Symbol)) ||
		((b.board[3] == p.Symbol) && (b.board[4] == p.Symbol) && (b.board[5] == p.Symbol)))
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
