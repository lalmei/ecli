package core

import (
	"os/exec"
)

func Open(url string) error {
	cc := exec.Command("open", url)
	return cc.Run()
}
