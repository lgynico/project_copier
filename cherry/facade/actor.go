package cherryFacade

import (
	"time"

	creflect "github.com/lgynico/project-copier/cherry/extend/reflect"
)

type (
	IActorSystem interface {
		GetIActor(id string) (IActor, bool)
		CreateActor(id string, handler IActorHandler) (IActor, error)

		PostRemote(m *Message) bool
		PostLocal(m *Message) bool

		Call(source, target, funcName string, arg any) int32
		CallWait(source, target, funcName string, arg, reply any) int32
		CallType(nodeType, actorID, funcName string, arg any) int32

		SetLocalInvoke(invoke InvokeFunc)
		SetRemoteInvoke(invoke InvokeFunc)
		SetCallTimeout(d time.Duration)
		SetArrivalTimeout(t int64)
		SetExecutionTimeout(t int64)
	}

	IActor interface {
		App() IApplication
		ActorID() string
		Path() *ActorPath

		Call(targetPath, funcName string, arg any) int32
		CallWait(targetPath, funcName string, arg, reply any) int32
		CallType(nodeType, actorID, funcName string, arg any) int32

		PostRemote(m *Message) bool
		PostLocal(m *Message) bool

		LastAt() int64
		Exit()
	}

	IActorHandler interface {
		AliasID() string
		OnInit()
		OnStop()
		OnLocalReceive(m *Message) (next, invoke bool)
		OnRemoteReceive(m *Message) (next, invoke bool)
		OnFindChild(m *Message) (IActor, bool)
	}

	IActorChlid interface {
		Create(id string, handler IActorHandler) (IActor, error)
		Get(id string) (IActor, bool)
		Remove(id string)
		Each(fun func(i IActor))
		Call(childID, funcName string, arg any)
		CallWait(childID, funcName string, arg, reply any) int32
	}
)

type (
	InvokeFunc func(app IApplication, fi *creflect.FuncInfo, m *Message)

	IEventData interface {
		Name() string
		UniqueID() int64
	}
)
