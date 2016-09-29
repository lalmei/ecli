package core

import (
	"os/exec"
)

func Open(url string) error {
	cc := exec.Command("open", url)
	if err := cc.Run(); err != nil {
		return err
	}
	return nil
}
