package docker

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/containous/structor/file"
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
func (d *DockerfileInformation) BuildImage(versionsInfo types.VersionsInformation, noCache, debug bool) (string, error) {
	err := ioutil.WriteFile(d.Path, d.Content, os.ModePerm)
	if err != nil {
		return "", err
	}

	dockerImageFullName := buildImageFullName(d.ImageName, versionsInfo.Current)

	// Build image
	output, err := Exec(debug, "build", "--no-cache="+strconv.FormatBool(noCache), "-t", dockerImageFullName, "-f", d.Path, versionsInfo.CurrentPath+"/")
	if err != nil {
		log.Println(output)
		return "", err
	}

	return dockerImageFullName, nil
}

// buildImageFullName returns the full docker image name, in the form image:tag.
// Please note that normalization is applied to avoid forbidden characters.
func buildImageFullName(imageName string, tagName string) string {
	r := strings.NewReplacer(":", "-", "/", "-")
	return r.Replace(imageName) + ":" + r.Replace(tagName)
}

// GetDockerfileFallback Downloads and creates the DockerfileInformation of the Dockerfile fallback.
func GetDockerfileFallback(dockerfileURL, imageName string) (DockerfileInformation, error) {
	dockerFileContent, err := file.Download(dockerfileURL)
	if err != nil {
		return DockerfileInformation{}, errors.Wrap(err, "failed to download Dockerfile")
	}

	return DockerfileInformation{
		Name:      fmt.Sprintf("%v.Dockerfile", time.Now().UnixNano()),
		Content:   dockerFileContent,
		ImageName: imageName,
	}, nil
}

// GetDockerfile Gets the effective Dockerfile.
func GetDockerfile(workingDirectory string, fallbackDockerfile DockerfileInformation, dockerfileName string) (*DockerfileInformation, error) {
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

			return &DockerfileInformation{
				Name:      dockerfileName,
				Path:      searchPath,
				ImageName: fallbackDockerfile.ImageName,
				Content:   dockerFileContent,
			}, nil
		}
	}

	log.Printf("Using fallback Dockerfile, written into %s", fallbackDockerfile.Path)
	return &fallbackDockerfile, nil
}

// Exec Executes a docker command.
func Exec(debug bool, args ...string) (string, error) {
	cmdName := "docker"

	if debug {
		log.Println(cmdName, strings.Join(args, " "))
	}

	output, err := exec.Command(cmdName, args...).CombinedOutput()

	return string(output), err
}
