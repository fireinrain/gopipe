package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

//接受客户端
type InboundCon struct {
	ClinetIp   string
	ClientPort int
	Protocol   string
}

//服务端
type OutBoundCon struct {
	ServerIp   string
	ServerPort int
	Protocol   string
}

/**
pattern map design
client ->|proxy-server| -> |remote server|
client <-|proxy-server| <- |remote server|

*/
var serverPort = 8000
var remoteAddr = "127.0.0.1:9000"
var remoteTcpAddr *net.TCPAddr

func main() {

	//当前本机的地址
	serverAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", serverPort))
	simpleCheckError(err)
	//远程需要代理的地址 端口
	targetAddr, err2 := net.ResolveTCPAddr("tcp", remoteAddr)
	remoteTcpAddr = targetAddr
	log.Println("targetAddr:", targetAddr.String())
	simpleCheckError(err2)

	//监听本机
	server, err := net.ListenTCP("tcp", serverAddr)
	messageCheckError(err, "failed to listen at "+serverAddr.String())
	defer server.Close()

	for {
		clientSocket, err := server.AcceptTCP()
		if err != nil {
			log.Printf(">>>>>>: error,%s", "failed to accept connection from client "+clientSocket.RemoteAddr().String())
			continue

		}
		//启动单独的协程处理
		go handleClientProxy(clientSocket)

	}

}

//处理转发
func handleClientProxy(client *net.TCPConn) {
	defer client.Close()
	log.Printf("client '%s' connected!\n", client.RemoteAddr().String())

	//保持链接
	_ = client.SetKeepAlive(true)
	_ = client.SetKeepAlivePeriod(time.Second * 15)

	//链接远程服务端口
	remoteServerCon, err := net.DialTCP("tcp", nil, remoteTcpAddr)
	messageCheckError(err, "failed to connect to "+remoteAddr)
	defer remoteServerCon.Close()

	//保持链接
	_ = remoteServerCon.SetKeepAlive(true)
	_ = remoteServerCon.SetKeepAlivePeriod(time.Second * 15)

	//接收数据
	//能否使用bufio来提高性能？
	var buffer = make([]byte, 100000)
	for {
		n, err := client.Read(buffer)
		if err != nil {
			break
		}
		//转发数据
		n, err = remoteServerCon.Write(buffer[:n])
		if err != nil {
			break
		}
	}

}

//检查error
func simpleCheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

//检查error，带描述信息
func messageCheckError(err error, message string) {
	if err != nil {
		log.Fatalf(">>>>>>: error,%s", message)
	}
}
