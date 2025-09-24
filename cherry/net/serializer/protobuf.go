package cherrySerializer

import (
	cerr "github.com/lgynico/project-copier/cherry/error"
	"google.golang.org/protobuf/proto"
)

type Protobuf struct {
}

func NewProtobuf() *Protobuf {
	return &Protobuf{}
}

func (p *Protobuf) Marshal(v any) ([]byte, error) {
	if data, ok := v.([]byte); ok {
		return data, nil
	}

	if pb, ok := v.(proto.Message); ok {
		return proto.Marshal(pb)
	}

	return nil, cerr.ProtobufWrongValueType
}

func (p *Protobuf) Unmarshal(data []byte, v any) error {
	if pb, ok := v.(proto.Message); ok {
		return proto.Unmarshal(data, pb)
	}

	return cerr.ProtobufWrongValueType
}

func (p *Protobuf) Name() string {
	return "protobuf"
}
