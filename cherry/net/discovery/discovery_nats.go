package cherryDiscovery

import (
	"fmt"
	"time"

	cfacade "github.com/lgynico/project-copier/cherry/facade"
	clog "github.com/lgynico/project-copier/cherry/logger"
	cnats "github.com/lgynico/project-copier/cherry/net/nats"
	cproto "github.com/lgynico/project-copier/cherry/net/proto"
	cprofile "github.com/lgynico/project-copier/cherry/profile"
	"github.com/nats-io/nats.go"
)

type DiscoveryNATS struct {
	DiscoveryDefault
	app               cfacade.IApplication
	thisMember        cfacade.IMember
	masterID          string
	prefix            string
	registerSubject   string
	unregisterSubject string
	addSubject        string
	checkSubject      string
}

func (p *DiscoveryNATS) Name() string {
	return "nats"
}

func (p *DiscoveryNATS) isMaster() bool {
	return p.app.NodeID() == p.masterID
}

func (p *DiscoveryNATS) isClient() bool {
	return p.app.NodeID() != p.masterID
}

func (p *DiscoveryNATS) buildSubject(subject string) string {
	return fmt.Sprintf(subject, p.prefix, p.masterID)
}

func (p *DiscoveryNATS) Load(app cfacade.IApplication) {
	p.DiscoveryDefault.PreInit()
	p.app = app
	p.loadMember()
	p.init()
}

func (p *DiscoveryNATS) thisMemberBytes() []byte {
	memberBytes, err := p.app.Serializer().Marshal(p.thisMember)
	if err != nil {
		clog.Warnf("Marshal member data error. err = %v", err)
		return nil
	}

	return memberBytes
}

func (p *DiscoveryNATS) loadMember() {
	p.thisMember = &cproto.Member{
		NodeID:   p.app.NodeID(),
		NodeType: p.app.NodeType(),
		Address:  p.app.RpcAddress(),
		Settings: make(map[string]string),
	}

	config := cprofile.GetConfig("cluster").GetConfig(p.Name())
	if config.LastError() != nil {
		clog.Fatalf("Nats config parameter not found. err = %v", config.LastError())
	}

	p.prefix = config.GetString("prefix", "node")

	p.masterID = config.GetString("master_node_id")
	if p.masterID == "" {
		clog.Fatalf("Master node id not in config.")
	}
}

func (p *DiscoveryNATS) init() {
	p.registerSubject = p.buildSubject("cherry.%s.discovery.%s.register")
	p.unregisterSubject = p.buildSubject("cherry.%s.discovery.%s.unregister")
	p.addSubject = p.buildSubject("cherry.%s.discovery.%s.addMember")
	p.checkSubject = p.buildSubject("cherry.%s.discovery.%s.check")

	p.subscribe(p.unregisterSubject, func(msg *nats.Msg) {
		var unregisterMember cproto.Member
		err := p.app.Serializer().Unmarshal(msg.Data, &unregisterMember)
		if err != nil {
			// TODO 打印占位符不对
			clog.Warnf("err = %s", err)
			return
		}

		if unregisterMember.NodeID == p.app.NodeID() {
			return
		}

		p.RemoveMember(unregisterMember.NodeID)
	})

	p.serverInit()
	p.clientInit()

	clog.Infof("[Discovery %s] is running.", p.Name())
}

func (p *DiscoveryNATS) serverInit() {
	if !p.isMaster() {
		return
	}

	p.AddMember(p.thisMember)

	p.subscribe(p.registerSubject, func(msg *nats.Msg) {
		newMember := &cproto.Member{}
		err := p.app.Serializer().Unmarshal(msg.Data, newMember)
		if err != nil {
			// TODO 打印占位符不对
			clog.Warnf("IMember Unmarshal [name = %s] error. dataLen = %+v, err = %s",
				p.app.Serializer().Name(),
				len(msg.Data),
				err)
			return
		}

		p.AddMember(newMember)

		memberList := &cproto.MemberList{}
		p.memberMap.Range(func(key, value any) bool {
			protoMember := value.(*cproto.Member)
			if protoMember.NodeID != newMember.NodeID {
				memberList.List = append(memberList.List, protoMember)
			}

			return true
		})

		rspData, err := p.app.Serializer().Marshal(memberList)
		if err != nil {
			// TODO 打印占位符不对
			clog.Warnf("Marshal fail. err = %s", err)
			return
		}

		err = msg.Respond(rspData)
		if err != nil {
			// TODO 打印占位符不对
			clog.Warnf("Respond fail. err = %s", err)
			return
		}

		err = cnats.GetConnect().Publish(p.addSubject, msg.Data)
		if err != nil {
			// TODO 打印占位符不对
			clog.Warnf("Publish fail. err = %s", err)
			return
		}
	})

	p.subscribe(p.checkSubject, func(msg *nats.Msg) {
		msg.Respond(nil)
	})
}

func (p *DiscoveryNATS) clientInit() {
	if !p.isClient() {
		return
	}

	p.subscribe(p.addSubject, func(msg *nats.Msg) {
		addMember := &cproto.Member{}
		err := p.app.Serializer().Unmarshal(msg.Data, addMember)
		if err != nil {
			// TODO 打印占位符不对
			clog.Warnf("err = %s", err)
			return
		}

		if _, ok := p.GetMember(addMember.NodeID); !ok {
			p.AddMember(addMember)
		}
	})

	go p.checkMaster()
}

func (p *DiscoveryNATS) checkMaster() {
	for {
		_, ok := p.GetMember(p.masterID)
		if !ok {
			p.registerToMaster()
		}

		time.Sleep(cnats.ReconnectDelay())
	}
}

func (p *DiscoveryNATS) registerToMaster() {
	natsData, err := cnats.GetConnect().Request(p.registerSubject, p.thisMemberBytes())
	if err != nil {
		// TODO 打印占位符不对
		clog.Warnf("Register node to [master = %s] fail. [err = %s]", p.masterID, err)
		return
	}

	// TODO 打印占位符不对
	clog.Infof("Register node to [master = %s]. [member = %s]", p.masterID, p.thisMember)

	memberList := &cproto.MemberList{}
	err = p.app.Serializer().Unmarshal(natsData, memberList)
	if err != nil {
		// TODO 打印占位符不对
		clog.Warnf("err = %s", err)
		return
	}

	for _, member := range memberList.List {
		p.AddMember(member)
	}
}

func (p *DiscoveryNATS) Stop() {
	err := cnats.GetConnect().Publish(p.unregisterSubject, p.thisMemberBytes())
	if err != nil {
		// TODO 打印占位符不对
		clog.Warnf("Publish fail. err = %s", err)
		return
	}

	clog.Debugf("[NodeID = %s] unregister node to [master = %s]", p.app.NodeID(), p.masterID)
}

func (p *DiscoveryNATS) subscribe(subject string, cb nats.MsgHandler) {
	err := cnats.GetConnect().Subscribe(subject, cb)
	if err != nil {
		// TODO 打印占位符不对
		clog.Warnf("Subscribe fail. err = %s", err)
		return
	}
}
