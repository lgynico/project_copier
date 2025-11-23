package cherryDiscovery

import (
	"math/rand"
	"sync"

	cherryError "github.com/lgynico/project-copier/cherry/error"
	cslice "github.com/lgynico/project-copier/cherry/extend/slice"
	cfacade "github.com/lgynico/project-copier/cherry/facade"
	clog "github.com/lgynico/project-copier/cherry/logger"
	cherryProto "github.com/lgynico/project-copier/cherry/net/proto"
	cprofile "github.com/lgynico/project-copier/cherry/profile"
)

type DiscoveryDefault struct {
	memberMap       sync.Map
	onAddListener   []cfacade.MemberListener
	onRemovListener []cfacade.MemberListener
}

// TODO 好像没什么用
func (p *DiscoveryDefault) PreInit() {
	p.memberMap = sync.Map{}
}

func (p *DiscoveryDefault) Load(_ cfacade.IApplication) {
	nodeConfig := cprofile.GetConfig("node")
	if nodeConfig == nil {
		clog.Error("`node` property not found in profile file.")
		return
	}

	for _, nodeType := range nodeConfig.Keys() {
		typeJson := nodeConfig.GetConfig(nodeType)
		for i := range typeJson.Size() {
			item := typeJson.Get(i)

			nodeID := item.Get("node_id").ToString()
			if nodeID == "" {
				clog.Errorf("nodeID is empty in nodeType = %s", nodeType)
				break
			}

			if _, ok := p.GetMember(nodeID); ok {
				clog.Errorf("nodeType = %s, nodeID = %s, duplicate nodeID", nodeType, nodeID)
				break
			}

			member := &cherryProto.Member{
				NodeID:   nodeID,
				NodeType: nodeType,
				Address:  item.Get("rpc_address").ToString(),
				Settings: make(map[string]string),
			}

			settings := item.Get("__settings__")
			for _, key := range settings.Keys() {
				member.Settings[key] = settings.Get(key).ToString()
			}

			p.memberMap.Store(member.NodeID, member)
		}
	}
}

func (p *DiscoveryDefault) Name() string {
	return "default"
}

func (p *DiscoveryDefault) Map() map[string]cfacade.IMember {
	memberMap := map[string]cfacade.IMember{}
	p.memberMap.Range(func(key, value any) bool {
		if member, ok := value.(cfacade.IMember); ok {
			memberMap[member.GetNodeID()] = member
		}

		return true
	})

	return memberMap
}

func (p *DiscoveryDefault) ListByType(nodeType string, filterNodeID ...string) []cfacade.IMember {
	var memberList []cfacade.IMember
	p.memberMap.Range(func(key, value any) bool {
		if member, ok := value.(cfacade.IMember); ok {
			if _, ok = cslice.StringIn(member.GetNodeID(), filterNodeID); !ok {
				memberList = append(memberList, member)
			}
		}

		return true
	})

	return memberList
}

func (p *DiscoveryDefault) Random(nodeType string) (cfacade.IMember, bool) {
	memberList := p.ListByType(nodeType)
	memberLen := len(memberList)

	if memberLen < 1 {
		return nil, false
	}

	if memberLen == 1 {
		return memberList[0], true
	}

	return memberList[rand.Intn(memberLen)], true
}

func (p *DiscoveryDefault) GetType(nodeID string) (string, error) {
	member, ok := p.GetMember(nodeID)
	if !ok {
		return "", cherryError.Errorf("nodeID = %s not found", nodeID)
	}

	return member.GetNodeType(), nil
}

func (p *DiscoveryDefault) GetMember(nodeID string) (cfacade.IMember, bool) {
	if nodeID == "" {
		return nil, false
	}

	member, ok := p.memberMap.Load(nodeID)
	if !ok {
		return nil, false
	}

	return member.(cfacade.IMember), true
}

func (p *DiscoveryDefault) AddMember(member cfacade.IMember) {
	_, loaded := p.memberMap.LoadOrStore(member.GetNodeID(), member)
	if loaded {
		clog.Warnf("Duplicate nodeID. [nodeType = %s, nodeID = %s, settings = %v]",
			member.GetNodeType(),
			member.GetNodeID(),
			member.GetSettings())
	}

	for _, listener := range p.onAddListener {
		listener(member)
	}

	// TODO 打印占位符不对
	clog.Debugf("AddMember new member. [member = %s]", member)
}

func (p *DiscoveryDefault) RemoveMember(nodeID string) {
	value, loaded := p.memberMap.LoadAndDelete(nodeID)
	if loaded {
		member := value.(cfacade.IMember)
		// TODO 打印占位符不对
		clog.Debugf("remove member. [member = %s]", member)

		for _, listener := range p.onRemovListener {
			listener(member)
		}
	}

}

func (p *DiscoveryDefault) OnAddMember(listener cfacade.MemberListener) {
	if listener == nil {
		return
	}

	p.onAddListener = append(p.onAddListener, listener)
}

func (p *DiscoveryDefault) OnRemoveMember(listener cfacade.MemberListener) {
	if listener == nil {
		return
	}

	p.onRemovListener = append(p.onRemovListener, listener)
}

func (p *DiscoveryDefault) Stop() {

}
