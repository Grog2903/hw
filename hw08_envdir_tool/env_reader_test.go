package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessEnvFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "envdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	tests := []struct {
		fileName           string
		content            string
		expectedValue      string
		expectedNeedRemove bool
	}{
		{"empty_file", "", "", true},
		{"simple_file", "HELLO", "HELLO", false},
		{"file_with_spaces", "WORLD  \t", "WORLD", false},
		{"file_with_null", "FOO\x00BAR", "FOO\nBAR", false},
	}

	for _, test := range tests {
		filePath := filepath.Join(dir, test.fileName)
		if err := os.WriteFile(filePath, []byte(test.content), 0o644); err != nil {
			t.Fatalf("Failed to write file %s: %v", test.fileName, err)
		}

		value, needRemove, err := processEnvFile(dir, test.fileName)
		if err != nil {
			t.Errorf("Unexpected error for file %s: %v", test.fileName, err)
		}
		if value != test.expectedValue {
			t.Errorf("For file %s, expected value %q, got %q", test.fileName, test.expectedValue, value)
		}
		if needRemove != test.expectedNeedRemove {
			t.Errorf("For file %s, expected needRemove %v, got %v", test.fileName, test.expectedNeedRemove, needRemove)
		}
	}
}

func TestReadDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "envdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	_ = os.WriteFile(filepath.Join(dir, "VAR1"), []byte("VALUE1"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "VAR2"), []byte("VALUE2  \t"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "EMPTY_VAR"), []byte(""), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "NULL_VAR"), []byte("DATA\x00NULL"), 0o644)

	env, err := ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}

	tests := map[string]EnvValue{
		"VAR1":      {Value: "VALUE1", NeedRemove: false},
		"VAR2":      {Value: "VALUE2", NeedRemove: false},
		"EMPTY_VAR": {Value: "", NeedRemove: true},
		"NULL_VAR":  {Value: "DATA\nNULL", NeedRemove: false},
	}

	for key, expected := range tests {
		value, ok := env[key]
		if !ok {
			t.Errorf("Expected environment variable %s to be present", key)
			continue
		}
		if value.Value != expected.Value || value.NeedRemove != expected.NeedRemove {
			t.Errorf("For key %s, expected %v, got %v", key, expected, value)
		}
	}
}
