package cherryProfile

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	cfacade "github.com/lgynico/project-copier/cherry/facade"
)

type Config struct {
	jsoniter.Any
}

func Wrap(v any) *Config {
	return &Config{
		Any: jsoniter.Wrap(v),
	}
}

func (p *Config) GetConfig(path ...any) cfacade.ProfileJSON {
	return &Config{
		Any: p.Any.Get(path...),
	}
}

func (p *Config) GetString(path any, defaultValue ...string) string {
	result := p.Get(path)
	if result.LastError() != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return ""
	}

	return result.ToString()
}

func (p *Config) GetBool(path any, defaultValue ...bool) bool {
	result := p.Get(path)
	if result.LastError() != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return false
	}

	return result.ToBool()
}

func (p *Config) GetInt(path any, defaultValue ...int) int {
	result := p.Get(path)
	if result.LastError() != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return 0
	}

	return result.ToInt()
}

func (p *Config) GetInt32(path any, defaultValue ...int32) int32 {
	result := p.Get(path)
	if result.LastError() != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return 0
	}

	return result.ToInt32()
}

func (p *Config) GetInt64(path any, defaultValue ...int64) int64 {
	result := p.Get(path)
	if result.LastError() != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return 0
	}

	return result.ToInt64()
}

func (p *Config) GetDuration(path any, defaultValue ...time.Duration) time.Duration {
	result := p.Get(path)
	if result.LastError() != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return 0
	}

	return time.Duration(result.ToInt64())
}

func (p *Config) Unmarshal(v any) error {
	if p.LastError() != nil {
		return p.LastError()
	}

	return jsoniter.UnmarshalFromString(p.ToString(), v)
}
