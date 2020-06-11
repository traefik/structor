package types

import "strings"

// NoOption empty struct.
type NoOption struct{}

// Configuration task configuration.
type Configuration struct {
	Owner                  string     `short:"o" description:"Repository owner. [required]"`
	RepositoryName         string     `short:"r" long:"repo-name" description:"Repository name. [required]"`
	Debug                  bool       `long:"debug" description:"Debug mode."`
	DockerfileURL          string     `short:"d" long:"dockerfile-url" description:"Use this Dockerfile when --dockerfile-name is not found. Can be a file path. [required]"`
	DockerfileName         string     `long:"dockerfile-name" description:"Search and use this Dockerfile in the repository (in './docs/' or in './') for building documentation."`
	ExperimentalBranchName string     `long:"exp-branch" description:"Build a branch as experimental."`
	ExcludedBranches       []string   `long:"exclude" description:"Exclude branches from the documentation generation."`
	DockerImageName        string     `long:"image-name" description:"Docker image name."`
	DockerBuildPath        string     `long:"docker-build-path" description:"Docker build path"`
	BranchPrefix           string     `long:"branch-prefix" description:"Branch prefix"`
	Menu                   *MenuFiles `long:"menu" description:"Menu templates files."`
	RequirementsURL        string     `long:"rqts-url" description:"Use this requirements.txt to merge with the current requirements.txt. Can be a file path."`
	RequirementsFullPaths  []string   `long:"rqts-full-paths" description:"Requirements full path."`
	NoCache                bool       `long:"no-cache" description:"Set to 'true' to disable the Docker build cache."`
	ForceEditionURI        bool       `long:"force-edit-url" description:"Add a dedicated edition URL for each version."`
}

// MenuFiles menu template files references.
type MenuFiles struct {
	JsURL   string `long:"js-url" description:"URL of the template of the JS file use for the multi version menu."`
	JsFile  string `long:"js-file" description:"File path of the template of the JS file use for the multi version menu."`
	CSSURL  string `long:"css-url" description:"URL of the template of the CSS file use for the multi version menu."`
	CSSFile string `long:"css-file" description:"File path of the template of the CSS file use for the multi version menu."`
}

// HasJsFile has JS file.
func (m *MenuFiles) HasJsFile() bool {
	return m != nil && len(m.JsFile) > 0 || len(m.JsURL) > 0
}

// HasCSSFile has CSS file.
func (m *MenuFiles) HasCSSFile() bool {
	return m != nil && len(m.CSSFile) > 0 || len(m.CSSURL) > 0
}

// VersionsInformation versions information.
type VersionsInformation struct {
	Current      string
	Latest       string
	Experimental string
	CurrentPath  string
	BranchPrefix string
}

// GetCurrent get current version and replace prefix by v.
func (v VersionsInformation) GetCurrent() string {
	return strings.Replace(v.Current, v.BranchPrefix, "v", 1)
}

// GetLatest get latest version and replace prefix by v.
func (v VersionsInformation) GetLatest() string {
	return strings.Replace(v.Latest, v.BranchPrefix, "v", 1)
}

// GetExperimental get experimental version and replace prefix by v.
func (v VersionsInformation) GetExperimental() string {
	return strings.Replace(v.Experimental, v.BranchPrefix, "v", 1)
}
