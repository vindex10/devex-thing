package main

import (
	"flag"
	"fmt"
	"os"
	"path"

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
	build := flags.Bool("build", false, "Build images before deploying")
	argsArr, _ := shlex.Split(args)
	flags.Parse(argsArr)

	Change{}.Validate()
	newImages := interpreter.GetNewImages()
	fmt.Println(newImages)

	if hasPatch() {
		interpreter.DoPatch()
	} else if hasManual() {
		interpreter.DoManual()
	} else {
		fmt.Println("No infrastructure changes introduced. Exit.")
		return
	}

	if *build {
		for _, newImage := range newImages {
			fmt.Println("Building: "+newImage.Name+":"+newImage.Version, " from ", newImage.Source)
			buildImage(newImage)
		}
	}

	if *deploy {
		emptyDeployments, err := common.IsDirEmpty(common.DEPLOYMENTS_DIR)
		if !emptyDeployments {
			sh.RunV("kubectl", "apply", "--prune", "--all", "-f", common.DEPLOYMENTS_DIR, "-R", "--namespace", common.NAMESPACE)
			fmt.Println("Deploy!")
		} else if emptyDeployments && err == nil {
			fmt.Println("Empty deployments dir. Cleaning up")
			sh.RunV("kubectl", "delete", "--all", "deployments", "--namespace", common.NAMESPACE)
		}
	}
}

func buildImage(image interpreter.ImageSpec) error {
	imageFullName := common.AssembleImageFullName(common.ARTIFACT_REGISTRY, image.Name, image.Version)
	imageExistsErr := sh.Run("docker", "manifest", "inspect", imageFullName)
	imageExists := sh.ExitStatus(imageExistsErr) == 0
	if imageExists {
		fmt.Println("Image exists. Skip build. " + imageFullName)
		return nil
	}
	fmt.Println("Image is not available in the registry. Building. " + imageFullName)
	repoDir, _ := common.MkdirTemp(image.Name)
	common.Checked(sh.RunV("git", "clone", image.Source, repoDir))
	os.Chdir(repoDir)
	defer os.Chdir(common.GIT_ROOT)
	common.Checked(sh.RunV("git", "reset", image.Version))
	testsErr := sh.RunV("go", "test", "./...")
	if testsErr != nil {
		fmt.Println("Tests failed for ", imageFullName)
		return testsErr
	}
	common.Checked(sh.RunV("docker", "build", "-f", path.Join(common.GIT_ROOT, "magefiles/assets/Dockerfile"), "-t", imageFullName, "."))
	common.Checked(sh.RunV("docker", "push", imageFullName))
	return nil
}

// -Validate Change for consistency
func (Change) Validate() error {
	os.Chdir(common.GIT_ROOT)
	if hasPatch() && hasManual() {
		return mg.Fatal(1, "Changelog patch can't be used together with the manual changes. Please split into different branches.")
	}
	return nil
}

func hasPatch() bool {
	var base string
	if onReleaseBranch() {
		base = "HEAD^"
	} else {
		// On PR
		base = common.RELEASE_BRANCH
	}
	_, has_patch_err := sh.Output("git", "diff", "--exit-code", base, "--", common.CHANGELOG_PATCH)
	has_patch := sh.ExitStatus(has_patch_err)
	return (has_patch != 0)
}

func hasManual() bool {
	var base string
	if onReleaseBranch() {
		base = "HEAD^"
	} else {
		// On PR
		base = common.RELEASE_BRANCH
	}
	_, has_manual_err := sh.Output("git", "diff", "--exit-code", base, "--", common.DEPLOYMENTS_DIR)
	has_manual := sh.ExitStatus(has_manual_err)
	return (has_manual != 0)
}

func onReleaseBranch() bool {
	currentBranch, _ := sh.Output("git", "rev-parse", "--abbrev-ref", "HEAD")
	return currentBranch == common.RELEASE_BRANCH
}
