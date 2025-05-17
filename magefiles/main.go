package main

import (
	"github.com/magefile/mage/sh"

	"github.com/vindex10/devex-thing/magefiles/common"
)

func Clean() {
	sh.Rm(common.TEMP_DIR)
}
