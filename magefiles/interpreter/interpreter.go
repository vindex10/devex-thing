package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vindex10/devex-thing/magefiles/common"
	"github.com/vindex10/devex-thing/magefiles/interpreter/commands"
)

func DoPatch() {
	fmt.Println("Begin doPatch")
	changelog, _ := os.OpenFile(common.CHANGELOG, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer changelog.Close()
	for onePatch := range iterChangelogPatch() {
		fmt.Println("Execute: ", onePatch.Command, " On ", onePatch.Deployment, ". Args: ", onePatch.Args)
		commands.CMD[onePatch.Command].ApplyFromStr(onePatch.Deployment, onePatch.Args)
		line := strings.Join([]string{onePatch.Command, onePatch.Deployment, onePatch.Args}, "\t")
		changelog.WriteString(line + "\n")
	}
	//sh.Run("rm", common.CHANGELOG_PATCH)
	fmt.Println("End doPatch")
}

func DoManual() {
	f, _ := os.OpenFile(common.CHANGELOG, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString("manual\n")
	fmt.Println("Added manual record to the Changelog")
}

type ChangelogPatchRaw = struct {
	Command    string
	Deployment string
	Args       string
}

func iterChangelogPatch() func(func(ChangelogPatchRaw) bool) {
	return func(yield func(ChangelogPatchRaw) bool) {
		f, _ := os.Open(common.CHANGELOG_PATCH)
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, "\t", 3)
			cmd, deployment, args := parts[0], parts[1], parts[2]
			if !yield(ChangelogPatchRaw{Command: cmd, Deployment: deployment, Args: args}) {
				return
			}
		}
	}
}

func WriteChangelogPatchCmd[T any](command string, deployment string, args T) {
	f, _ := os.OpenFile(common.CHANGELOG_PATCH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	argsBytes := commands.SerializeArgs(args)
	f.WriteString(command + "\t" + deployment + "\t" + string(argsBytes) + "\n")
	f.Close()
}

type ImageSpec = struct {
	Name    string
	Version string
	Source  string
}

func GetNewImages() []ImageSpec {
	res := []ImageSpec{}
	for onePatch := range iterChangelogPatch() {
		if onePatch.Command == commands.KEYS.DeploymentInit {
			args := commands.DeserializeArgs[commands.DeploymentInitArgs](onePatch.Args)
			if args.ImageRegistry != common.ARTIFACT_REGISTRY {
				continue
			}
			res = append(res, ImageSpec{Name: args.ImageName, Version: args.ImageVersion, Source: args.ImageSource})
		}
	}
	return res
}
