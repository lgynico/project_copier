package cherryError

import (
	"errors"
	"fmt"
)

func Error(text string) error {
	return errors.New(text)
}

func Errorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

var (
	RouteFieldCantEmpty = Error("Route field can not be empty")
	RouteInvalid        = Error("Invalid route")
)

var (
	PacketWrongType              = Error("Wrong packet type")
	PacketSizeExceed             = Error("Codec: packet size exceed")
	PacketConnectClosed          = Error("Client connection closed")
	PakcetInvalidHeader          = Error("Invalid header")
	PakcetMsgSmallerThanExpected = Error("Received less data than expected, EOF?")
)

var (
	MessageWrongType     = Error("Wrong message type")
	MessageInvalid       = Error("Invalid message")
	MessageRouteNotFound = Error("Route info not found in dictionary")
)

var (
	ProtobufWrongValueType = Error("Convert on wrong value type")
)
var (
	ClusterClientIsStop           = Error("Cluster client is stop")
	ClusterRequestTimeout         = Error("Cluster request timeout")
	ClusterPatcketMarshalFail     = Error("Cluster packet marshal fail")
	ClusterPacketUnmarshalFail    = Error("Cluster packet unmarshal fail")
	ClusterPublishFail            = Error("Cluster publish fail")
	ClusterRequesFail             = Error("Cluster request fail")
	ClusterNodeTypeIsNil          = Error("Cluster node type is nil")
	ClusterNodeTypeMemberNotFound = Error("Cluster node type member not found")
)

var (
	DiscoveryNotFoundNode = Error("Discovery not found node")
)

var (
	ActorPathError = Error("Actor path is error.")
)

var (
	FuncIsNil     = Error("Func is nil")
	FuncTypeError = Error("Func type error")
)
