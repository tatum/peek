package main

import "testing"

func TestPagerCommand(t *testing.T) {
	cmd, args := pagerCommand()
	if cmd == "" {
		t.Error("expected non-empty pager command")
	}
	_ = args // just verify it doesn't panic
}

func TestPagerCommandFromEnv(t *testing.T) {
	t.Setenv("PAGER", "more")
	cmd, args := pagerCommand()
	if cmd != "more" {
		t.Errorf("expected 'more' from PAGER env, got %s", cmd)
	}
	if len(args) != 0 {
		t.Errorf("expected no args for 'more', got %v", args)
	}
}

func TestPagerCommandDefault(t *testing.T) {
	t.Setenv("PAGER", "")
	cmd, args := pagerCommand()
	if cmd != "less" {
		t.Errorf("expected 'less' as default, got %s", cmd)
	}
	found := false
	for _, a := range args {
		if a == "-R" {
			found = true
		}
	}
	if !found {
		t.Error("expected -R flag for less")
	}
}
