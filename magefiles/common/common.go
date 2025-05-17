package common

import (
	"os"
	"path"

	"github.com/magefile/mage/sh"
)

const RELEASE_BRANCH = "main"
const GIT_REMOTE = "origin"
const ARTIFACT_REGISTRY = "europe-north1-docker.pkg.dev/coop-test-459821/prod"
const IMAGE_SOURCE = "https://github.com/vindex10"
const TEMP_DIR = "tmp"

var GIT_ROOT, _ = sh.Output("git", "rev-parse", "--show-toplevel")
var CHANGELOG_PATCH = path.Join(GIT_ROOT, "changelog.patch")
var CHANGELOG = path.Join(GIT_ROOT, "changelog")
var DEPLOYMENTS_DIR = path.Join(GIT_ROOT, "deployments")

func AssembleImageFullName(registry string, name string, version string) string {
	return path.Join(registry, name) + ":" + version
}

func Checked(err error) {
	if err != nil {
		panic("Crash")
	}
}

func MkdirTemp(tag string) (string, error) {
	os.Mkdir(TEMP_DIR, 0755)
	return os.MkdirTemp(TEMP_DIR, tag+"-*")
}
