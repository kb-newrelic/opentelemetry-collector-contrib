// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoInstall(t *testing.T) {
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "telemetrygen")
	t.Log("Using temporary GOBIN for test:", tmpDir)

	_, err := os.Stat(binPath)
	require.ErrorIs(t, err, os.ErrNotExist, "telemetrygen should not exist before install")

	shaCmd := exec.Command("git", "rev-parse", "HEAD")
	shaOut, err := shaCmd.Output()
	require.NoError(t, err, "failed to get git SHA")
	gitSHA := strings.TrimSpace(string(shaOut))
	require.Regexp(t, `^[0-9a-f]{40}$`, gitSHA, "git SHA should be a valid 40-character hex string")

	// #nosec G204 -- gitSHA is validated against regex to eliminate command injection risk
	installCmd := exec.Command("go", "install",
		"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/telemetrygen@"+gitSHA)
	installCmd.Env = append(os.Environ(), "GOBIN="+tmpDir)
	output, err := installCmd.CombinedOutput()
	require.NoError(t, err, "go install with sha %s should succeed: %s", gitSHA, string(output))

	helpCmd := exec.Command(binPath, "--help")
	helpOut, err := helpCmd.CombinedOutput()
	require.NoError(t, err, "telemetrygen --help should succeed: %s", string(helpOut))

	helpOutput := string(helpOut)
	require.Contains(t, helpOutput, "telemetrygen", "help output should mention telemetrygen")
	require.Contains(t, helpOutput, "traces", "help output should mention traces command")
}
