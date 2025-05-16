package main

import (
	"github.com/magefile/mage/mg"
	"github.com/vindex10/devex-thing/magefiles/interpreter"
	"github.com/vindex10/devex-thing/magefiles/interpreter/commands"
)

type Deployment mg.Namespace

// -Create new Deployment
func (Deployment) Init(deployment string, imageName string, imageVersion string) {
	interpreter.WriteChangelogPatchCmd(commands.KEYS.DeploymentInit, deployment, commands.DeploymentInitArgs{ImageName: imageName, ImageVersion: imageVersion})
}

type Resources mg.Namespace

// -Set resources
func (Resources) Set(deployment string, ram string, cpu string) {
	interpreter.WriteChangelogPatchCmd(commands.KEYS.ResourcesSet, deployment, commands.ResourcesSetArgs{Ram: ram, Cpu: cpu})
}
