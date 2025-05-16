package commands

import (
	"github.com/vindex10/devex-thing/magefiles/common"
)

var KEYS = struct {
	DeploymentInit string
	ResourcesSet   string
}{
	"deployment_init",
	"resources_set",
}

var CMD = map[string]common.StrCommander{
	KEYS.DeploymentInit: DeploymentInit{Key: KEYS.ResourcesSet},
	KEYS.ResourcesSet:   ResourcesSet{Key: KEYS.ResourcesSet},
}
