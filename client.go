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

type Player *triki.Player

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

func listenConnection(conn net.Conn) *triki.Board {
	dec := gob.NewDecoder(conn)
	board := &triki.Board{}
	dec.Decode(board)
	return board
}

func play(vsPlayer Player, me Player, conn net.Conn, board *triki.Board) {
	if board == nil {
		board = triki.NewBoard()
	}
	enc := gob.NewEncoder(conn)
	var posx, posy int
	for board.Available > 0 {
		board.Print()
		fmt.Println("Posicion de su Jugada: ")
		_, err := fmt.Scanf("%d %d", &posx, &posy)
		if err != nil {
			continue
		}
		board.Play(me, float64(posx), float64(posy))
		enc.Encode(*board)
		board = listenConnection(conn)
	}
}

func callPlayer(vsPlayer Player, me Player) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8081")
	assertNoError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	assertNoError(err)
	defer conn.Close()
	err = conn.SetNoDelay(true)
	assertNoError(err)

	// send me for the other player know me
	enc := gob.NewEncoder(conn)
	enc.Encode(*me)

	play(vsPlayer, me, conn, nil)
}

func getVsPlayer(conn net.Conn) Player {
	dec := gob.NewDecoder(conn)
	vsPlayer := &triki.Player{}
	dec.Decode(vsPlayer)
	return vsPlayer
}

func startServer(quit chan int, me Player) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8081")
	assertNoError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)

	for {
		conn, err := listener.Accept()
		me.Symbol = triki.X
		hasVs <- true
		if err != nil {
			continue
		}
		vsPlayer := getVsPlayer(conn)
		board := listenConnection(conn)
		play(vsPlayer, me, conn, board)
		conn.Close()
		quit <- 1
	}
}

func getIpAddress() string {
	addrs, err := net.InterfaceAddrs()
	assertNoError(err)
	for _, a := range addrs {
		ipnet, ok := a.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			fmt.Println("your ip address is: ", ipnet.IP.String())
			return ipnet.IP.String()
		}
	}
	return ""
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
	quit := make(chan int)

	fmt.Print("Username: ")
	fmt.Scanf("%s", &username)

	me := &triki.Player{Id: encryptId(),
		Status: AVAILABLE,
		Uname:  username,
		Ip:     getIpAddress()}

	var reply []Player
	go startServer(quit, me)
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
	var vsPlayer Player
	found := false
	has_conn := false
	_select := func(p Player) {
		if vs == p.Uname {
			vsPlayer = p
			found = true
		}
	}
	for !found && !has_conn {
		fmt.Println("Escoger Jugador: ")
		each(reply, noPrint(me))
		fmt.Scanf("%s", &vs)
		each(reply, _select)
		select {
		case <-hasVs:
			has_conn = true
			break
		default:
			continue
		}
	}
	if found {
		vsPlayer.Symbol = triki.X
		me.Symbol = triki.O
		callPlayer(vsPlayer, me)
		close(quit)
		// err = client.Call("SessionManager.SelectPlayer", []string{vs, me.Id}, nil)
		// assertNoError(err)
	}
	<-quit
}
