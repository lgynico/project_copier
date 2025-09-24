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
	ErrActorPath = Error("Actor path is error.")
)

var (
	ProtobufWrongValueType = Error("Convert on wrong value type")
)

var (
	ClusterRequestTimeout = Error("Cluster request timeout")
)
