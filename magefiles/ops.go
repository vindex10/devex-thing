package main

import (
	"flag"
	"strings"

	"github.com/google/shlex"
	"github.com/magefile/mage/mg"

	"github.com/vindex10/devex-thing/magefiles/interpreter"
	"github.com/vindex10/devex-thing/magefiles/interpreter/commands"
)

type Deployment mg.Namespace

// -Create new Deployment
func (Deployment) Init(deployment string, ImageName string, ImageVersion string, rest string) {
	ImageRegistry, ImageName := parseImageRegistry(ImageName)
	args := commands.NewDeploymentInitArgs(ImageName, ImageVersion)
	flags := flag.NewFlagSet("deployment:init", flag.PanicOnError)
	AppName := flags.String("app-name", args.AppName, "")
	Replicas := flags.Int("replicas", args.Replicas, "")
	ContainerPort := flags.Int("container-port", args.ContainerPort, "")
	ImageSource := flags.String("image-source", args.ImageSource, "Where to pull src from")
	argsArr, _ := shlex.Split(rest)
	flags.Parse(argsArr)
	args.AppName = *AppName
	args.Replicas = *Replicas
	args.ContainerPort = *ContainerPort
	args.ImageSource = *ImageSource
	if len(ImageRegistry) > 0 {
		args.ImageRegistry = ImageRegistry
	}
	interpreter.WriteChangelogPatchCmd(commands.KEYS.DeploymentInit, deployment, args)
}

func parseImageRegistry(imageName string) (string, string) {
	if !strings.Contains(imageName, "/") {
		return "", imageName
	}
	parts := strings.Split(imageName, "/")
	imageRegistry := strings.Join(parts[:len(parts)-1], "/") + "/"
	return imageRegistry, parts[len(parts)-1]
}

func (Deployment) Delete(deployment string) {
	interpreter.WriteChangelogPatchCmd(commands.KEYS.DeploymentDelete, deployment, commands.DummyArgs{})
}

type Resources mg.Namespace

// -Set resources
func (Resources) SetLimits(deployment string, memory string, cpu string) {
	interpreter.WriteChangelogPatchCmd(commands.KEYS.ResourcesSetLimits, deployment, commands.ResourcesSetLimitsArgs{Memory: memory, Cpu: cpu})
}
