package main

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"
	"triki"
)

func assertNoError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

var (
	hasVs = make(chan bool)
)

const (
	PLAYING   = 0
	AVAILABLE = 1
)

type Player triki.Player

func encryptId() string {
	crutime := time.Now().Unix()
	hash := sha256.New()
	hash.Write([]byte(strconv.FormatInt(crutime, 10)))
	sid := hex.EncodeToString(hash.Sum(nil))
	return sid
}

func each(arr []Player, f func(p Player)) {
	for _, p := range arr {
		f(p)
	}
}

func noPrint(me Player) func(p Player) {
	return func(p Player) {
		if p.Uname != me.Uname {
			fmt.Println("* " + p.Uname)
		}
	}
}

func listenConnection(conn net.Conn) {
	defer conn.Close()
	dec := gob.NewDecoder(conn)
	board := &triki.Board{}
	dec.Decode(board)
	fmt.Print("Error")
	fmt.Println(board)
}

func callPlayer(vsPlayer *Player, me *Player) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8081")
	assertNoError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	assertNoError(err)
	defer conn.Close()
	err = conn.SetNoDelay(true)
	assertNoError(err)
	board := triki.NewBoard()
	fmt.Println(board)
	enc := gob.NewEncoder(conn)
	enc.Encode(board)
}

func startServer() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8081")
	assertNoError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)

	for {
		conn, err := listener.Accept()
		hasVs <- true
		if err != nil {
			continue
		}
		listenConnection(conn)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "server")
		os.Exit(1)
	}
	serverAddr := os.Args[1]

	var username string
	var client *rpc.Client
	var err error

	client, err = rpc.Dial("tcp", serverAddr+":8080")
	assertNoError(err)

	go startServer()

	fmt.Print("Username: ")
	fmt.Scanf("%s", &username)
	me := &Player{Id: encryptId(),
		Status: AVAILABLE,
		Uname:  username}
	var reply []Player
	for {
		err = client.Call("SessionManager.SessionStart", me, &reply)
		assertNoError(err)

		if len(reply)-1 == 0 {
			fmt.Println("No hay jugadores disponibles\nEsperando Jugadores...")
			time.Sleep(5000 * time.Millisecond)
		} else {
			break
		}
	}

	var vs string
	var vsPlayer *Player
	found := false
	_select := func(p Player) {
		if vs == p.Uname {
			vsPlayer = &p
			found = true
		}
	}
	for !found {
		fmt.Println("Escoger Jugador: ")
		each(reply, noPrint(*me))
		fmt.Scanf("%s", &vs)
		each(reply, _select)
		select {
		case <-hasVs:
			break
		default:
			continue
		}
	}
	if found {
		vsPlayer.Symbol = triki.X
		me.Symbol = triki.O
		callPlayer(vsPlayer, me)
		// err = client.Call("SessionManager.SelectPlayer", []string{vs, me.Id}, nil)
		// assertNoError(err)
	}
}
