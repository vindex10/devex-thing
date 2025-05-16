package commands

import (
	"encoding/json"
)

var KEYS = struct {
	DeploymentInit string
	ResourcesSet   string
}{
	"deployment_init",
	"resources_set",
}

var CMD = map[string]StrCommander{
	KEYS.DeploymentInit: DeploymentInit{Key: KEYS.ResourcesSet},
	KEYS.ResourcesSet:   ResourcesSet{Key: KEYS.ResourcesSet},
}

type Commander[T any] interface {
	Apply(string, T) error
}

type StrCommander interface {
	ApplyFromStr(string, string) error
}

type Command[T any] struct {
	Commander[T]
	StrCommander
	Key string
}

func CommandApplyFromStr[T any](c Commander[T], deployment string, args string) error {
	return c.Apply(deployment, DeserializeArgs[T](args))
}

func SerializeArgs[T any](args T) string {
	res, _ := json.Marshal(args)
	return string(res)
}

func DeserializeArgs[T any](args string) T {
	var r T
	json.Unmarshal([]byte(args), &r)
	return r
}
