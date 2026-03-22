package cmd

import "testing"

func TestInitCommandRegistered(t *testing.T) {
	t.Parallel()

	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "init" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected init command to be registered on root")
	}
}
