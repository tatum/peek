package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runFzf(lang, themeName string, noLines bool) (string, error) {
	if _, err := exec.LookPath("fzf"); err != nil {
		return "", errors.New("fzf not found on PATH")
	}

	self, err := exec.LookPath(os.Args[0])
	if err != nil {
		self = os.Args[0]
	}

	previewArgs := []string{"--force-color"}
	if lang != "" {
		previewArgs = append(previewArgs, "-l", shellQuote(lang))
	}
	if themeName != "" {
		previewArgs = append(previewArgs, "-t", shellQuote(themeName))
	}
	if noLines {
		previewArgs = append(previewArgs, "-n")
	}
	preview := fmt.Sprintf("%s %s {}", shellQuote(self), strings.Join(previewArgs, " "))

	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", fmt.Errorf("open /dev/tty: %w", err)
	}
	defer tty.Close()

	cmd := exec.Command("fzf", "--preview", preview, "--preview-window", "right:60%")
	cmd.Stdin = tty
	cmd.Stderr = tty
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 130 {
			return "", nil // user cancelled
		}
		return "", fmt.Errorf("fzf: %w", err)
	}

	return strings.TrimRight(string(out), "\n"), nil
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
