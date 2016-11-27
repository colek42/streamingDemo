package main

import (
	"flag"
	"github.com/colek42/streamingDemo/packetsender"
)

var (
	uri string
)

func init() {
	flag.StringVar(&uri, "f", "udp://@234.5.5.5:8209", "ex. udp://@345.5.5.5:8209")
	flag.Parse()
}

func main() {
	packetsender.OpenStream(uri)
}
