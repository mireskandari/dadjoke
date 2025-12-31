//go:build !windows

package pdf

import "os/exec"

// hideWindow is a no-op on non-Windows platforms
func hideWindow(cmd *exec.Cmd) {
	// Nothing to do on macOS/Linux
}
