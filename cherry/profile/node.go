package cherryProfile

import (
	"fmt"
	"regexp"
	"strings"

	cerr "github.com/lgynico/project-copier/cherry/error"
	cfacade "github.com/lgynico/project-copier/cherry/facade"
)

type Node struct {
	nodeID     string
	nodeType   string
	address    string
	rpcAddress string
	settings   cfacade.ProfileJSON
	enabled    bool
}

func (p *Node) NodeID() string {
	return p.nodeID
}

func (p *Node) NodeType() string {
	return p.nodeType
}

func (p *Node) Address() string {
	return p.address
}

func (p *Node) RpcAddress() string {
	return p.rpcAddress
}

func (p *Node) Settings() cfacade.ProfileJSON {
	return p.settings
}

func (p *Node) Enabled() bool {
	return p.enabled
}

func (p *Node) String() string {
	return fmt.Sprintf("nodeID = %s, nodeType = %s, address = %s, rpcAddress = %s, enabled = %v",
		p.nodeID, p.nodeType, p.address, p.rpcAddress, p.enabled)
}

func GetNodeWithConfig(config *Config, nodeID string) (cfacade.INode, error) {
	nodeConfig := config.GetConfig("node")
	if nodeConfig.LastError() != nil {
		return nil, cerr.Error("`nodes` property not found in profile file.")
	}

	for _, nodeType := range nodeConfig.Keys() {
		typeJson := nodeConfig.GetConfig(nodeType)
		for i := range typeJson.Size() {
			item := typeJson.GetConfig(i)
			if !findNodeID(nodeID, item.GetConfig("node_id")) {
				continue
			}

			node := &Node{
				nodeID:     nodeID,
				nodeType:   nodeType,
				address:    item.GetString("address"),
				rpcAddress: item.GetString("rpc_address"),
				settings:   item.GetConfig("__settings__"),
				enabled:    item.GetBool("enabled"),
			}

			return node, nil
		}
	}

	return nil, cerr.Errorf("nodeID = %s not found.", nodeID)
}

func LoadNode(nodeID string) (cfacade.INode, error) {
	return GetNodeWithConfig(cfg.jsonConfig, nodeID)
}

func findNodeID(nodeID string, nodeIDJson cfacade.ProfileJSON) bool {
	configNodeID := nodeIDJson.ToString()
	if configNodeID == nodeID {
		return true
	}

	if isRegexNodeID(nodeID, configNodeID) {
		return true
	}

	for i := range nodeIDJson.Size() {
		if nodeIDJson.GetString(i) == nodeID {
			return true
		}
	}

	return false
}

func isRegexNodeID(nodeID, regexNodeID string) bool {
	if !strings.HasPrefix(regexNodeID, "^") || !strings.HasSuffix(regexNodeID, "$") {
		return false
	}

	regex, err := regexp.Compile(regexNodeID)
	if err != nil {
		return false
	}

	return regex.MatchString(nodeID)
}
