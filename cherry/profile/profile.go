package cherryProfile

import (
	"path/filepath"

	cerr "github.com/lgynico/project-copier/cherry/error"
	cfile "github.com/lgynico/project-copier/cherry/extend/file"
	cjson "github.com/lgynico/project-copier/cherry/extend/json"
	cstring "github.com/lgynico/project-copier/cherry/extend/string"
	cfacade "github.com/lgynico/project-copier/cherry/facade"
)

var cfg = &struct {
	profilePath string
	profileName string
	jsonConfig  *Config
	env         string
	debug       bool
	printLevel  string
}{}

func Path() string {
	return cfg.profilePath
}

func Name() string {
	return cfg.profileName
}

func Env() string {
	return cfg.env
}

func Debug() bool {
	return cfg.debug
}

func PrintLevel() string {
	return cfg.printLevel
}

// TODO 这个方法应该把配置初始化和获取节点分开
func Init(filePath, nodeID string) (cfacade.INode, error) {
	if filePath == "" {
		return nil, cerr.Error("File path is nil.")
	}

	if nodeID == "" {
		return nil, cerr.Error("NodeID is nil.")
	}

	judgePath, ok := cfile.JudgeFile(filePath)
	if !ok {
		return nil, cerr.Errorf("File path error. filePath = %s", filePath)
	}

	p, f := filepath.Split(judgePath)
	jsonConfig, err := loadFile(p, f)
	if err != nil || jsonConfig.Any == nil || jsonConfig.LastError() != nil {
		return nil, cerr.Errorf("Load profile file error. [err = %v]", err)
	}

	node, err := GetNodeWithConfig(jsonConfig, nodeID)
	if err != nil {
		return nil, cerr.Errorf("Failed to get node config from profile file. [err = %v]", err)
	}

	cfg.profilePath = p
	cfg.profileName = f
	cfg.jsonConfig = jsonConfig
	cfg.env = jsonConfig.GetString("env", "default")
	cfg.debug = jsonConfig.GetBool("debug", true)
	cfg.printLevel = jsonConfig.GetString("print_level", "debug")

	return node, nil
}

func GetConfig(path ...any) cfacade.ProfileJSON {
	return cfg.jsonConfig.GetConfig(path...)
}

func loadFile(filePath, fileName string) (*Config, error) {
	var (
		profileMaps = make(map[string]any)
		includeMaps = make(map[string]any)
		rootMaps    = make(map[string]any)
	)

	fileNamePath := filepath.Join(filePath, fileName)
	if err := cjson.ReadMaps(fileNamePath, profileMaps); err != nil {
		return nil, err
	}

	if v, found := profileMaps["include"].([]any); found {
		paths := cstring.ToStringSlice(v)
		for _, p := range paths {
			includePath := filepath.Join(filePath, p)
			if err := cjson.ReadMaps(includePath, includeMaps); err != nil {
				return nil, err
			}
		}
	}

	mergeMap(rootMaps, includeMaps)
	mergeMap(rootMaps, profileMaps)

	return Wrap(rootMaps), nil
}

// TODO 快乐路径
func mergeMap(dst, src map[string]any) {
	for key, value := range src {
		if v, ok := dst[key]; ok {
			if m1, ok := v.(map[string]any); ok {
				if m2, ok := value.(map[string]any); ok {
					mergeMap(m1, m2)
				} else {
					dst[key] = value
				}
			} else {
				dst[key] = value
			}
		} else {
			dst[key] = value
		}
	}
}
