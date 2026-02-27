package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func pagerCommand() (string, []string) {
	pager := os.Getenv("PAGER")
	if pager == "" {
		return "less", []string{"-R"}
	}

	parts := strings.Fields(pager)
	if len(parts) == 1 {
		return parts[0], nil
	}
	return parts[0], parts[1:]
}

func outputWithPager(content string) error {
	cmd, args := pagerCommand()
	pager := exec.Command(cmd, args...)
	pager.Stdout = os.Stdout
	pager.Stderr = os.Stderr

	stdin, err := pager.StdinPipe()
	if err != nil {
		return fmt.Errorf("pager stdin: %w", err)
	}

	if err := pager.Start(); err != nil {
		// Fallback: just print if pager fails
		fmt.Print(content)
		return nil
	}

	_, _ = io.WriteString(stdin, content)
	stdin.Close()

	return pager.Wait()
}
