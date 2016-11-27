package packetsender

import (
	"github.com/Novetta/pwcop/lib/messaging"
	"github.com/gorilla/mux"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL: %v", r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", 404)
		return
	}

	http.Error(w, "Do This", 404)
	//TODO, server homepage
}

func startVideo(room VideoRoom) {
	OpenStream(room.uri, room.tsPackets)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uri := params["uri"]
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
	router := mux.NewRouter
	router.HandleFunc("/", serveHome)
	router.HandleFunc("/ws/{uri}", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func messageWriter(c client, v VideoRoom) {
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
