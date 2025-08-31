package cherryFacade

import "net"

type (
	IConnector interface {
		IComponent
		Start()
		Stop()
		OnConnect(fn OnConnectFunc)
	}

	OnConnectFunc func(conn net.Conn)
)
