package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestZ85DecodeCmd(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty input",
			input: "",
			want:  "",
		},
		{
			name:  "partial chunk 2 chars",
			input: "li",
			want:  "\x42",
		},
		{
			name:  "partial chunk 3 chars",
			input: "lp7",
			want:  "\x42\x42",
		},
		{
			name:  "partial chunk 4 chars",
			input: "lpa0",
			want:  "\x42\x42\x42",
		},
		{
			name:  "exact chunk",
			input: "Hello",
			want:  "\x86\x4F\xD2\x6F",
		},
		{
			name:  "two chunks",
			input: "HelloWorld",
			want:  "\x86\x4F\xD2\x6F\xB5\x59\xF7\x5B",
		},
		{
			name:  "text with partial chunk",
			input: "nm=QNzY<mxA+]nfaP",
			want:  "Hello world!!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", ".")
			cmd.Stdin = strings.NewReader(tt.input)
			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("command failed: %v", err)
			}
			if string(out) != tt.want {
				t.Errorf("got %q, want %q", out, tt.want)
			}
		})
	}
}
