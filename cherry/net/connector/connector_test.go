package cherryConnector

import (
	"fmt"
	"net"
	"sync"
	"testing"

	clog "github.com/lgynico/project-copier/cherry/logger"
)

func TestNewTCPConnector(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	tcp := NewTCP(":9071")
	tcp.OnConnect(func(conn net.Conn) {
		clog.Infof("new net.Conn = %s", conn.RemoteAddr())
	})

	tcp.Start()

	wg.Wait()
}

func TestNewWSConnector(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	ws := NewWS(":9071")
	ws.OnConnect(func(conn net.Conn) {
		clog.Infof("new net.Conn = %s", conn.RemoteAddr())
		go func() {
			for {
				buf := make([]byte, 2048)
				for {
					n, err := conn.Read(buf)
					if err != nil {
						return
					}

					fmt.Println(string(buf[:n]))
				}
			}
		}()
	})

	ws.Start()

	wg.Wait()
}
