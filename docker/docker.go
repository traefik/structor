package docker

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/containous/structor/types"
	"github.com/pkg/errors"
)

// DockerfileInformation Dockerfile information.
type DockerfileInformation struct {
	Name      string
	Path      string
	Content   []byte
	ImageName string
}

// BuildImage Builds a Docker image.
func BuildImage(config *types.Configuration, fallbackDockerfile DockerfileInformation, versionsInfo types.VersionsInformation) (string, error) {
	baseDockerfile, err := getDockerfile(fallbackDockerfile, versionsInfo.CurrentPath, config.DockerfileName)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(baseDockerfile.Path, baseDockerfile.Content, os.ModePerm)
	if err != nil {
		return "", err
	}

	dockerImageFullName := buildImageFullName(baseDockerfile.ImageName, versionsInfo.Current)

	// Build image
	output, err := Exec(config.Debug, "build", "--no-cache="+strconv.FormatBool(config.NoCache), "-t", dockerImageFullName, "-f", baseDockerfile.Path, versionsInfo.CurrentPath+"/")
	if err != nil {
		log.Println(output)
		return "", err
	}

	return dockerImageFullName, nil
}

// Exec Executes a docker command.
func Exec(debug bool, args ...string) (string, error) {
	cmdName := "docker"

	if debug {
		log.Println(cmdName, strings.Join(args, " "))
	}

	cmd := exec.Command(cmdName, args...)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

func getDockerfile(fallbackDockerfile DockerfileInformation, workingDirectory string, dockerfileName string) (*DockerfileInformation, error) {
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
			baseDockerfile := DockerfileInformation{
				Name:      dockerfileName,
				Path:      searchPath,
				ImageName: fallbackDockerfile.ImageName,
				Content:   dockerFileContent,
			}
			return &baseDockerfile, nil
		}
	}

	log.Printf("Using fallback Dockerfile, written into %s", fallbackDockerfile.Path)
	return &fallbackDockerfile, nil
}

// buildImageFullName returns the full docker image name, in the form image:tag.
// Please note that normalization is applied to avoid forbidden characters.
func buildImageFullName(imageName string, tagName string) string {
	r := strings.NewReplacer(":", "-", "/", "-")
	return r.Replace(imageName) + ":" + r.Replace(tagName)
}
