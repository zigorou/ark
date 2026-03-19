// Package vault provides utilities for managing the ark vault directory.
package vault

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrVaultNotInitialized is returned when the vault directory is not initialized.
var ErrVaultNotInitialized = errors.New("vault not initialized")

// CheckInitialized reports whether the vault at dir is fully initialized.
// A vault is considered initialized when dir exists, contains a .git/ directory,
// and contains a .sops.yaml file. Any missing condition returns an error wrapping
// ErrVaultNotInitialized.
func CheckInitialized(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("%w: directory %s does not exist", ErrVaultNotInitialized, dir)
	}

	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("%w: %s is not a git repository", ErrVaultNotInitialized, dir)
	}

	sopsYAML := filepath.Join(dir, ".sops.yaml")
	if _, err := os.Stat(sopsYAML); os.IsNotExist(err) {
		return fmt.Errorf("%w: .sops.yaml not found in %s", ErrVaultNotInitialized, dir)
	}

	return nil
}

// DefaultDir returns the default vault directory path (~/.ark).
func DefaultDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("vault: resolve home directory: %w", err)
	}
	return filepath.Join(home, ".ark"), nil
}

// ResolveDir resolves the vault directory using the following priority:
// flagValue (if non-empty) > ARK_VAULT_DIR env var > ~/.ark (default).
func ResolveDir(flagValue string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}
	if env := os.Getenv("ARK_VAULT_DIR"); env != "" {
		return env, nil
	}
	return DefaultDir()
}
