package system

import (
	"log"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

func runHidden(cmd *exec.Cmd) {
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	cmd.Run()
}

func startHidden(cmd *exec.Cmd) {
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	cmd.Start()
}

func ControlVolume(command string) string {
	var cmd *exec.Cmd

	if strings.Contains(command, "up") || strings.Contains(command, "increase") {
		// Increase volume by 10%
		cmd = exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]175)")
		runHidden(cmd)
		return "Increasing volume"
	} else if strings.Contains(command, "down") || strings.Contains(command, "decrease") {
		// Decrease volume by 10%
		cmd = exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]174)")
		runHidden(cmd)
		return "Decreasing volume"
	} else if strings.Contains(command, "mute") {
		// Toggle mute
		cmd = exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]173)")
		runHidden(cmd)
		return "Toggling mute"
	}

	return "Volume command not recognized"
}

func ControlBrightness(command string) string {
	log.Println("Brightness command:", command)
	return "Adjusting brightness..."
}

func PowerControl(command string) string {
	var cmd *exec.Cmd
	if strings.Contains(command, "shutdown") {
		cmd = exec.Command("shutdown", "/s", "/t", "0")
		startHidden(cmd)
		return "Shutting down the system"
	} else if strings.Contains(command, "restart") {
		cmd = exec.Command("shutdown", "/r", "/t", "0")
		startHidden(cmd)
		return "Restarting the system"
	} else if strings.Contains(command, "lock") {
		cmd = exec.Command("rundll32.exe", "user32.dll,LockWorkStation")
		startHidden(cmd)
		return "Locking the screen"
	} else if strings.Contains(command, "sleep") {
		// Put system to sleep
		cmd = exec.Command("rundll32.exe", "powrprof.dll,SetSuspendState", "0", "1", "0")
		startHidden(cmd)
		return "Putting the system to sleep"
	}
	return "Invalid power command"
}

func GetStatus() map[string]interface{} {
	// Simulated system stats
	return map[string]interface{}{
		"cpu":    12,
		"ram":    48,
		"status": "Online",
	}
}
