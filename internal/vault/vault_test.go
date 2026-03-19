package vault_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/zigorou/ark/internal/vault"
)

func TestCheckInitialized(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(t *testing.T) string // returns dir path
		wantErr bool
	}{
		{
			name: "fully initialized",
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				must(t, os.Mkdir(filepath.Join(dir, ".git"), 0o755))
				must(t, os.WriteFile(filepath.Join(dir, ".sops.yaml"), []byte(""), 0o644))
				return dir
			},
			wantErr: false,
		},
		{
			name: "dir does not exist",
			setup: func(t *testing.T) string {
				t.Helper()
				return filepath.Join(t.TempDir(), "nonexistent")
			},
			wantErr: true,
		},
		{
			name: "no .git directory",
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				must(t, os.WriteFile(filepath.Join(dir, ".sops.yaml"), []byte(""), 0o644))
				return dir
			},
			wantErr: true,
		},
		{
			name: "no .sops.yaml",
			setup: func(t *testing.T) string {
				t.Helper()
				dir := t.TempDir()
				must(t, os.Mkdir(filepath.Join(dir, ".git"), 0o755))
				return dir
			},
			wantErr: true,
		},
		{
			name: "empty dir",
			setup: func(t *testing.T) string {
				t.Helper()
				return t.TempDir()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := tc.setup(t)
			err := vault.CheckInitialized(dir)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, vault.ErrVaultNotInitialized) {
					t.Errorf("expected ErrVaultNotInitialized, got %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestResolveDir(t *testing.T) {
	defaultDir, err := vault.DefaultDir()
	if err != nil {
		t.Fatalf("DefaultDir: %v", err)
	}

	tests := []struct {
		name      string
		flagValue string
		envValue  string
		want      string
	}{
		{
			name:      "flag set",
			flagValue: "/tmp/x",
			want:      "/tmp/x",
		},
		{
			name:     "env set",
			envValue: "/tmp/y",
			want:     "/tmp/y",
		},
		{
			name:      "flag wins over env",
			flagValue: "/tmp/x",
			envValue:  "/tmp/y",
			want:      "/tmp/x",
		},
		{
			name: "neither — use default",
			want: defaultDir,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envValue != "" {
				t.Setenv("ARK_VAULT_DIR", tc.envValue)
			} else {
				t.Setenv("ARK_VAULT_DIR", "")
			}

			got, err := vault.ResolveDir(tc.flagValue)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("ResolveDir(%q) = %q, want %q", tc.flagValue, got, tc.want)
			}
		})
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
