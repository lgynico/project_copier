package cherryFacade

import (
	"time"

	cproto "github.com/lgynico/project-copier/cherry/net/proto"
)

type (
	ICluster interface {
		Init()
		Stop()

		PublishLocal(nodeID string, packet *cproto.ClusterPacket) error
		PublishRemote(nodeID string, packet *cproto.ClusterPacket) error
		PublishRemoteType(nodeType string, packet *cproto.ClusterPacket) error
		RequestRemote(nodeID string, packet *cproto.ClusterPacket, timeout ...time.Duration) ([]byte, int32)
	}

	IDiscovery interface {
		Name() string
		Load(app IApplication)
		Stop()

		GetMember(nodeID string) (IMember, bool)
		AddMember(member IMember)
		RemoveMember(nodeID string)
		GetType(nodeID string) (nodeType string, err error)

		Map() map[string]IMember
		ListByType(nodeType string, filterNodeID ...string) []IMember
		Random(nodeType string) (IMember, bool)

		OnAddMember(listener MemberListener)
		OnRemoveMember(listener MemberListener)
	}

	IMember interface {
		GetNodeID() string
		GetNodeType() string
		GetAddress() string
		GetSettings() map[string]string
	}

	MemberListener func(member IMember)
)
