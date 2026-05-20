package system

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

func OpenApp(appName string) string {
	var cmd *exec.Cmd

	appName = strings.ToLower(appName)

	if runtime.GOOS == "windows" {
		switch appName {
		case "chrome", "browser":
			cmd = exec.Command("cmd", "/c", "start chrome")
		case "vscode", "vs code", "code":
			cmd = exec.Command("cmd", "/c", "code")
		case "notepad":
			cmd = exec.Command("cmd", "/c", "start notepad")
		case "calculator", "calc":
			cmd = exec.Command("cmd", "/c", "start calc")
		case "file explorer", "explorer":
			cmd = exec.Command("cmd", "/c", "start explorer")
		case "paint":
			cmd = exec.Command("cmd", "/c", "start mspaint")
		case "camera":
			cmd = exec.Command("cmd", "/c", "start microsoft.windows.camera:")
		case "settings":
			cmd = exec.Command("cmd", "/c", "start ms-settings:")
		case "task manager":
			cmd = exec.Command("cmd", "/c", "start taskmgr")
		default:
			// Fallback: Use PowerShell to search for a custom shortcut in Start Menu / Desktop
			psScript := fmt.Sprintf(`
$appName = "%s"
$paths = @(
    "$env:ProgramData\Microsoft\Windows\Start Menu\Programs",
    "$env:APPDATA\Microsoft\Windows\Start Menu\Programs",
    "$env:PUBLIC\Desktop",
    "$env:USERPROFILE\Desktop"
)
foreach ($path in $paths) {
    if (Test-Path $path) {
        $shortcut = Get-ChildItem -Path $path -Filter "*$appName*.lnk" -Recurse -ErrorAction SilentlyContinue | Select-Object -First 1
        if ($shortcut) {
            Start-Process $shortcut.FullName
            exit 0
        }
    }
}
# If we didn't find a shortcut, fallback to trying the name directly
Start-Process $appName -ErrorAction SilentlyContinue
`, appName)

			cmd = exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", psScript)
		}

		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	} else {
		return "App opening is currently optimized for Windows."
	}

	err := cmd.Start()
	if err != nil {
		log.Printf("Failed to open %s: %v", appName, err)
		return "Failed to open " + appName
	}
	return "Opening " + appName
}
