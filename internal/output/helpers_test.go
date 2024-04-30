package output_test

import (
	"testing"

	"github.com/google/osv-scanner/pkg/models"
)

type outputTestCaseArgs struct {
	vulnResult *models.VulnerabilityResults
}

type outputTestCase struct {
	name string
	args outputTestCaseArgs
}

type outputTestRunner = func(t *testing.T, args outputTestCaseArgs)

func testOutputWithVulnerabilities(t *testing.T, run outputTestRunner) {
	t.Helper()

	tests := []outputTestCase{
		{
			name: "no sources",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					Results: []models.PackageSource{},
				},
			},
		},
		{
			name: "one source with no packages",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					Results: []models.PackageSource{
						{
							Source:   models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{},
						},
					},
				},
			},
		},
		{
			name: "multiple sources with no packages",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					Results: []models.PackageSource{
						{
							Source:   models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{},
						},
						{
							Source:   models.SourceInfo{Path: "path/to/my/second/lockfile"},
							Packages: []models.PackageVulns{},
						},
						{
							Source:   models.SourceInfo{Path: "path/to/my/third/lockfile"},
							Packages: []models.PackageVulns{},
						},
					},
				},
			},
		},
		{
			name: "one source with one package, no vulnerabilities",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Vulnerabilities: models.Vulnerabilities{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "multiple sources with a mixed count of packages, no vulnerabilities",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Vulnerabilities: models.Vulnerabilities{},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/second/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine2",
										Version:   "3.2.5",
										Ecosystem: "npm",
									},
									Vulnerabilities: models.Vulnerabilities{},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine3",
										Version:   "0.4.1",
										Ecosystem: "npm",
									},
									Vulnerabilities: models.Vulnerabilities{},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/third/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.3.5",
										Ecosystem: "npm",
									},
									Vulnerabilities: models.Vulnerabilities{},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Vulnerabilities: models.Vulnerabilities{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "one source with one package and one vulnerability",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{IDs: []string{"OSV-1"}}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-1",
											Summary:  "Something scary!",
											Severity: []models.Severity{{Type: "high", Score: "1"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "one source with one package and one vulnerability (dev)",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									DepGroups: []string{"dev"},
									Groups:    []models.GroupInfo{{IDs: []string{"OSV-1"}}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-1",
											Summary:  "Something scary!",
											Severity: []models.Severity{{Type: "high", Score: "1"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "one source with one package and two aliases of a single vulnerability",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{
										IDs:     []string{"OSV-1", "GHSA-123"},
										Aliases: []string{"OSV-1", "GHSA-123"},
									}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-1",
											Summary:  "Something scary!",
											Severity: []models.Severity{{Type: "high", Score: "1"}},
										},
										{
											ID:       "GHSA-123",
											Summary:  "Something scary!",
											Aliases:  []string{"OSV-1"},
											Severity: []models.Severity{{Type: "high", Score: "1"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "two sources with packages, one vulnerability",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{IDs: []string{"OSV-1"}}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-1",
											Summary:  "Something scary!",
											Severity: []models.Severity{{Type: "high", Score: "1"}},
										},
									},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/second/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine2",
										Version:   "5.9.0",
										Ecosystem: "npm",
									},
									Vulnerabilities: models.Vulnerabilities{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "multiple sources with a mixed count of packages, some vulnerabilities",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{IDs: []string{"OSV-1"}}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-1",
											Summary:  "Something scary!",
											Severity: []models.Severity{{Type: "high", Score: "1"}},
										},
									},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/second/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine2",
										Version:   "3.2.5",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{IDs: []string{"OSV-2"}}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-2",
											Summary:  "Something less scary!",
											Severity: []models.Severity{{Type: "low", Score: "1"}},
										},
									},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine3",
										Version:   "0.4.1",
										Ecosystem: "npm",
									},
									Vulnerabilities: models.Vulnerabilities{},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/third/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.3.5",
										Ecosystem: "npm",
									},
									Vulnerabilities: models.Vulnerabilities{},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{IDs: []string{"OSV-1"}}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-1",
											Summary:  "Something scary!",
											Severity: []models.Severity{{Type: "high", Score: "1"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "multiple sources with a mixed count of packages, and multiple vulnerabilities",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{
										{IDs: []string{"OSV-1"}},
										{IDs: []string{"OSV-5"}},
									},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-1",
											Summary:  "Something scary!",
											Severity: []models.Severity{{Type: "high", Score: "1"}},
										},
										{
											ID:       "OSV-5",
											Summary:  "Something scarier!",
											Severity: []models.Severity{{Type: "extreme", Score: "1"}},
										},
									},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.2",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{IDs: []string{"OSV-1"}}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-1",
											Summary:  "Something scary!",
											Severity: []models.Severity{{Type: "high", Score: "1"}},
										},
									},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/second/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine2",
										Version:   "3.2.5",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{IDs: []string{"OSV-2"}}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-2",
											Summary:  "Something less scary!",
											Severity: []models.Severity{{Type: "low", Score: "1"}},
										},
									},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine3",
										Version:   "0.4.1",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{
										{IDs: []string{"OSV-3"}},
										{IDs: []string{"OSV-5"}},
									},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-3",
											Summary:  "Something mildly scary!",
											Severity: []models.Severity{{Type: "medium", Score: "1"}},
										},
										{
											ID:       "OSV-5",
											Summary:  "Something scarier!",
											Severity: []models.Severity{{Type: "extreme", Score: "1"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "multiple sources with a mixed count of grouped packages, and multiple vulnerabilities",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									DepGroups: []string{"dev", "optional"},
									Groups: []models.GroupInfo{
										{IDs: []string{"OSV-1"}},
										{IDs: []string{"OSV-5"}},
									},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-1",
											Summary:  "Something scary!",
											Severity: []models.Severity{{Type: "high", Score: "1"}},
										},
										{
											ID:       "OSV-5",
											Summary:  "Something scarier!",
											Severity: []models.Severity{{Type: "extreme", Score: "1"}},
										},
									},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.2",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{IDs: []string{"OSV-1"}}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-1",
											Summary:  "Something scary!",
											Severity: []models.Severity{{Type: "high", Score: "1"}},
										},
									},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/second/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine2",
										Version:   "3.2.5",
										Ecosystem: "npm",
									},
									DepGroups: []string{"dev"},
									Groups:    []models.GroupInfo{{IDs: []string{"OSV-2"}}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-2",
											Summary:  "Something less scary!",
											Severity: []models.Severity{{Type: "low", Score: "1"}},
										},
									},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine3",
										Version:   "0.4.1",
										Ecosystem: "npm",
									},
									DepGroups: []string{"build"},
									Groups: []models.GroupInfo{
										{IDs: []string{"OSV-3"}},
										{IDs: []string{"OSV-5"}},
									},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-3",
											Summary:  "Something mildly scary!",
											Severity: []models.Severity{{Type: "medium", Score: "1"}},
										},
										{
											ID:       "OSV-5",
											Summary:  "Something scarier!",
											Severity: []models.Severity{{Type: "extreme", Score: "1"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "one source with vulnerabilities, some missing content",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{IDs: []string{"OSV-1"}}},
									Vulnerabilities: models.Vulnerabilities{
										{ID: "OSV-1", Details: "This vulnerability allows for some very scary stuff to happen - seriously, you'd not believe it!"},
									},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine3",
										Version:   "0.10.2-rc",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{IDs: []string{"OSV-2"}}},
									Vulnerabilities: models.Vulnerabilities{
										{ID: "OSV-2"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			run(t, tt.args)
		})
	}
}

func testOutputWithLicenseViolations(t *testing.T, run outputTestRunner) {
	t.Helper()

	experimentalAnalysisConfig := models.ExperimentalAnalysisConfig{
		Licenses: models.ExperimentalLicenseConfig{Summary: false, Allowlist: []models.License{"ISC"}},
	}

	tests := []outputTestCase{
		{
			name: "no sources",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					ExperimentalAnalysisConfig: experimentalAnalysisConfig,
					Results:                    []models.PackageSource{},
				},
			},
		},
		{
			name: "one source with no packages",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					ExperimentalAnalysisConfig: experimentalAnalysisConfig,
					Results: []models.PackageSource{
						{
							Source:   models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{},
						},
					},
				},
			},
		},
		{
			name: "multiple sources with no packages",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					ExperimentalAnalysisConfig: experimentalAnalysisConfig,
					Results: []models.PackageSource{
						{
							Source:   models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{},
						},
						{
							Source:   models.SourceInfo{Path: "path/to/my/second/lockfile"},
							Packages: []models.PackageVulns{},
						},
						{
							Source:   models.SourceInfo{Path: "path/to/my/third/lockfile"},
							Packages: []models.PackageVulns{},
						},
					},
				},
			},
		},
		{
			name: "one source with one package, no license violations",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					ExperimentalAnalysisConfig: experimentalAnalysisConfig,
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"ISC"},
									LicenseViolations: []models.License{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "multiple sources with a mixed count of packages, no license violations",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					ExperimentalAnalysisConfig: experimentalAnalysisConfig,
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"ISC"},
									LicenseViolations: []models.License{},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/second/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine2",
										Version:   "3.2.5",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"ISC"},
									LicenseViolations: []models.License{},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine3",
										Version:   "0.4.1",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"ISC"},
									LicenseViolations: []models.License{},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/third/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.3.5",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"ISC"},
									LicenseViolations: []models.License{},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"ISC"},
									LicenseViolations: []models.License{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "one source with one package and one license violation",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					ExperimentalAnalysisConfig: experimentalAnalysisConfig,
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"MIT"},
									LicenseViolations: []models.License{"MIT"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "one source with one package and one license violation (dev)",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					ExperimentalAnalysisConfig: experimentalAnalysisConfig,
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									DepGroups:         []string{"dev"},
									Licenses:          []models.License{"MIT"},
									LicenseViolations: []models.License{"MIT"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "two sources with packages, one license violation",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					ExperimentalAnalysisConfig: experimentalAnalysisConfig,
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"MIT"},
									LicenseViolations: []models.License{"MIT"},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/second/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine2",
										Version:   "5.9.0",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"ISC"},
									LicenseViolations: []models.License{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "multiple sources with a mixed count of packages, some license violations",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					ExperimentalAnalysisConfig: experimentalAnalysisConfig,
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"MIT"},
									LicenseViolations: []models.License{"MIT"},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/second/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine2",
										Version:   "3.2.5",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"Apache-2.0"},
									LicenseViolations: []models.License{"Apache-2.0"},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine3",
										Version:   "0.4.1",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"ISC"},
									LicenseViolations: []models.License{},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/third/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.3.5",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"ISC"},
									LicenseViolations: []models.License{},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"MIT"},
									LicenseViolations: []models.License{"MIT"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "multiple sources with a mixed count of packages and groups, some license violations",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					ExperimentalAnalysisConfig: experimentalAnalysisConfig,
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									DepGroups:         []string{"dev", "optional"},
									Licenses:          []models.License{"MIT"},
									LicenseViolations: []models.License{"MIT"},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/second/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine2",
										Version:   "3.2.5",
										Ecosystem: "npm",
									},
									DepGroups:         []string{"dev", "optional"},
									Licenses:          []models.License{"Apache-2.0"},
									LicenseViolations: []models.License{"Apache-2.0"},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine3",
										Version:   "0.4.1",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"ISC"},
									LicenseViolations: []models.License{},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/third/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.3.5",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"ISC"},
									LicenseViolations: []models.License{},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									DepGroups:         []string{"build"},
									Licenses:          []models.License{"MIT"},
									LicenseViolations: []models.License{"MIT"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "one source with one package and multiple license violations",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					ExperimentalAnalysisConfig: experimentalAnalysisConfig,
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Licenses:          []models.License{"MIT", "Apache-2.0"},
									LicenseViolations: []models.License{"MIT", "Apache-2.0"},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			run(t, tt.args)
		})
	}
}

func testOutputWithMixedIssues(t *testing.T, run outputTestRunner) {
	t.Helper()

	experimentalAnalysisConfig := models.ExperimentalAnalysisConfig{
		Licenses: models.ExperimentalLicenseConfig{Summary: false, Allowlist: []models.License{"ISC"}},
	}

	tests := []outputTestCase{
		{
			name: "one source with one package, one vulnerability, and one license violation",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					ExperimentalAnalysisConfig: experimentalAnalysisConfig,
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{IDs: []string{"OSV-1"}}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-1",
											Summary:  "Something scary!",
											Severity: []models.Severity{{Type: "high", Score: "1"}},
										},
									},
									Licenses:          []models.License{"MIT"},
									LicenseViolations: []models.License{"MIT"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "two sources with packages, one vulnerability, one license violation",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					ExperimentalAnalysisConfig: experimentalAnalysisConfig,
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{IDs: []string{"OSV-1"}}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-1",
											Summary:  "Something scary!",
											Severity: []models.Severity{{Type: "high", Score: "1"}},
										},
									},
									Licenses:          []models.License{"ISC"},
									LicenseViolations: []models.License{},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/second/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine2",
										Version:   "5.9.0",
										Ecosystem: "npm",
									},
									Groups:            []models.GroupInfo{},
									Vulnerabilities:   models.Vulnerabilities{},
									Licenses:          []models.License{"MIT"},
									LicenseViolations: []models.License{"MIT"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "multiple sources with a mixed count of packages, some vulnerabilities and license violations",
			args: outputTestCaseArgs{
				vulnResult: &models.VulnerabilityResults{
					ExperimentalAnalysisConfig: experimentalAnalysisConfig,
					Results: []models.PackageSource{
						{
							Source: models.SourceInfo{Path: "path/to/my/first/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{IDs: []string{"OSV-1"}}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-1",
											Summary:  "Something scary!",
											Severity: []models.Severity{{Type: "high", Score: "1"}},
										},
									},
									Licenses:          []models.License{"MIT"},
									LicenseViolations: []models.License{"MIT"},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/second/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine2",
										Version:   "3.2.5",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{IDs: []string{"OSV-2"}}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-2",
											Summary:  "Something less scary!",
											Severity: []models.Severity{{Type: "low", Score: "1"}},
										},
									},
									Licenses:          []models.License{"ISC"},
									LicenseViolations: []models.License{},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine3",
										Version:   "0.4.1",
										Ecosystem: "npm",
									},
									Vulnerabilities:   models.Vulnerabilities{},
									Licenses:          []models.License{"ISC"},
									LicenseViolations: []models.License{},
								},
							},
						},
						{
							Source: models.SourceInfo{Path: "path/to/my/third/lockfile"},
							Packages: []models.PackageVulns{
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.3.5",
										Ecosystem: "npm",
									},
									Vulnerabilities:   models.Vulnerabilities{},
									Licenses:          []models.License{"MIT"},
									LicenseViolations: []models.License{"MIT"},
								},
								{
									Package: models.PackageInfo{
										Name:      "mine1",
										Version:   "1.2.3",
										Ecosystem: "npm",
									},
									Groups: []models.GroupInfo{{IDs: []string{"OSV-1"}}},
									Vulnerabilities: models.Vulnerabilities{
										{
											ID:       "OSV-1",
											Summary:  "Something scary!",
											Severity: []models.Severity{{Type: "high", Score: "1"}},
										},
									},
									Licenses:          []models.License{"Apache-2.0"},
									LicenseViolations: []models.License{"Apache-2.0"},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			run(t, tt.args)
		})
	}
}
