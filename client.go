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

func encryptId() string {
	crutime := time.Now().Unix()
	hash := sha256.New()
	hash.Write([]byte(strconv.FormatInt(crutime, 10)))
	sid := hex.EncodeToString(hash.Sum(nil))
	return sid
}

func each(arr []*triki.Player, f func(p *triki.Player)) {
	for _, p := range arr {
		f(p)
	}
}

func noPrint(me *triki.Player) func(p *triki.Player) {
	return func(p *triki.Player) {
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

func play(me *triki.Player, conn net.Conn, board *triki.Board) bool {
	if board == nil {
		board = triki.NewBoard()
	}
	var posx, posy int
	var enc *gob.Encoder
	for board.Available > 0 {
		board.Print()
		fmt.Println("Posicion de su Jugada: ")
		_, err := fmt.Scanf("%d %d", &posx, &posy)
		if err != nil {
			continue
		}
		board.Play(me, posx, posy)
		board.Print()
		if board.CheckWin(me) {
			fmt.Println("Has Ganado")
			return true
		}
		enc = gob.NewEncoder(conn)
		enc.Encode(*board)
		board = listenConnection(conn)
	}
	return false
}

func callPlayer(vsPlayer *triki.Player, me *triki.Player) {
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

	won := play(me, conn, nil)
	if !won {
		fmt.Println("Has perdido")
	}
}

func getVsPlayer(conn net.Conn) *triki.Player {
	dec := gob.NewDecoder(conn)
	vsPlayer := &triki.Player{}
	dec.Decode(vsPlayer)
	return vsPlayer
}

func startServer(quit chan int, me triki.Player) {
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
		fmt.Println("Incomming connection from" + vsPlayer.Uname)
		board := listenConnection(conn)
		won := play(&me, conn, board)
		if !won {
			fmt.Println("Has perdido")
		}
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

	var reply []*triki.Player
	go startServer(quit, *me)
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
	var vsPlayer *triki.Player
	found := false
	has_conn := false
	_select := func(p *triki.Player) {
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
		me.Symbol = triki.O
		err = client.Call("SessionManager.SelectPlayer", []string{vsPlayer.Id, me.Id}, nil)
		assertNoError(err)
		callPlayer(vsPlayer, me)
		close(quit)
	}
	<-quit
}
