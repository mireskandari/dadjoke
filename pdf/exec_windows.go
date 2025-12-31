//go:build windows

package pdf

import (
	"os/exec"
	"syscall"
)

// hideWindow configures the command to run without a visible console window
func hideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
}
