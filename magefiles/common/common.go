package common

import (
	"io"
	"os"
	"path"

	"github.com/magefile/mage/sh"
	yaml "gopkg.in/yaml.v3"
)

const RELEASE_BRANCH = "main"
const GIT_REMOTE = "origin"
const ARTIFACT_REGISTRY = "europe-north1-docker.pkg.dev/coop-test-459821/prod"
const IMAGE_SOURCE = "https://github.com/vindex10"
const TEMP_DIR = "tmp"
const NAMESPACE = "default"

var GIT_ROOT, _ = sh.Output("git", "rev-parse", "--show-toplevel")
var CHANGELOG_PATCH = path.Join(GIT_ROOT, "changelog.patch")
var CHANGELOG = path.Join(GIT_ROOT, "changelog")
var DEPLOYMENTS_DIR = path.Join(GIT_ROOT, "deployments")
var DEPLOYMENT_FILE = "deployment.yaml"

func AssembleImageFullName(registry string, name string, version string) string {
	return path.Join(registry, name) + ":" + version
}

func Checked(err error) {
	if err != nil {
		panic("Crash")
	}
}

func MkdirTemp(tag string) (string, error) {
	os.Mkdir(TEMP_DIR, 0755)
	return os.MkdirTemp(TEMP_DIR, tag+"-*")
}

func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func ReadYaml(filepath string) (yaml.Node, error) {
	var res yaml.Node
	dat, err := os.ReadFile(filepath)
	if err != nil {
		return res, err
	}
	yaml.Unmarshal(dat, &res)
	return res, nil
}

func YamlFindValue(value string, node *yaml.Node) *yaml.Node {
	var buf string
	for _, v := range node.Content {
		if v.Tag == "!!map" && buf == value {
			return v
		}
		buf = v.Value
	}
	for _, v := range node.Content {
		if v.Tag != "!!map" {
			continue
		}
		if child := YamlFindValue(value, v); child != nil {
			return child
		}
	}
	return nil
}

func YamlFindPath(nodepath []string, doc *yaml.Node) *yaml.Node {
	res := doc
	for _, key := range nodepath {
		res = YamlFindValue(key, res)
		if res == nil {
			return nil
		}
	}
	return res
}

func YamlGetPath(nodepath []string, doc *yaml.Node) *yaml.Node {
	res := doc
	for _, key := range nodepath {
		newRes := YamlFindValue(key, res)
		if newRes == nil {
			newNode := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
			res.Content = append(res.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key})
			res.Content = append(res.Content, newNode)
			newRes = newNode
		}
		res = newRes
	}
	return res
}
