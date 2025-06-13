package os_test

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"

	museos "github.com/RWAs-labs/muse/pkg/os"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestResolveHome(t *testing.T) {
	usr, err := user.Current()
	require.NoError(t, err)

	testCases := []struct {
		name     string
		pathIn   string
		expected string
		fail     bool
	}{
		{
			name:     `should resolve home with leading "~/"`,
			pathIn:   "~/tmp/file.json",
			expected: filepath.Clean(filepath.Join(usr.HomeDir, "tmp/file.json")),
		},
		{
			name:     "should resolve '~'",
			pathIn:   `~`,
			expected: filepath.Clean(filepath.Join(usr.HomeDir, "")),
		},
		{
			name:     "should not resolve '~someuser/tmp'",
			pathIn:   `~someuser/tmp`,
			expected: `~someuser/tmp`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			pathOut, err := museos.ExpandHomeDir(tc.pathIn)
			require.NoError(t, err)
			require.Equal(t, tc.expected, pathOut)
		})
	}
}

func TestFileExists(t *testing.T) {
	path := sample.CreateTempDir(t)

	// create a test file
	existingFile := filepath.Join(path, "test.txt")
	_, err := os.Create(existingFile)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		file     string
		expected bool
	}{
		{
			name:     "should return true for existing file",
			file:     existingFile,
			expected: true,
		},
		{
			name:     "should return false for non-existing file",
			file:     filepath.Join(path, "non-existing.txt"),
			expected: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			exists := museos.FileExists(tc.file)
			require.Equal(t, tc.expected, exists)
		})
	}
}
