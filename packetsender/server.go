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
	log.Printf("serveWs")
	uri := "udp://234.5.5.5:8209"
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
	//r.Path("/").Handler(http.FileServer(http.Dir(webroot)))

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
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})

			}
			c.conn.WriteMessage(websocket.BinaryMessage, pkt.data)
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, webroot+r.URL.Path[1:])
}
