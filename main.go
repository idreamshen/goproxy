package main

import (
	"os"
	"os/signal"

	"idreamshen.com/goproxy/metrics"
	"idreamshen.com/goproxy/proxy"
)

func main() {
	metrics.StartMetricsServer()
	proxy.StartProxyServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
