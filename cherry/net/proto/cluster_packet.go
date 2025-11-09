package cherryProto

import (
	"fmt"
	"sync"

	ctime "github.com/lgynico/project-copier/cherry/extend/time"
	"google.golang.org/protobuf/proto"
)

var clusterPacketPool = &sync.Pool{
	New: func() any {
		return new(ClusterPacket)
	},
}

func GetClusterPacket() *ClusterPacket {
	pkg := clusterPacketPool.Get().(*ClusterPacket)
	pkg.BuildTime = ctime.Now().ToMillisecond()
	return pkg
}

func UnmarshalPacket(data []byte) (*ClusterPacket, error) {
	packet := GetClusterPacket()
	err := proto.Unmarshal(data, packet)
	return packet, err
}

func BuildClusterPacket(source, target, funcName string) *ClusterPacket {
	packet := GetClusterPacket()
	packet.SourcePath = source
	packet.TargetPath = target
	packet.FuncName = funcName

	return packet
}

func (p *ClusterPacket) Recycle() {
	p.BuildTime = 0
	p.SourcePath = ""
	p.TargetPath = ""
	p.FuncName = ""
	p.ArgBytes = nil
	p.Session = nil

	clusterPacketPool.Put(p)
}

func (p *ClusterPacket) PrintLog() string {
	return fmt.Sprintf("source = %s, target = %s, funcName = %s, bytesLen = %d, session = %v",
		p.SourcePath, p.TargetPath, p.FuncName, len(p.ArgBytes), p.Session)
}
