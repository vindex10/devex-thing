package commands

import (
	"fmt"
)

type DeploymentInitArgs struct {
	ImageName    string
	ImageVersion string
}
type DeploymentInit Command[DeploymentInitArgs]

func (c DeploymentInit) Apply(deployment string, args DeploymentInitArgs) error {
	fmt.Println("Apply ", c.Key, " to ", deployment, " with args: ", args.ImageName+":"+args.ImageVersion)
	return nil
}

func (c DeploymentInit) ApplyFromStr(deployment string, args string) error {
	return CommandApplyFromStr(c, deployment, args)
}
