package packetsender

import (
	"github.com/gorilla/websocket"
	"html/template"
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
	homeTemplate.Execute(w, "ws://"+r.Host+"/ws")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<head>
<meta charset="utf-8">
<script>
window.addEventListener("load", function(evt) {
	var packetNum = 0;
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
		ws.binaryType = 'arraybuffer';
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
			var arr = new Uint8Array(evt.data);
			console.log(packetNum++);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server,
"Send" to send a message to the server and "Close" to close the connection.
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
