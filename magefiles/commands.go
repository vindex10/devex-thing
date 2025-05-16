package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/vindex10/devex-thing/magefiles/commands"
)

type Deployment mg.Namespace

// Create new Deployment
func (Deployment) Init(deploymentName string, imageName string, imageVersion string) {
	fmt.Println(deploymentName, commands.KEYS.DeploymentInit)
}

type Resources mg.Namespace

// Replace resources
func (Resources) Set(deploymentName string, ram string, cpu string) {
	fmt.Println(deploymentName, commands.KEYS.ResourcesSet)
}
