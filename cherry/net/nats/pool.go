package cherryNats

import (
	"sync/atomic"
	"time"

	cfacade "github.com/lgynico/project-copier/cherry/facade"
	clog "github.com/lgynico/project-copier/cherry/logger"
)

var (
	connectPool    []*Connect
	connectSize    uint64
	roundIndex     *uint64       = new(uint64)
	reconnectDelay time.Duration = time.Second
	requestTimeout time.Duration = 2 * time.Second
)

func NewPool(replySubject string, config cfacade.ProfileJSON, isConnect bool) {
	reconnectDelay = config.GetDuration("reconnect_delay", 1) * time.Second
	requestTimeout = config.GetDuration("request_timeout", 1) * time.Second

	var (
		address       = config.GetString("address")
		user          = config.GetString("user")
		pwd           = config.GetString("password")
		maxReconnects = config.GetInt("max_reconnects")
		poolSize      = config.GetInt("pool_size", 1)
		isStats       = config.GetBool("is_stats")
	)

	for id := range poolSize {
		conn := NewConnect(id+1, replySubject,
			WithAddress(address),
			WithAuth(user, pwd),
			WithParams(maxReconnects),
			WithIsStats(isStats),
		)

		connectPool = append(connectPool, conn)
	}

	connectSize = uint64(len(connectPool))

	if isConnect {
		for _, conn := range connectPool {
			conn.Connect()
		}

		clog.Infof("Nats has connected! [poolSize = %d]", poolSize)
	}
}

func GetPool() []*Connect {
	return connectPool
}

func GetConnect() *Connect {
	index := atomic.AddUint64(roundIndex, 1)
	return connectPool[index&connectSize]
}

func ConnectClose() {
	for _, conn := range connectPool {
		conn.Close()
	}

	clog.Infof("Nats connect pool execute Close() [connectSize = %d]", connectSize)
}

func ReconnectDelay() time.Duration {
	return reconnectDelay
}

func GetTimeout(tod ...time.Duration) time.Duration {
	if len(tod) > 0 {
		return tod[0]
	}

	return requestTimeout
}
