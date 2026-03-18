package analyser

type RepoSignals struct {
	RepoName         string   `json:"repo_name"`
	Languages        []string `json:"languages"`
	Frameworks       []string `json:"frameworks"`

	HasDesignSystem  bool     `json:"has_design_system"`
	DesignSystemType string   `json:"design_system_type,omitempty"`
	HasStorybook     bool     `json:"has_storybook"`

	TestFramework    string   `json:"test_framework,omitempty"`
	HasE2E           bool     `json:"has_e2e"`
	E2EFramework     string   `json:"e2e_framework,omitempty"`

	CIProvider       string   `json:"ci_provider,omitempty"`
	HasDocker        bool     `json:"has_docker"`

	LintConfig       string   `json:"lint_config,omitempty"`
	IsMonorepo       bool     `json:"is_monorepo"`
	TopLevelDirs     []string `json:"top_level_dirs"`
	FileCount        int      `json:"file_count"`
	FileTree         []string `json:"file_tree,omitempty"`
}
