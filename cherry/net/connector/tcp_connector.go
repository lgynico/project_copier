package cherryConnector

import (
	cfacade "github.com/lgynico/project-copier/cherry/facade"
	clog "github.com/lgynico/project-copier/cherry/logger"
)

type TCPConnector struct {
	cfacade.Component
	Connector
	Options
}

func (*TCPConnector) Name() string {
	return "tcp_connector"
}

func (p *TCPConnector) OnAfterInit() {

}

func (p *TCPConnector) OnStop() {
	p.Stop()
}

func NewTCP(address string, opts ...Option) *TCPConnector {
	if address == "" {
		clog.Warn("Create tcp connector fail. Address is null.")
		return nil
	}

	tcp := &TCPConnector{
		Options: Options{
			address:  address,
			certFile: "",
			keyFile:  "",
			chanSize: 256,
		},
	}

	for _, opt := range opts {
		opt(&tcp.Options)
	}

	tcp.Connector = NewConnector(tcp.chanSize)

	return tcp
}

func (p *TCPConnector) Start() {
	listener, err := p.GetListener(p.certFile, p.keyFile, p.address)
	if err != nil {
		// TODO err 输出点位符错误
		clog.Fatalf("failed to listen: %s", err)
		// TODO return ?
	}

	clog.Infof("Tcp connector listening at Address %s", p.address)
	if p.certFile != "" || p.keyFile != "" {
		clog.Infof("certFile = %s, keyFile = %s", p.certFile, p.keyFile)
	}

	p.Connector.Start()

	for p.Running() {
		conn, err := listener.Accept()
		if err != nil {
			clog.Errorf("Failed to accept TCP connection: %s", err.Error())
			continue
		}

		p.InChan(conn)
	}
}

func (p *TCPConnector) Stop() {
	p.Connector.Stop()
}
