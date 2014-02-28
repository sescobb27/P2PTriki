package triki

import (
	"fmt"
	"os"
)

const (
	X string = "[X]"
	O string = "[O]"
	E string = "[ ]"
)

type Player struct {
	Id     string
	Uname  string
	Symbol string
	Status int
	Ip     string
}

func (p *Player) ask() {
	var response string
	fmt.Print("Continuar? (Y/n): ")
	fmt.Scanf("%s", &response)
	switch response {
	case "n":
		os.Exit(0)
	}
}
