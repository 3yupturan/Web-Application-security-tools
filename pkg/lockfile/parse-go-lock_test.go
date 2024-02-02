package lockfile_test

import (
	"io/fs"
	"testing"

	"github.com/google/osv-scanner/pkg/models"

	"github.com/google/osv-scanner/pkg/lockfile"
)

func TestGoLockExtractor_ShouldExtract(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "",
			path: "",
			want: false,
		},
		{
			name: "",
			path: "go.mod",
			want: true,
		},
		{
			name: "",
			path: "path/to/my/go.mod",
			want: true,
		},
		{
			name: "",
			path: "path/to/my/go.mod/file",
			want: false,
		},
		{
			name: "",
			path: "path/to/my/go.mod.file",
			want: false,
		},
		{
			name: "",
			path: "path.to.my.go.mod",
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := lockfile.GoLockExtractor{}
			got := e.ShouldExtract(tt.path)
			if got != tt.want {
				t.Errorf("Extract() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseGoLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGoLock("fixtures/go/does-not-exist")

	expectErrIs(t, err, fs.ErrNotExist)
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseGoLock_Invalid(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGoLock("fixtures/go/not-go-mod.txt")

	expectErrContaining(t, err, "could not extract from")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseGoLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGoLock("fixtures/go/empty.mod")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseGoLock_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGoLock("fixtures/go/one-package.mod")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:         "github.com/BurntSushi/toml",
			Version:      "1.0.0",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 4, End: 4},
		},
	})
}

func TestParseGoLock_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGoLock("fixtures/go/two-packages.mod")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:         "github.com/BurntSushi/toml",
			Version:      "1.0.0",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 6, End: 6},
		},
		{
			Name:         "gopkg.in/yaml.v2",
			Version:      "2.4.0",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 7, End: 7},
		},
		{
			Name:      "stdlib",
			Version:   "1.17.0",
			Ecosystem: lockfile.GoEcosystem,
			CompareAs: lockfile.GoEcosystem,
		},
	})
}

func TestParseGoLock_IndirectPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGoLock("fixtures/go/indirect-packages.mod")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:         "github.com/BurntSushi/toml",
			Version:      "1.0.0",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 6, End: 6},
		},
		{
			Name:         "gopkg.in/yaml.v2",
			Version:      "2.4.0",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 7, End: 7},
		},
		{
			Name:         "github.com/mattn/go-colorable",
			Version:      "0.1.9",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 11, End: 11},
		},
		{
			Name:         "github.com/mattn/go-isatty",
			Version:      "0.0.14",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 12, End: 12},
		},
		{
			Name:         "golang.org/x/sys",
			Version:      "0.0.0-20210630005230-0f9fa26af87c",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 13, End: 13},
		},
		{
			Name:      "stdlib",
			Version:   "1.17.0",
			Ecosystem: lockfile.GoEcosystem,
			CompareAs: lockfile.GoEcosystem,
		},
	})
}

func TestParseGoLock_Replacements_One(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGoLock("fixtures/go/replace-one.mod")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:         "example.com/fork/net",
			Version:      "1.4.5",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 5, End: 5},
		},
	})
}

func TestParseGoLock_Replacements_Mixed(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGoLock("fixtures/go/replace-mixed.mod")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:         "example.com/fork/net",
			Version:      "1.4.5",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 7, End: 7},
		},
		{
			Name:         "golang.org/x/net",
			Version:      "0.5.6",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 3, End: 3},
		},
	})
}

func TestParseGoLock_Replacements_Local(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGoLock("fixtures/go/replace-local.mod")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:         "./fork/net",
			Version:      "",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 7, End: 7},
		},
		{
			Name:         "github.com/BurntSushi/toml",
			Version:      "1.0.0",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 3, End: 3},
		},
	})
}

func TestParseGoLock_Replacements_Different(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGoLock("fixtures/go/replace-different.mod")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:         "example.com/fork/foe",
			Version:      "1.4.5",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 7, End: 7},
		},
		{
			Name:         "example.com/fork/foe",
			Version:      "1.4.2",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 8, End: 8},
		},
	})
}

func TestParseGoLock_Replacements_NotRequired(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGoLock("fixtures/go/replace-not-required.mod")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:         "golang.org/x/net",
			Version:      "0.5.6",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 2, End: 2},
		},
		{
			Name:         "github.com/BurntSushi/toml",
			Version:      "1.0.0",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 3, End: 3},
		},
	})
}

func TestParseGoLock_Replacements_NoVersion(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGoLock("fixtures/go/replace-no-version.mod")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:         "example.com/fork/net",
			Version:      "1.4.5",
			Ecosystem:    lockfile.GoEcosystem,
			CompareAs:    lockfile.GoEcosystem,
			LinePosition: models.FilePosition{Start: 7, End: 7},
		},
	})
}
