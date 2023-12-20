package testsnapshot

import (
	"encoding/json"
	"runtime"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

type Snapshot struct {
	Path                string
	WindowsReplacements map[string]string
}

// applyWindowsReplacements will replace any matching strings if on Windows
func (s Snapshot) applyWindowsReplacements(content string) string {
	if //goland:noinspection GoBoolExpressions
	runtime.GOOS == "windows" {
		for replacement, match := range s.WindowsReplacements {
			content = strings.ReplaceAll(content, match, replacement)
		}
	}

	return content
}

func New(path string, windowsReplacements map[string]string) Snapshot {
	return Snapshot{Path: path, WindowsReplacements: windowsReplacements}
}

// MatchJSON asserts the existing snapshot matches what was gotten in the test,
// after being marshalled as JSON
func MatchJSON[V any](t *testing.T, snapshot Snapshot, got V) {
	t.Helper()

	j, err := json.Marshal(got)

	if err != nil {
		t.Fatalf("Failed to marshal JSON: %s", err)
	}

	snaps.MatchSnapshot(t, snapshot.applyWindowsReplacements(string(j)))
}

// MatchText asserts the existing snapshot matches what was gotten in the test
func MatchText(t *testing.T, snapshot Snapshot, got string) {
	t.Helper()

	snaps.MatchSnapshot(t, snapshot.applyWindowsReplacements(got))
}
