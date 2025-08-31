package cherryFacade

import (
	"time"

	jsoniter "github.com/json-iterator/go"
)

type (
	// INode 节点信息
	INode interface {
		NodeID() string        // 节点ID
		NodeType() string      // 节点类型
		Address() string       // 对外网络监听地址（前端节点用）
		RpcAddress() string    // Rpc监听地址（未使用）
		Settings() ProfileJSON // 节点配置
		Enabled() bool         // 是否启用
	}

	IApplication interface {
		INode
		Running() bool
		DieChan() chan bool
		IsFrontend() bool

		Register(components ...IComponent)
		Find(name string) IComponent
		Remove(name string) IComponent
		All() []IComponent

		OnShutdown(fn ...func())
		Startup()
		Shutdown()

		Serializer() ISerializer
		Discovery() IDiscovery
		Cluster() ICluster
		ActorSystem() IActorSystem
	}

	ProfileJSON interface {
		jsoniter.Any
		GetConfig(path ...any) ProfileJSON
		GetString(path any, defaultValue ...string) string
		GetBool(path any, defaultValue ...bool) bool
		GetInt(path any, defaultValue ...int) int
		GetInt32(path any, defaultValue ...int32) int32
		GetInt64(path any, defaultValue ...int64) int64
		GetDuration(path any, defaultValue ...time.Duration) time.Duration

		Unmarshal(ptrVal any) error
	}
)
