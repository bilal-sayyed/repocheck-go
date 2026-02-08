package report

// RepoScanResult holds the complete result of a repository scan
type RepoScanResult struct {
	Summary      RepoSummary      `json:"summary"`
	Readiness    OnboardingStatus `json:"onboarding"`
	Dependencies DependencyStats  `json:"dependencies"`
	Hygiene      HygieneChecks    `json:"hygiene"`

	// Paid Features
	LicenseRisk    *LicenseRisk    `json:"license_risk,omitempty"`
	ReadinessScore *ReadinessScore `json:"readiness_score,omitempty"`
	CIDevOps       *CIDevOps       `json:"ci_devops,omitempty"`
	Ownership      *Ownership      `json:"ownership,omitempty"`

	ExecutiveSummary string   `json:"executive_summary,omitempty"`
	ActionItems      []string `json:"action_items,omitempty"`
}

type RepoSummary struct {
	Language    string   `json:"language"`
	Type        string   `json:"type"` // API, CLI, Library, Frontend, etc.
	GitPresent  bool     `json:"git_present"`
	FileCount   int      `json:"file_count"`
	EntryPoints []string `json:"entry_points,omitempty"`
}

type OnboardingStatus struct {
	HasReadme            bool     `json:"has_readme"`
	HasRunInstructions   bool     `json:"has_run_instructions"`
	HasSetupInstructions bool     `json:"has_setup_instructions"`
	MissingElements      []string `json:"missing_elements"` // List of what's missing
}

type DependencyStats struct {
	Manager     string `json:"manager"` // go.mod, package.json
	DirectCount int    `json:"direct_count"`
	TotalCount  int    `json:"total_count"`
}

type HygieneChecks struct {
	HasEnvExample    bool     `json:"has_env_example"`
	HardcodedSecrets []string `json:"hardcoded_secrets,omitempty"` // List of files/lines
	ConfigPresent    bool     `json:"config_present"`
	HasGitignore     bool     `json:"has_gitignore"` // New V2
	EnvIgnored       bool     `json:"env_ignored"`   // New V2
}

type Ownership struct {
	HasCodeowners  bool `json:"has_codeowners"`
	HasMaintainers bool `json:"has_maintainers"` // package.json
}

type LicenseRisk struct {
	RiskLevel       string   `json:"risk_level"` // LOW, MEDIUM, HIGH
	RiskyLicenses   []string `json:"risky_licenses"`
	OldDependencies []string `json:"old_dependencies"`
	Reasons         []string `json:"reasons,omitempty"`
}

type ReadinessScore struct {
	Score     int             `json:"score"`
	Grade     string          `json:"grade"` // A, B, C, D, F
	Breakdown *ScoreBreakdown `json:"breakdown,omitempty"`
}

type ScoreBreakdown struct {
	Documentation int `json:"documentation"`
	Dependencies  int `json:"dependencies"`
	CIDevOps      int `json:"ci_devops"`
	Security      int `json:"security"`
}

type CIDevOps struct {
	HasDockerfile bool     `json:"has_dockerfile"`
	HasCI         bool     `json:"has_ci"` // GitHub Actions, GitLab, etc.
	HasTests      bool     `json:"has_tests"`
	Tools         []string `json:"tools"`
}
