package cherryConnector

import (
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	cfacade "github.com/lgynico/project-copier/cherry/facade"
	clog "github.com/lgynico/project-copier/cherry/logger"
)

type (
	WSConnector struct {
		cfacade.Component
		Connector
		Options
		upgrade *websocket.Upgrader
	}

	WSConn struct {
		*websocket.Conn
		typ    int
		reader io.Reader
	}
)

func (*WSConnector) Name() string {
	return "websocket_connector"
}

func (p *WSConnector) OnAfterInit() {
}

func (p *WSConnector) OnStop() {
	p.Stop()
}

func NewWS(address string, opts ...Option) *WSConnector {
	if address == "" {
		return nil
	}

	ws := &WSConnector{
		Options: Options{
			address:  address,
			certFile: "",
			keyFile:  "",
			chanSize: 256,
		},
		upgrade: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(_ *http.Request) bool {
				return true
			},
		},
	}

	for _, opt := range opts {
		opt(&ws.Options)
	}

	ws.Connector = NewConnector(ws.chanSize)

	return ws
}

func (p *WSConnector) Start() {
	listener, err := p.GetListener(p.certFile, p.keyFile, p.address)
	if err != nil {
		// TODO err 输出占位符错误
		clog.Fatalf("failed to listen: %s", err)
		// TODO return ?
	}

	clog.Infof("websocket connector listening at Address %s", p.address)
	if p.certFile != "" || p.keyFile != "" {
		clog.Infof("certFile = %s, keyFile = %s", p.certFile, p.keyFile)
	}

	p.Connector.Start()

	http.Serve(listener, p)
}

func (p *WSConnector) Stop() {
	p.Connector.Stop()
}

func (p *WSConnector) SetUpgrade(upgrade *websocket.Upgrader) {
	if upgrade != nil {
		p.upgrade = upgrade
	}
}

func (p *WSConnector) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	wsConn, err := p.upgrade.Upgrade(rw, r, nil)
	if err != nil {
		clog.Infof("Upgrade fail,URI = %s, Error = %s", r.RequestURI, err.Error())
		return
	}

	conn := NewWSConn(wsConn)
	p.InChan(&conn)
}

func NewWSConn(conn *websocket.Conn) WSConn {
	return WSConn{
		Conn: conn,
	}
}

func (p *WSConn) Read(b []byte) (int, error) {
	if p.reader == nil {
		t, r, err := p.NextReader()
		if err != nil {
			return 0, err
		}
		p.typ = t
		p.reader = r
	}

	n, err := p.reader.Read(b)
	if err != nil {
		if err != io.EOF {
			return n, err
		}

		_, r, err := p.NextReader()
		if err != nil {
			return 0, err
		}

		p.reader = r
	}

	return n, nil
}

func (p *WSConn) Write(b []byte) (int, error) {
	err := p.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (p *WSConn) SetDeadline(t time.Time) error {
	if err := p.SetReadDeadline(t); err != nil {
		return err
	}

	return p.SetWriteDeadline(t)
}
