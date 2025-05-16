package commands

import (
	"fmt"
	"github.com/vindex10/devex-thing/magefiles/common"
	"strings"
)

type ResourcesSetArgs struct {
	ram string
	cpu string
}
type ResourcesSet common.Command[ResourcesSetArgs]

func (ResourcesSet) ParseArgs(argstr string) ResourcesSetArgs {
	parts := strings.SplitN(argstr, " ", 2)
	return ResourcesSetArgs{parts[0], parts[1]}
}

func (ResourcesSet) Apply(args ResourcesSetArgs) error {
	fmt.Println("Set resources. RAM: ", args.ram, " CPU:", args.cpu)
	return nil
}

func (c ResourcesSet) ApplyFromStr(args string) error {
	return common.CommandApplyFromStr[ResourcesSetArgs](c, args)
}
