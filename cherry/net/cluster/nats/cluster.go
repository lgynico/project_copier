package cherryNatsCluster

import (
	"time"

	ccode "github.com/lgynico/project-copier/cherry/code"
	cerror "github.com/lgynico/project-copier/cherry/error"
	cfacade "github.com/lgynico/project-copier/cherry/facade"
	clog "github.com/lgynico/project-copier/cherry/logger"
	cnats "github.com/lgynico/project-copier/cherry/net/nats"
	cproto "github.com/lgynico/project-copier/cherry/net/proto"
	cprofile "github.com/lgynico/project-copier/cherry/profile"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Cluster struct {
	app               cfacade.IApplication
	prefix            string
	localSubject      string
	remoteSubject     string
	replySubject      string
	remoteTypeSubject string
}

func New(app cfacade.IApplication) cfacade.ICluster {
	return &Cluster{
		app: app,
	}
}

func (p *Cluster) loadNatsConfig() {
	natsConfig := cprofile.GetConfig("cluster").GetConfig("nats")
	if natsConfig.LastError() != nil {
		panic("cluster nats config not found.")
	}

	p.prefix = natsConfig.GetString("prefix", "node")
	p.localSubject = GetLocalSubject(p.prefix, p.app.NodeType(), p.app.NodeID())
	p.remoteSubject = GetRemoteSubject(p.prefix, p.app.NodeType(), p.app.NodeID())
	p.remoteTypeSubject = GetRemoteTypeSubject(p.prefix, p.app.NodeType())
	p.replySubject = GetReplySubject(p.prefix, p.app.NodeType(), p.app.NodeID())

	cnats.NewPool(p.replySubject, natsConfig, true)
}

func (p *Cluster) Init() {
	p.loadNatsConfig()

	p.localProcess()
	p.remoteProcess()
	p.remoteTypeProcess()

	clog.Info("Nats cluster execute OnInit()")
}

func (p *Cluster) Stop() {
	cnats.ConnectClose()

	clog.Info("Nats cluster execute OnStop()")
}

func (p *Cluster) localProcess() {
	processor := func(natsMsg *nats.Msg) {
		packet, err := cproto.UnmarshalPacket(natsMsg.Data)
		defer packet.Recycle()

		if err != nil {
			// TODO err 打印不对
			clog.Warnf("[localProcess] Unmarshal fail. [subject = %s, %s, err = %s]", natsMsg.Subject, packet.PrintLog(), err)
			return
		}

		message := cfacade.BuildClusterMessage(packet)
		p.app.ActorSystem().PostLocal(&message)
	}

	conn := cnats.GetConnect()
	err := conn.Subscribe(p.localSubject, processor)
	if err != nil {
		clog.Errorf("[localProcess] Create subscribe failed. [subject = %s, err = %v]", p.localSubject, err)
	}
}

func (p *Cluster) remoteProcess() {
	processor := func(natsMsg *nats.Msg) {
		packet, err := cproto.UnmarshalPacket(natsMsg.Data)
		defer packet.Recycle()

		if err != nil {
			clog.Warnf("[remoteProcess] Unmarshal fail. [subject = %s, %s, err = %v]", natsMsg.Subject, packet.PrintLog(), err)
			return
		}

		message := cfacade.BuildClusterMessage(packet)
		if len(natsMsg.Reply) > 0 {
			message.Header = natsMsg.Header
			message.Reply = natsMsg.Reply
		}

		p.app.ActorSystem().PostRemote(&message)
	}

	conn := cnats.GetConnect()
	err := conn.Subscribe(p.remoteSubject, processor)
	if err != nil {
		clog.Errorf("[remoteProcess] Create subscribe failed. [subject = %s, err = %v]", p.remoteSubject, err)
		return
	}
}

func (p *Cluster) remoteTypeProcess() {
	processor := func(natsMsg *nats.Msg) {
		packet, err := cproto.UnmarshalPacket(natsMsg.Data)
		defer packet.Recycle()

		if err != nil {
			clog.Warnf("[remoteTypeProcess] Unmarshal fail. [subject = %s, %s, err = %v]", natsMsg.Subject, packet.PrintLog(), err)
			return
		}

		message := cfacade.BuildClusterMessage(packet)

		p.app.ActorSystem().PostRemote(&message)
	}

	conn := cnats.GetConnect()
	err := conn.Subscribe(p.remoteTypeSubject, processor)
	if err != nil {
		clog.Errorf("[remoteTypeProcess] Create subscribe failed. [subject = %s, err = %v]", p.remoteTypeSubject, err)
		return
	}
}

func (p *Cluster) PublishLocal(nodeID string, packet *cproto.ClusterPacket) error {
	defer packet.Recycle()

	nodeType, err := p.app.Discovery().GetType(nodeID)
	if err != nil {
		clog.Warnf("[PublishLocal] Get node type fail. [nodeID = %s, pakcet = %s, err = %v],",
			nodeID, packet.PrintLog(), err)

		return cerror.DiscoveryNotFoundNode
	}

	bytes, err := proto.Marshal(packet)
	if err != nil {
		clog.Warnf("[PublishLocal] Marshal error. [nodeID = %s, packet = %s, err = %v]",
			nodeID, packet.PrintLog(), err)

		return cerror.ClusterPatcketMarshalFail
	}

	subject := GetLocalSubject(p.prefix, nodeType, nodeID)
	err = cnats.GetConnect().Publish(subject, bytes)
	if err != nil {
		clog.Warnf("[PublishLocal] Nats publish fail. [nodeID = %s, %s, err = %v]",
			nodeID, packet.PrintLog(), err)

		return cerror.ClusterPublishFail
	}

	return nil
}

func (p *Cluster) PublishRemote(nodeID string, packet *cproto.ClusterPacket) error {
	defer packet.Recycle()

	nodeType, err := p.app.Discovery().GetType(nodeID)
	if err != nil {
		clog.Warnf("[PublishRemote] Get node type fail. [nodeID = %s, %s, err = %v]",
			nodeID, packet.PrintLog(), err)

		return cerror.DiscoveryNotFoundNode
	}

	bytes, err := proto.Marshal(packet)
	if err != nil {
		clog.Warnf("[PublishRemote] Marshal error. [nodeID = %s, packet = %s, err = %vv]",
			nodeID, packet.PrintLog(), err)

		return cerror.ClusterPatcketMarshalFail
	}

	subject := GetRemoteSubject(p.prefix, nodeType, nodeID)
	err = cnats.GetConnect().Publish(subject, bytes)
	if err != nil {
		clog.Warnf("[PublishRemote] Nats publish fail. [nodeID = %s, %s, err = %v]",
			nodeID, packet.PrintLog(), err)

		return cerror.ClusterPublishFail
	}

	return nil
}

func (p *Cluster) PublishRemoteType(nodeType string, packet *cproto.ClusterPacket) error {
	defer packet.Recycle()

	bytes, err := proto.Marshal(packet)
	if err != nil {
		clog.Warnf("[PublishRemoteType] Marshal error. [nodeType = %s, packet = %s, err = %v]",
			nodeType, packet.PrintLog(), err)

		return cerror.ClusterPatcketMarshalFail
	}

	// TODO 往上面移动减少 proto 序列化
	if nodeType == "" {
		return cerror.ClusterNodeTypeIsNil
	}

	if members := p.app.Discovery().ListByType(nodeType); len(members) < 1 {
		return cerror.ClusterNodeTypeMemberNotFound
	}

	subject := GetRemoteTypeSubject(p.prefix, nodeType)
	err = cnats.GetConnect().Publish(subject, bytes)
	if err != nil {
		clog.Warnf("[PublishRemoteType] Nats publish fail. [nodeType = %s, %s, err = %v]",
			nodeType, packet.PrintLog(), err)

		return cerror.ClusterPublishFail
	}

	return nil
}

func (p *Cluster) RequestRemote(nodeID string, packet *cproto.ClusterPacket, timeout ...time.Duration) ([]byte, int32) {
	defer packet.Recycle()

	nodeType, err := p.app.Discovery().GetType(nodeID)
	if err != nil {
		clog.Warnf("[RequestRemote] Get node type fail. [nodeID = %s, %s, err = %v]",
			nodeID, packet.PrintLog(), err)

		return nil, ccode.DiscoveryNotFoundNode
	}

	bytes, err := proto.Marshal(packet)
	if err != nil {
		clog.Warnf("[RequestRemote] Marshal fail. [nodeID = %s, %s, err = %v]",
			nodeID, packet.PrintLog(), err)

		return nil, ccode.RPCMarshalError
	}

	subject := GetRemoteSubject(p.prefix, nodeType, nodeID)
	natsData, err := cnats.GetConnect().RequestSync(subject, bytes, timeout...)
	if err != nil {
		clog.Warnf("[RequestRemote] Nats request fail. [nodeID = %s, %s, err = %v]",
			nodeID, packet.PrintLog(), err)

		return nil, ccode.RPCRemoteExecuteError
	}

	rsp := &cproto.Response{}
	if err = proto.Unmarshal(natsData, rsp); err != nil {
		clog.Warnf("[RequestRemote] Unmarshal fail. [nodeID = %s, %s, rsp, err = %v]",
			nodeID, packet.PrintLog(), rsp, err)
	}

	return rsp.Data, rsp.Code
}
