package cmd

import (
	"bytes"
	"testing"
)

func TestVersion_PrintsVersionString(t *testing.T) {
	SetVersion("dev")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--version"})
	t.Cleanup(func() {
		rootCmd.SetOut(nil)
		rootCmd.SetArgs(nil)
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	want := "Version: dev\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestVersion_ShorthandPrintsVersionString(t *testing.T) {
	SetVersion("dev")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"-v"})
	t.Cleanup(func() {
		rootCmd.SetOut(nil)
		rootCmd.SetArgs(nil)
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	want := "Version: dev\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
