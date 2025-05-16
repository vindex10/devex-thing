package main

import (
	"flag"
	"fmt"

	"github.com/google/shlex"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/vindex10/devex-thing/magefiles/common"
	"github.com/vindex10/devex-thing/magefiles/interpreter"
)

type Change mg.Namespace

// -Prepare branch for new change
// Usage: change:new target-branch
func (Change) New(branch string) {
	// validate repo is clean before introducing new changes
	// validate branch name here
	sh.Run("git", "checkout", "-b", branch)
}

// -Push change for review
func (Change) Push() {
	Change{}.Validate()
	sh.Run("git", "push", "-u", common.GIT_REMOTE, "HEAD")
}

// -Apply change locally and optionally send to the cloud
func (Change) Apply(args string) {
	flags := flag.NewFlagSet("change:apply", flag.PanicOnError)
	deploy := flags.Bool("deploy", false, "Deploy changes to kubernetes cluster")
	argsArr, _ := shlex.Split(args)
	flags.Parse(argsArr)

	Change{}.Validate()
	if hasPatch() {
		interpreter.DoPatch()
	}
	if hasManual() {
		interpreter.DoManual()
	}

	if *deploy {
		//sh.RunV("kubectl", "apply", "--prune", "-f", common.DEPLOYMENTS_DIR, "-R")
		fmt.Println("Deploy!")
	}
}

// -Validate Change for consistency
func (Change) Validate() error {
	sh.Run("cd", common.GIT_ROOT)
	has_patch := hasPatch()
	has_manual := hasManual()
	if has_patch && has_manual {
		return mg.Fatal(1, "Changelog patch can't be used together with the manual changes. Please split into different branches.")
	}
	return nil
}

func hasPatch() bool {
	_, has_patch_err := sh.Output("git", "diff", "--exit-code", common.RELEASE_BRANCH, "--", common.CHANGELOG_PATCH)
	has_patch := sh.ExitStatus(has_patch_err)
	return (has_patch != 0)
}

func hasManual() bool {
	_, has_manual_err := sh.Output("git", "diff", "--exit-code", common.RELEASE_BRANCH, "--", common.DEPLOYMENTS_DIR)
	has_manual := sh.ExitStatus(has_manual_err)
	return (has_manual != 0)
}
