package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestZ85EncodeCmd(t *testing.T) {
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
			name:  "single byte",
			input: "\x42",
			want:  "li",
		},
		{
			name:  "two bytes",
			input: "\x42\x42",
			want:  "lp7",
		},
		{
			name:  "three bytes",
			input: "\x42\x42\x42",
			want:  "lpa0",
		},
		{
			name:  "exact chunk",
			input: "\x86\x4F\xD2\x6F",
			want:  "Hello",
		},
		{
			name:  "two chunks",
			input: "\x86\x4F\xD2\x6F\xB5\x59\xF7\x5B",
			want:  "HelloWorld",
		},
		{
			name:  "text with partial chunk",
			input: "Hello world!!",
			want:  "nm=QNzY<mxA+]nfaP",
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
