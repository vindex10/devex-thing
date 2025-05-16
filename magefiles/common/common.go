package common

import (
	"path"

	"github.com/magefile/mage/sh"
)

const RELEASE_BRANCH = "main"
const GIT_REMOTE = "origin"
const ARTIFACT_REGISTRY = "europe-north1-docker.pkg.dev/coop-test-459821/prod"

var GIT_ROOT, _ = sh.Output("git", "rev-parse", "--show-toplevel")
var CHANGELOG_PATCH = path.Join(GIT_ROOT, "changelog.patch")
var CHANGELOG = path.Join(GIT_ROOT, "changelog")
var DEPLOYMENTS_DIR = path.Join(GIT_ROOT, "deployments")
