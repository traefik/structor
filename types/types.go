package types

// Configuration task configuration.
type Configuration struct {
	Owner                  string     `short:"o" description:"Repository owner. [required]"`
	RepositoryName         string     `short:"r" long:"repo-name" description:"Repository name. [required]"`
	Debug                  bool       `long:"debug" description:"Debug mode."`
	DockerfileURL          string     `short:"d" long:"dockerfile-url" description:"Dockerfile URL. [required]"`
	ExperimentalBranchName string     `long:"exp-branch" description:"Build a branch as experimental."`
	DockerImageName        string     `long:"image-name" description:"Docker image name."`
	Menu                   *MenuFiles `long:"menu" description:"Menu templates files."`
}

// MenuFiles menu template files references
type MenuFiles struct {
	JsURL   string `long:"js-url" description:"URL of the template of the JS file use for the multi version menu."`
	JsFile  string `long:"js-file" description:"File path of the template of the JS file use for the multi version menu."`
	CSSURL  string `long:"css-url" description:"URL of the template of the CSS file use for the multi version menu."`
	CSSFile string `long:"css-file" description:"File path of the template of the CSS file use for the multi version menu."`
}

// HasJsFile has JS file
func (m *MenuFiles) HasJsFile() bool {
	return m != nil && len(m.JsFile) > 0 || len(m.JsURL) > 0
}

// HasCSSFile has CSS file
func (m *MenuFiles) HasCSSFile() bool {
	return m != nil && len(m.CSSFile) > 0 || len(m.CSSURL) > 0
}

// MenuContent Content of menu files
type MenuContent struct {
	Js  []byte
	CSS []byte
}

// VersionsInformation versions information
type VersionsInformation struct {
	Current      string
	Latest       string
	Experimental string
	CurrentPath  string
}

// RepoID Repository identifier.
type RepoID struct {
	Owner          string
	RepositoryName string
}

// NoOption empty struct.
type NoOption struct{}
