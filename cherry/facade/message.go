package cherryFacade

import (
	"strings"

	cconst "github.com/lgynico/project-copier/cherry/const"
	cerr "github.com/lgynico/project-copier/cherry/error"
	cstring "github.com/lgynico/project-copier/cherry/extend/string"
	cherryTime "github.com/lgynico/project-copier/cherry/extend/time"
	cproto "github.com/lgynico/project-copier/cherry/net/proto"
	"github.com/nats-io/nats.go"
)

type (
	Message struct {
		BuildTime  int64
		PostTime   int64
		Source     string
		Target     string
		FuncName   string
		Session    *cproto.Session
		Args       any
		Header     nats.Header
		Reply      string
		IsCluster  bool
		ChanResult chan any
		targetPath *ActorPath
	}

	ActorPath struct {
		NodeID  string
		ActorID string
		ChildID string
	}
)

func GetMessage() Message {
	return Message{
		BuildTime: cherryTime.Now().ToMillisecond(),
	}
}

func BuildClusterMessage(packet *cproto.ClusterPacket) Message {
	return Message{
		BuildTime: packet.BuildTime,
		Source:    packet.SourcePath,
		Target:    packet.TargetPath,
		FuncName:  packet.FuncName,
		IsCluster: true,
		Session:   packet.Session,
		Args:      packet.ArgBytes,
	}
}

func NewActorPath(nodeID, actorID, childID string) *ActorPath {
	return &ActorPath{
		NodeID:  nodeID,
		ActorID: actorID,
		ChildID: childID,
	}
}

func NewChildPath(nodeID, actorID, childID any) string {
	if childID == "" {
		return NewPath(nodeID, actorID)
	}

	return cstring.ToString(nodeID) + cconst.DOT + cstring.ToString(actorID) + cconst.DOT + cstring.ToString(childID)
}

func NewPath(nodeID, actorID any) string {
	return cstring.ToString(nodeID) + cconst.DOT + cstring.ToString(actorID)
}

func ToActorPath(path string) (*ActorPath, error) {
	if path == "" {
		return nil, cerr.ErrActorPath
	}

	p := strings.Split(path, cconst.DOT)
	pLen := len(p)

	if pLen == 2 {
		return NewActorPath(p[0], p[1], ""), nil
	}

	if pLen == 3 {
		return NewActorPath(p[0], p[1], p[2]), nil
	}

	return nil, cerr.ErrActorPath
}

func (p *Message) TargetPath() *ActorPath {
	if p.targetPath == nil {
		p.targetPath, _ = ToActorPath(p.Target)
	}

	return p.targetPath
}

func (p *Message) IsReply() bool {
	return p.Reply != ""
}

func (p *Message) Destory() {
	p.targetPath = nil
	p.Session = nil
	p.Args = nil
	p.Header = nil
	p.ChanResult = nil
}

func (p *ActorPath) IsChild() bool {
	return p.ChildID != ""
}

func (p *ActorPath) IsParent() bool {
	return p.ChildID == ""
}

func (p *ActorPath) String() string {
	return NewChildPath(p.NodeID, p.ActorID, p.ChildID)
}
