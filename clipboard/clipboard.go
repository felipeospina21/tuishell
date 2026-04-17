package clipboard

import (
	"os/exec"
	"runtime"
	"strings"
)

// CopyToClipboard copies text to the system clipboard.
func CopyToClipboard(text string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("clip")
	default:
		return
	}

	cmd.Stdin = strings.NewReader(text)
	_ = cmd.Run()
}
