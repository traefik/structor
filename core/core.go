package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/containous/structor/copy"
	"github.com/containous/structor/gh"
	"github.com/containous/structor/menu"
	"github.com/containous/structor/repository"
	"github.com/containous/structor/types"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/ldez/go-git-cmd-wrapper/worktree"
	"github.com/pkg/errors"
)

const (
	baseRemote              = "origin/"
	envVarStructorLatestTag = "STRUCTOR_LATEST_TAG"
)

type dockerfileInformation struct {
	name      string
	path      string
	content   []byte
	imageName string
}

// Execute core process
func Execute(config *types.Configuration) error {
	workDir, err := ioutil.TempDir("", "structor")
	if err != nil {
		return err
	}

	defer func() {
		err = cleanAll(workDir, config.Debug)
		if err != nil {
			log.Println("Error during cleanning: ", err)
		}
	}()

	if config.Debug {
		log.Printf("Temp directory: %s", workDir)
	}

	menuContent := getMenuTemplateContent(config.Menu)

	fallbackDockerFileContent, err := downloadFile(config.DockerfileURL)
	if err != nil {
		return errors.Wrap(err, "failed to download Dockerfile")
	}

	fallbackDockerfile := dockerfileInformation{
		name:      fmt.Sprintf("%v.Dockerfile", time.Now().UnixNano()),
		content:   fallbackDockerFileContent,
		imageName: config.DockerImageName,
	}

	repoID := types.RepoID{
		Owner:          config.Owner,
		RepositoryName: config.RepositoryName,
	}

	requirementsContent, err := getRequirementsContent(config.RequirementsURL)
	if err != nil {
		return err
	}

	return process(workDir, repoID, fallbackDockerfile, menuContent, requirementsContent, config)
}

func process(workDir string, repoID types.RepoID, fallbackDockerfile dockerfileInformation,
	menuContent types.MenuContent, requirementsContent []byte, config *types.Configuration) error {
	latestTagName, err := getLatestReleaseTagName(repoID)
	if err != nil {
		return err
	}

	log.Printf("Latest tag: %s", latestTagName)

	branches, err := getBranches(config.ExperimentalBranchName, config.Debug)
	if err != nil {
		return err
	}

	siteDir, err := buildSiteDirectory()
	if err != nil {
		return err
	}

	for _, branchRef := range branches {
		versionName := strings.Replace(branchRef, baseRemote, "", 1)
		versionCurrentPath := filepath.Join(workDir, versionName)

		log.Printf("Generating doc for version %s", versionName)

		err := repository.CreateWorkTree(versionCurrentPath, branchRef, config.Debug)
		if err != nil {
			return err
		}

		versionDocsRoot, err := getDocumentationRoot(versionCurrentPath)
		if err != nil {
			return err
		}

		versionsInfo := types.VersionsInformation{
			Current:      versionName,
			Latest:       latestTagName,
			Experimental: config.ExperimentalBranchName,
			CurrentPath:  versionDocsRoot,
		}
		fallbackDockerfile.path = filepath.Join(versionsInfo.CurrentPath, fallbackDockerfile.name)

		err = buildDocumentation(branches, versionsInfo, fallbackDockerfile, menuContent, requirementsContent, config)
		if err != nil {
			return err
		}

		outputDir := siteDir
		if strings.HasPrefix(latestTagName, versionName) {
			err = copy.Copy(filepath.Join(versionsInfo.CurrentPath, "site"), outputDir)
			if err != nil {
				return err
			}
		}

		outputDir = filepath.Join(outputDir, versionName)
		err = copy.Copy(filepath.Join(versionsInfo.CurrentPath, "site"), outputDir)
		if err != nil {
			return err
		}
	}

	return nil
}

// getDocumentationRoot returns the path to the documentation's root by searching for "${menu.ManifestFileName}".
// Search is done from the docsRootSearchPath, relatively to the provided repository path.
// An additional sanity checking is done on the file named "requirements.txt" which must be located in the same directory.
func getDocumentationRoot(repositoryRoot string) (string, error) {
	if repositoryRoot == "" {
		return "", errors.New("repositoryRoot is undefined")
	}
	if _, err := os.Stat(repositoryRoot); os.IsNotExist(err) {
		return "", err
	}

	var docsRootSearchPaths = [...]string{"/", "docs/"}

	for _, docsRootSearchPath := range docsRootSearchPaths {
		candidateDocsRootPath := filepath.Join(repositoryRoot, docsRootSearchPath)

		if _, err := os.Stat(filepath.Join(candidateDocsRootPath, menu.ManifestFileName)); !os.IsNotExist(err) {
			log.Printf("Found %s for building documentation in %s.", menu.ManifestFileName, candidateDocsRootPath)

			// Sanity checking: the file "requirements.txt" must be in the candidate directory, next to manifest file
			if _, err := os.Stat(filepath.Join(candidateDocsRootPath, "requirements.txt")); os.IsNotExist(err) {
				return "", err
			}

			return candidateDocsRootPath, nil
		}
	}

	return "", errors.New("no file " + menu.ManifestFileName + " found in " + repositoryRoot + " (search path was: " + strings.Join(docsRootSearchPaths[:], ",") + ")")
}

func buildDocumentation(branches []string, versionsInfo types.VersionsInformation,
	fallbackDockerfile dockerfileInformation, menuTemplateContent types.MenuContent, requirementsContent []byte,
	config *types.Configuration) error {

	err := menu.Build(versionsInfo, branches, menuTemplateContent)
	if err != nil {
		return err
	}

	err = buildRequirements(versionsInfo, requirementsContent)
	if err != nil {
		return err
	}

	baseDockerfile, err := getDockerfile(fallbackDockerfile, versionsInfo.CurrentPath, config.DockerfileName)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(baseDockerfile.path, baseDockerfile.content, os.ModePerm)
	if err != nil {
		return err
	}

	dockerImageFullName := getDockerImageFullName(baseDockerfile.imageName, versionsInfo.Current)

	// Build image
	output, err := dockerCmd(config.Debug, "build", "--no-cache="+strconv.FormatBool(config.NoCache), "-t", dockerImageFullName, "-f", baseDockerfile.path, versionsInfo.CurrentPath+"/")
	if err != nil {
		log.Println(output)
		return err
	}

	// Run image
	output, err = dockerCmd(config.Debug, "run", "--rm", "-v", versionsInfo.CurrentPath+":/mkdocs", dockerImageFullName, "mkdocs", "build")
	if err != nil {
		log.Println(output)
		return err
	}

	return nil
}

// getDockerImageFullName returns the full docker image name, in the form image:tag. Please note that normalization is applied to avoid fordbidden characters
func getDockerImageFullName(imageName string, tagName string) string {
	r := strings.NewReplacer(":", "-", "/", "-")
	return r.Replace(imageName) + ":" + r.Replace(tagName)
}

func getLatestReleaseTagName(repoID types.RepoID) (string, error) {
	latest := os.Getenv(envVarStructorLatestTag)
	if len(latest) > 0 {
		return latest, nil
	}

	return gh.GetLatestReleaseTagName(repoID)
}

func getBranches(experimentalBranchName string, debug bool) ([]string, error) {
	var branches []string

	if len(experimentalBranchName) > 0 {
		branches = append(branches, baseRemote+experimentalBranchName)
	}

	gitBranches, err := repository.ListBranches(debug)
	if err != nil {
		return nil, err
	}
	branches = append(branches, gitBranches...)

	if len(branches) == 0 {
		log.Println("[WARN] no branch.")
	}

	return branches, nil
}

func buildSiteDirectory() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	siteDir := filepath.Join(currentDir, "site")
	err = createDirectory(siteDir)
	if err != nil {
		return "", err
	}

	return siteDir, nil
}

func getMenuTemplateContent(menu *types.MenuFiles) types.MenuContent {
	var content types.MenuContent

	if menu.HasJsFile() {
		jsContent, err := getMenuFileContent(menu.JsFile, menu.JsURL)
		if err != nil {
			return types.MenuContent{}
		}
		content.Js = jsContent
	}

	if menu.HasCSSFile() {
		cssContent, err := getMenuFileContent(menu.CSSFile, menu.CSSURL)
		if err != nil {
			return types.MenuContent{}
		}
		content.CSS = cssContent
	}

	return content
}

func getMenuFileContent(f string, u string) ([]byte, error) {
	if len(f) > 0 {
		content, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get template menu file content")
		}
		return content, nil
	}

	content, err := downloadFile(u)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download menu template")
	}
	return content, nil
}

func getRequirementsContent(requirementsURL string) ([]byte, error) {
	var content []byte

	if len(requirementsURL) > 0 {
		_, err := os.Stat(requirementsURL)
		if err != nil {
			content, err = downloadFile(requirementsURL)
			if err != nil {
				return nil, errors.Wrap(err, "failed to download Requirements file")
			}
		} else {
			content, err = ioutil.ReadFile(requirementsURL)
			if err != nil {
				return nil, errors.Wrap(err, "failed to read Requirements file")
			}
		}
	}
	return content, nil
}

func buildRequirements(versionsInfo types.VersionsInformation, customContent []byte) error {
	if len(customContent) > 0 {
		requirementsPath := filepath.Join(versionsInfo.CurrentPath, "requirements.txt")

		baseContent, err := ioutil.ReadFile(requirementsPath)
		if err != nil {
			return err
		}

		reqBase, err := parseRequirements(baseContent)
		if err != nil {
			return err
		}

		reqCustom, err := parseRequirements(customContent)
		if err != nil {
			return err
		}

		// merge
		for key, value := range reqCustom {
			reqBase[key] = value
		}

		file, err := os.Create(requirementsPath)
		if err != nil {
			return err
		}
		defer safeClose(file.Close)

		for key, value := range reqBase {
			fmt.Fprintf(file, "%s%s\n", key, value)
		}
	}
	return nil
}

func parseRequirements(content []byte) (map[string]string, error) {
	exp := regexp.MustCompile(`([\w-_]+)([=|>|<].+)`)

	result := make(map[string]string)

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) > 0 {
			submatch := exp.FindStringSubmatch(line)
			if len(submatch) != 3 {
				return nil, errors.Errorf("invalid line format: %s", line)
			}

			result[submatch[1]] = submatch[2]
		}
	}

	return result, nil
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, resp.Body.Close()
}

func cleanAll(workDir string, debug bool) error {
	err := os.RemoveAll(workDir)
	if err != nil {
		return err
	}

	output, err := git.Worktree(worktree.Prune, git.Debugger(debug))
	if err != nil {
		log.Println(output)
		return err
	}

	return nil
}

func getDockerfile(fallbackDockerfile dockerfileInformation, workingDirectory string, dockerfileName string) (*dockerfileInformation, error) {
	if workingDirectory == "" {
		return nil, errors.New("workingDirectory is undefined")
	}
	if _, err := os.Stat(workingDirectory); os.IsNotExist(err) {
		return nil, err
	}

	searchPaths := []string{
		filepath.Join(workingDirectory, dockerfileName),
		filepath.Join(workingDirectory, "docs", dockerfileName),
	}

	for _, searchPath := range searchPaths {
		if _, err := os.Stat(searchPath); !os.IsNotExist(err) {
			log.Printf("Found Dockerfile for building documentation in %s.", searchPath)

			var dockerFileContent []byte
			dockerFileContent, err = ioutil.ReadFile(searchPath)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get dockerfile file content.")
			}
			baseDockerfile := dockerfileInformation{
				name:      dockerfileName,
				path:      searchPath,
				imageName: fallbackDockerfile.imageName,
				content:   dockerFileContent,
			}
			return &baseDockerfile, nil
		}
	}

	log.Printf("Using fallback Dockerfile, written into %s", fallbackDockerfile.path)
	return &fallbackDockerfile, nil
}

func dockerCmd(debug bool, args ...string) (string, error) {
	name := "docker"

	if debug {
		log.Println(name, strings.Join(args, " "))
	}

	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

func createDirectory(directoryPath string) error {
	_, err := os.Stat(directoryPath)
	if !os.IsNotExist(err) {
		err = os.RemoveAll(directoryPath)
		if err != nil {
			return err
		}
	}

	return os.MkdirAll(directoryPath, os.ModePerm)
}

func safeClose(fn func() error) {
	err := fn()
	if err != nil {
		log.Println(err)
	}
}
