package commands

import (
	"fmt"
	"github.com/vindex10/devex-thing/magefiles/common"
	"strings"
)

type DeploymentInitArgs struct {
	imageName    string
	imageVersion string
}
type DeploymentInit common.Command[DeploymentInitArgs]

func (DeploymentInit) ParseArgs(argstr string) DeploymentInitArgs {
	parts := strings.SplitN(argstr, " ", 2)
	return DeploymentInitArgs{parts[0], parts[1]}
}

func (DeploymentInit) Apply(args DeploymentInitArgs) error {
	fmt.Println(args.imageName, args.imageVersion)
	return nil
}

func (c DeploymentInit) ApplyFromStr(args string) error {
	return common.CommandApplyFromStr[DeploymentInitArgs](c, args)
}
