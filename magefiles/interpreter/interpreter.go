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
	f, _ := os.Open(common.CHANGELOG_PATCH)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "\t", 3)
		cmd, deployment, args := parts[0], parts[1], parts[2]
		fmt.Println("Execute: ", cmd, " On ", deployment, ". Args: ", args)
		commands.CMD[cmd].ApplyFromStr(deployment, args)
		changelog.WriteString(line + "\n")
	}
	f.Close()
	//sh.Run("rm", common.CHANGELOG_PATCH)
	fmt.Println("End doPatch")
}

func DoManual() {
	f, _ := os.OpenFile(common.CHANGELOG, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString("manual\n")
	fmt.Println("Added manual record to the Changelog")
}

func WriteChangelogPatchCmd[T any](command string, deployment string, args T) {
	f, _ := os.OpenFile(common.CHANGELOG_PATCH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	argsBytes := commands.SerializeArgs(args)
	f.WriteString(command + "\t" + deployment + "\t" + string(argsBytes) + "\n")
	f.Close()
}
