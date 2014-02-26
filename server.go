package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"triki/server"
)

var (
	sessionManager *server.SessionManager
)

func assertNoError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func startListening(listener *net.TCPListener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error Accepting connection: ", err.Error())
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func startServer() {
	var err error
	var tcpAddr *net.TCPAddr
	var listener *net.TCPListener

	tcpAddr, err = net.ResolveTCPAddr("tcp", ":8080")
	assertNoError(err)
	listener, err = net.ListenTCP("tcp", tcpAddr)
	assertNoError(err)

	startListening(listener)
}

func main() {
	sessionManager = server.InitializeSessionManager()
	rpc.Register(sessionManager)
	startServer()
}
