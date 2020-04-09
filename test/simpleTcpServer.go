package test

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

func RunLocalTcpServer(port int) {
	var ipAddr = "127.0.0.1:" + strconv.Itoa(port)
	addr, err := net.ResolveTCPAddr("tcp", ipAddr)
	if err != nil {
		log.Fatalln(err)
	}

	server, err2 := net.ListenTCP("tcp", addr)
	if err2 != nil {
		log.Fatalln(err2)
	}
	log.Println(">>>>>>: server listen at: ", addr.String())
	for {
		clientSocket, err := server.AcceptTCP()
		if err != nil {
			log.Println(">>>>>>>: error ", err)
			continue
		}
		log.Println(">>>>>>>: recieve a new con from " + clientSocket.RemoteAddr().String())
		_, _ = io.Copy(os.Stdout, clientSocket)
	}

}
