package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

func ioCopy(dst io.Writer, src io.Reader) {
	// defer func() {
	// 	src.Close()
	// 	fmt.Println("close")
	// }()

	fmt.Println("copy start")
	w, err := io.Copy(dst, src)
	if err != nil {
		panic(err)
	}
	fmt.Printf("copy %d", w)
}

func handleClient(c net.Conn) {
	reader := bufio.NewReader(c)

	reqLine, err := reader.ReadString('\r')

	if err != nil {
		panic(err)
	}

	reqLine = strings.TrimSpace(reqLine)
	parts := strings.Split(reqLine, " ")
	if len(parts) < 3 {
		log.Println("Invalid request line:", reqLine)
		return
	}

	method, target := parts[0], parts[1]
	log.Println("Incoming request:", method, target)

	if method == "CONNECT" {
		reader.Discard(reader.Buffered())
		fmt.Println("discard")
	} else {

	}

	var addr string
	if method == "CONNECT" {
		addr = target
	} else {
		u, err := url.Parse(target)
		if err != nil {
			panic(err)
		}

		if u.Port() == "" {
			addr = u.Hostname() + ":80"
		} else {
			addr = u.Hostname() + ":" + u.Port()
		}
	}

	//addr = "192.168.89.99:8080"
	fmt.Println(addr)
	remote, err := net.Dial("tcp", addr) //建立服务端和代理服务器的tcp连接
	if err != nil {
		panic(err)
	}

	log.Printf("establisted")

	//_, err = c.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	if method != "CONNECT" {
		remote.Write([]byte(method + " " + "/" + " " + "HTTP/1.1"))
		remote.Write([]byte("\r\n"))
	}

	reader.ReadBytes('\n')

	go ioCopy(remote, reader)
	go ioCopy(c, remote)

}

func tcpServer() {
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	for {
		client, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go handleClient(client)
	}
}

func main() {
	tcpServer()
}
