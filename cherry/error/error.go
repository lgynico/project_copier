package cherryError

import "errors"

func Error(text string) error {
	return errors.New(text)
}

var (
	ErrActorPath = Error("Actor path is error.")
)

var (
	ProtobufWrongValueType = Error("Convert on wrong value type")
)

var (
	ClusterRequestTimeout = Error("Cluster request timeout")
)
