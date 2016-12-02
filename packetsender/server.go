package packetsender

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type VideoRoom struct {
	uri       string
	tsPackets chan tspacket
}

type client struct {
	conn *websocket.Conn
}

type tspacket struct {
	pts       int64
	dts       int64
	timeStamp int64
	data      []byte
}

var webroot = "/home/cole/go/src/github.com/colek42/streamingDemo/packetrecevier"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func startVideo(room *VideoRoom) {
	log.Printf("starting Video")
	OpenStream(room.uri, room.tsPackets)
}

func serveWs(w http.ResponseWriter, r *http.Request) {

	uri := r.URL.Query().Get("uri")
	log.Println(uri)
	// if uri == "" {
	// 	return
	// }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error %v", err)
		return
	}

	pktChan := make(chan tspacket, 10)

	vr := &VideoRoom{
		uri:       uri,
		tsPackets: pktChan,
	}
	go startVideo(vr)

	c := &client{
		conn: conn,
	}

	messageWriter(c, vr)
	log.Println("Closed")
}

func Serve() {
	http.HandleFunc("/ws", serveWs)
	http.HandleFunc("/", home)
	err := http.ListenAndServe("0.0.0.0:8787", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func messageWriter(c *client, v *VideoRoom) {
	for {
		select {
		case pkt, ok := <-v.tsPackets:
			if !ok {
				continue
			}
			err := c.conn.WriteMessage(websocket.BinaryMessage, pkt.data)
			if err != nil {
				log.Printf("Err Sending Websocket: %v", err)
			}

		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	log.Printf("Req: %v", webroot+r.URL.Path)
	http.ServeFile(w, r, webroot+r.URL.Path)
}
