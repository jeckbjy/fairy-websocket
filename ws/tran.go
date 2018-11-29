package ws

import (
	"net"
	"net/http"

	"github.com/jeckbjy/fairy"

	"github.com/gorilla/websocket"
	"github.com/jeckbjy/fairy/base"
)

// NewTran create websocket transport
func NewTran() fairy.ITran {
	return base.NewTran(&wsTran{})
}

type httpServer struct {
	upgrader *websocket.Upgrader
	cb       base.OnAccept
}

func (server *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := server.upgrader.Upgrade(w, r, nil)
	wconn := &wsConn{Conn: conn}
	server.cb(wconn, err)
}

type wsTran struct {
	websocket.Upgrader
}

func (*wsTran) Connect(host string) (net.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		return nil, err
	}

	return &wsConn{Conn: conn}, nil
}

func (*wsTran) Listen(host string) (net.Listener, error) {
	return net.Listen("tcp", host)
}

func (wt *wsTran) Serve(l net.Listener, cb base.OnAccept) {
	server := http.Server{Handler: &httpServer{upgrader: &wt.Upgrader, cb: cb}}
	server.Serve(l)
}
