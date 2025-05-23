package proxy

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"idreamshen.com/goproxy/metrics"
)

func StartProxyServer() {
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			client, err := l.Accept()
			if err != nil {
				panic(err)
			}
			go handleClient(client)
		}
	}()
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

	pTimer := prometheus.NewTimer(metrics.HistogramProxyProcessSec.WithLabelValues(addr))
	defer pTimer.ObserveDuration()

	//addr = "192.168.89.99:8080"
	remote, err := net.Dial("tcp", addr) //建立服务端和代理服务器的tcp连接
	if err != nil {
		panic(err)
	}

	log.Printf("establisted tcp")

	metrics.CounterProxyProcessTotal.Inc()

	if method == "CONNECT" {
		_, err = c.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	} else {
		remote.Write([]byte(method + " " + "/" + " " + "HTTP/1.1"))
		remote.Write([]byte("\r\n"))
		reader.ReadBytes('\n')
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go ioCopy(&wg, remote, reader)
	go ioCopy(&wg, c, remote)

	wg.Wait()

	if err := remote.Close(); err != nil {
		log.Printf("remote close err: %s", err.Error())
	}

	if err := c.Close(); err != nil {
		log.Printf("client close err: %s", err.Error())
	}
}

func ioCopy(wg *sync.WaitGroup, dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	wg.Done()

	if err != nil {
		panic(err)
	}
}
