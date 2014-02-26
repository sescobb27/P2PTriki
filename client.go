package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"
	// "triki"
)

func assertNoError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

var (
	me    *Player
	hasVs = make(chan bool)
)

const (
	PLAYING   = 0
	AVAILABLE = 1
)

type Player struct {
	Uname  string
	Status int
	Id     string
	Ip     string
}

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

func print(p Player) {
	if p.Uname != me.Uname {
		fmt.Println("* " + p.Uname)
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

func listenConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println(conn)
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	assertNoError(err)
	fmt.Println(string(buf))
}

func callPlayer(vsPlayer *Player) {
	// addr := vsPlayer.Ip + ":8080"
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8081")
	assertNoError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	assertNoError(err)
	defer conn.Close()
	err = conn.SetNoDelay(true)
	assertNoError(err)
	conn.Write([]byte("Hello\n"))
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
	me = &Player{Id: encryptId(),
		Status: AVAILABLE,
		Uname:  username,
		Ip:     getIpAddress()}
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
		each(reply, print)
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
		callPlayer(vsPlayer)
		// err = client.Call("SessionManager.SelectPlayer", []string{vs, me.Id}, nil)
		// assertNoError(err)
	}
}
