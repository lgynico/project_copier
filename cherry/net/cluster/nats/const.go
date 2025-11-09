package cherryNatsCluster

import "fmt"

const (
	localSubjectFormat      = "cherry-%s.local.%s.%s"   // cherry-{env}.local.{nodeType}.{nodeID}
	remoteSubjectFormat     = "cherry-%s.remote.%s.%s"  // cherry-{env}.remote.{nodeType}.{nodeID}
	remoteTypeSubjectFormat = "cherry-%s.remoteType.%s" // cherry-{env}.remoteType.{nodeType}
	replySubjectFormat      = "cherry-%s.reply.%s.%s"   // cherry-{env}.reply.{nodeType}.{nodeID}
)

func GetLocalSubject(prefix, nodeType, nodeID string) string {
	return fmt.Sprintf(localSubjectFormat, prefix, nodeType, nodeID)
}

func GetRemoteSubject(prefix, nodeType, nodeID string) string {
	return fmt.Sprintf(remoteSubjectFormat, prefix, nodeType, nodeID)
}

func GetRemoteTypeSubject(prefix, nodeType string) string {
	return fmt.Sprintf(remoteTypeSubjectFormat, prefix, nodeType)
}

func GetReplySubject(prefix, nodeType, nodeID string) string {
	return fmt.Sprintf(replySubjectFormat, prefix, nodeType, nodeID)
}
