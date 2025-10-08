//go:build windows
// +build windows

package main

import (
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// stop stops the command and all its child processes.
func stop(cmd *exec.Cmd) {
	// https://stackoverflow.com/a/44551450
	if cmd.Process == nil {
		return
	}

	exec.Command("taskkill.exe", "/t", "/f", "/pid", strconv.Itoa(cmd.Process.Pid)).Run()

	// Give the process a moment to terminate gracefully
	time.Sleep(100 * time.Millisecond)

	// Check if the process is still running and force kill if necessary
	if isProcessRunningWindows(cmd.Process.Pid) {
		exec.Command("taskkill.exe", "/t", "/f", "/pid", strconv.Itoa(cmd.Process.Pid)).Run()
	}
}

// isProcessRunningWindows checks if a process with the given PID is still running on Windows
func isProcessRunningWindows(pid int) bool {
	// Use tasklist to check if the process is still running
	cmd := exec.Command("tasklist.exe", "/fi", "PID eq "+strconv.Itoa(pid))
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	// If the process is found, tasklist will contain the PID in its output
	return strings.Contains(string(output), strconv.Itoa(pid))
}

// setpgid is a no-op on windows.
func setpgid(cmd *exec.Cmd) {}

// joinArgs joins the arguments of the command into a string which can then be
// passed to `exec.Command("pwsh.exe", "-command", $STRING)`. Examples:
//
// ["echo", "foo"] => echo foo
//
// ["echo", "hello goodbye"] => echo 'hello goodbye'
func joinArgs(args []string) string {
	// references:
	// https://www.rlmueller.net/PowerShellEscape.htm
	// https://stackoverflow.com/a/11231504
	var b strings.Builder
	for i, arg := range args {
		if i == 0 {
			b.WriteString(arg)
			continue
		}
		b.WriteString(" ")
		if arg == "" {
			b.WriteString("''")
			continue
		}
		if !strings.ContainsAny(arg, " '`$(){}<>|&;*") {
			b.WriteString(arg)
			continue
		}
		b.WriteString("'" + strings.ReplaceAll(arg, "'", "''") + "'")
	}
	return b.String()
}
