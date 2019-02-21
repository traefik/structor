package requirements

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/containous/structor/file"
	"github.com/containous/structor/types"
	"github.com/pkg/errors"
)

// Check return an error if the requirements file is not found in the doc root directory.
func Check(docRoot string) error {
	_, err := os.Stat(filepath.Join(docRoot, "requirements.txt"))
	return err
}

func GetContent(requirementsURL string) ([]byte, error) {
	var content []byte

	if len(requirementsURL) > 0 {
		_, err := os.Stat(requirementsURL)
		if err != nil {
			content, err = file.Download(requirementsURL)
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

func Build(versionsInfo types.VersionsInformation, customContent []byte) error {
	if len(customContent) > 0 {
		requirementsPath := filepath.Join(versionsInfo.CurrentPath, "requirements.txt")

		baseContent, err := ioutil.ReadFile(requirementsPath)
		if err != nil {
			return err
		}

		reqBase, err := parse(baseContent)
		if err != nil {
			return err
		}

		reqCustom, err := parse(customContent)
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

func parse(content []byte) (map[string]string, error) {
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

func safeClose(fn func() error) {
	err := fn()
	if err != nil {
		log.Println(err)
	}
}
