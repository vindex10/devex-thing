# Devex-thing

This repository is an example of GitOps setup for Kubernetes deployments,
and could serve as an interface for interaction between
Dev Teams developing applications and Platform Team maintaining the infrastructure.

It solves two problems:

1. For Developers, who don't need to be experts in Kubernetes to configure their application
2. For Platform Team, who want to save time on reviewing config files focusing on semantical changes not code changes.

Indeed, infrastructure as code is typically quite verbose.
Verbosity brings flexibility but also increases complexity for working with the code base.

The solution presented in this repo introduces a workflow which is aware of a set of basic operations, like
creating a new app, changing its resource limits, deleting an app. The list can be expanded :)

A typical workflow of a developer would be:

1. Create a new branch
2. Stage sequence of basic operations for submission using a user friendly CLI
3. Open pull request

Then it comes to the Platform Team to review the sequence of changes (instead of diffing the yamls), then merge Pull Request.

Automation then does the rest of the processing:

1. Validate and execute the set of operations - this results in applying actual changes to kubernetes yaml configs
2. Checking if requested images exist, and if not - pulling the source code, running tests and releasing an image
3. Applying the changes to the kubernetes cluster
4. Creating a release commit containing the actual diffs to the kubernetes configs
reflecting the current state of the cluster, and updating the changelog.

## Example 1. Interpreted changes

```
mage change:new new_app1  # checkout new branch
mage deployment:init app1 devex-service v1.0.0 '--AppName app123'  # create deployment
mage resources:setLimits app1 100m 2  # memory, cpu set. default is controlled by cluster config
git add changelog.patch
git commit -m "release app1"
mage change:push  # validate, push to origin
```

as a result, a new branch `new_app1` will be created. The changes are written to a `changelog.patch` in the form
(command, deployment, args):

```
deployment_init	app1	{"ImageName":"devex-service","ImageVersion":"v1.0.0","AppName":"app123","Replicas":1,"ContainerPort":8080,"ImageSource":"https://github.com/vindex10/devex-service","ImageRegistry":"europe-north1-docker.pkg.dev/coop-test-459821/prod"}
resources_set_limits	app1	{"Cpu":"2","Memory":"100m"}
```

This is written to file `changelog.patch` and should be pushed for review.

When the PR is merged, and automation would run `mage change:apply '--build --deploy'` to execute the commands in the patch,
the patch will be removed in the release commit.

Developer can also run `mage change:apply` locally.
Without additional flags it is a dry-run version to verify the code changes if needed.


## Example 2. Manual changes

```
mage change:new new_app1
# do the changes
git add deployments/
git commit -m "do manual changes"
mage change:push
```

When you push manual changes to deployments, automation will just apply them.
It is developer's responsibility then to make sure the images are released.

A changelog entry 'manual' will be added in a release commit.

It is not allowed to do both manual and interpreted changes in the same PR.


## Design thoughts


This approach is my first attempt to implememt the workflow involving the interpreted changes.
It was a tracing bullet project for me, to get an overview of possible issues, technical challenges
and generally to experience the problem we are solving.

In this variation, only a list of interpreted changes is pushed for review.

Alternatively, one could push both: `changelog.patch` and the changes applied to the yamls inside `deployments/`.
Then PR automation could validate that the changes in the changelog actually reproduce the changes pushed inside `deployments/`.

I think this could:

* Improve experience for reviewers. Because when approving PR you know exactly what will be deployed without second guesses.
* Reduce risk of race conditions when two PRs are merged and the automations run concurrently.


## What is missing

1. Workflow to validate changes inside pull request can be useful
2. Better command descriptions inside `mage -l`
3. Autocomplete for CLI
4. More commands for the interpreter and better validations inside existing ones
