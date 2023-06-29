package offline

import (
	"fmt"
	"os"
	"path"

	"github.com/google/osv-scanner/pkg/lockfile"
	"github.com/google/osv-scanner/pkg/models"
	"github.com/google/osv-scanner/pkg/osv"
	"github.com/google/osv-scanner/pkg/reporter"
)

func loadDB(dbBasePath string, ecosystem lockfile.Ecosystem, offline bool) (*ZipDB, error) {
	return NewZippedDB(dbBasePath, string(ecosystem), fmt.Sprintf("https://osv-vulnerabilities.storage.googleapis.com/%s/all.zip", ecosystem), offline)
}

func toPackageDetails(query *osv.Query) (lockfile.PackageDetails, error) {
	if query.Package.PURL != "" {
		pkg, err := models.PURLToPackage(query.Package.PURL)

		if err != nil {
			return lockfile.PackageDetails{}, err
		}

		return lockfile.PackageDetails{
			Name:      pkg.Name,
			Version:   pkg.Version,
			Ecosystem: lockfile.Ecosystem(pkg.Ecosystem),
			CompareAs: lockfile.Ecosystem(pkg.Ecosystem),
		}, nil
	}

	return lockfile.PackageDetails{
		Name:      query.Package.Name,
		Version:   query.Version,
		Commit:    query.Commit,
		Ecosystem: query.Package.Ecosystem,
		CompareAs: query.Package.Ecosystem,
	}, nil
}

func setupLocalDBDirectory() (string, error) {
	var localDBPath string
	var err error
	var envSet bool

	if p, envSet := os.LookupEnv("OSV_SCANNER_LOCAL_DB_CACHE_DIRECTORY"); envSet {
		localDBPath = p
	} else {
		localDBPath, err = os.UserCacheDir()

		if err != nil {
			localDBPath = os.TempDir()
		}
	}

	err = os.Mkdir(path.Join(localDBPath, "osv-scanner"), 0750)

	// if the scanner cannot create its subdirectory when an explicit local db path
	// has been provided, then it should error rather than fallback to another path
	//
	// otherwise, it should fall back to the temp directory before erroring
	if err != nil && envSet {
		return path.Join(localDBPath, "osv-scanner"), err
	} else if localDBPath == os.TempDir() {
		localDBPath = os.TempDir()

		err = os.Mkdir(path.Join(localDBPath, "osv-scanner"), 0750)

		if err != nil {
			return path.Join(localDBPath, "osv-scanner"), err
		}
	}

	return path.Join(localDBPath, "osv-scanner"), nil
}

func Check(r reporter.Reporter, query osv.BatchedQuery, offline bool) (*osv.HydratedBatchedResponse, error) {
	results := make([]osv.Response, 0, len(query.Queries))
	dbs := make(map[lockfile.Ecosystem]*ZipDB)

	dbBasePath, err := setupLocalDBDirectory()

	if err != nil {
		return &osv.HydratedBatchedResponse{}, fmt.Errorf("could not create %s: %w", dbBasePath, err)
	}

	loadDBFromCache := func(ecosystem lockfile.Ecosystem) (*ZipDB, error) {
		if db, ok := dbs[ecosystem]; ok {
			return db, nil
		}

		db, err := loadDB(dbBasePath, ecosystem, offline)

		if err != nil {
			return nil, err
		}

		dbs[ecosystem] = db

		return db, nil
	}

	for _, query := range query.Queries {
		pkg, err := toPackageDetails(query)

		if err != nil {
			// currently, this will actually only error if the PURL cannot be parses
			r.PrintError(fmt.Sprintf("skipping %s as it is not a valid PURL: %v\n", query.Package.PURL, err))
			results = append(results, osv.Response{Vulns: []models.Vulnerability{}})

			continue
		}

		db, err := loadDBFromCache(pkg.Ecosystem)

		if err != nil {
			// currently, this will actually only error if the PURL cannot be parses
			r.PrintError(fmt.Sprintf("could not load db for %s ecosystem: %v\n", pkg.Ecosystem, err))
			results = append(results, osv.Response{Vulns: []models.Vulnerability{}})

			continue
		}

		results = append(results, osv.Response{Vulns: db.VulnerabilitiesAffectingPackage(pkg)})
	}

	return &osv.HydratedBatchedResponse{Results: results}, nil
}
