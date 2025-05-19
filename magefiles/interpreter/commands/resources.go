package commands

import (
	"errors"
	"os"
	"path"

	"github.com/vindex10/devex-thing/magefiles/common"
	yaml "gopkg.in/yaml.v3"
)

type ResourcesSetLimitsArgs struct {
	Cpu    string
	Memory string
}

type ResourcesSetLimits Command[ResourcesSetLimitsArgs]

func (c ResourcesSetLimits) Apply(deployment string, args ResourcesSetLimitsArgs) error {
	//fmt.Println("Apply ", c.Key, " for ", deployment, " with args: ", args.Cpu, "  ; ", args.Ram)
	filename := path.Join(common.DEPLOYMENTS_DIR, deployment, common.DEPLOYMENT_FILE)
	node, readErr := common.ReadYaml(filename)
	if readErr != nil {
		return errors.New("can't read deployment.yaml")
	}
	limits := common.YamlGetPath([]string{"spec", "resources", "limits"}, node.Content[0])
	cpu := common.YamlFindValue("cpu", limits)
	if cpu == nil {
		limits.Content = append(limits.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: "cpu", Tag: "!!str"})
		limits.Content = append(limits.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: args.Cpu, Tag: "!!str"})
	}
	memory := common.YamlFindValue("memory", limits)
	if memory == nil {
		limits.Content = append(limits.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: "memory", Tag: "!!str"})
		limits.Content = append(limits.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: args.Memory, Tag: "!!str"})
	}
	dump, _ := yaml.Marshal(node.Content[0])
	os.WriteFile(filename, dump, os.FileMode(os.O_WRONLY))
	return nil
}

func (c ResourcesSetLimits) ApplyFromStr(deployment string, args string) error {
	return CommandApplyFromStr(c, deployment, args)
}
