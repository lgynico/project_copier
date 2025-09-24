package cherrySerializer

import jsoniter "github.com/json-iterator/go"

type JSON struct {
}

func NewJSON() *JSON {
	return &JSON{}
}

func (p *JSON) Marshal(v any) ([]byte, error) {
	if data, ok := v.([]byte); ok {
		return data, nil
	}

	return jsoniter.Marshal(v)
}

func (p *JSON) Unmarshal(data []byte, v any) error {
	return jsoniter.Unmarshal(data, v)
}

func (p *JSON) Name() string {
	return "json"
}
