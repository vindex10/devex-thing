package commands

import (
	"fmt"
)

type ResourcesSetArgs struct {
	Ram string
	Cpu string
}

type ResourcesSet Command[ResourcesSetArgs]

func (c ResourcesSet) Apply(deployment string, args ResourcesSetArgs) error {
	fmt.Println("Apply ", c.Key, " for ", deployment, " with args: ", args.Cpu, "  ; ", args.Ram)
	return nil
}

func (c ResourcesSet) ApplyFromStr(deployment string, args string) error {
	return CommandApplyFromStr(c, deployment, args)
}
