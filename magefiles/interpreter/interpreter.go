package interpreter

import (
	"fmt"
	"os"
	"strings"

	"github.com/magefile/mage/sh"

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
	sh.Run("rm", common.CHANGELOG_PATCH)
	fmt.Println("End doPatch")
}

func DoManual() {
	f, _ := os.OpenFile(common.CHANGELOG, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString("manual\n")
	fmt.Println("Added manual record to the Changelog")
}
