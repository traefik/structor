package requirements

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/traefik/structor/file"
	"github.com/traefik/structor/types"
)

const filename = "requirements.txt"

// Check return an error if the requirements file is not found in the doc root directory.
func Check(docRoot string) error {
	_, err := os.Stat(filepath.Join(docRoot, filename))
	return err
}

// GetContent Gets the content of the "requirements.txt".
func GetContent(requirementsPath string) ([]byte, error) {
	if len(requirementsPath) == 0 {
		return nil, nil
	}

	if _, errStat := os.Stat(requirementsPath); errStat == nil {
		content, err := ioutil.ReadFile(requirementsPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read Requirements file: %w", err)
		}
		return content, nil
	}

	content, err := file.Download(requirementsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to download Requirements file: %w", err)
	}
	return content, nil
}

// Build Builds a "requirements.txt" file.
func Build(versionsInfo types.VersionsInformation, customContent []byte) error {
	if len(customContent) == 0 {
		return nil
	}

	requirementsPath := filepath.Join(versionsInfo.CurrentPath, filename)

	baseContent, err := ioutil.ReadFile(requirementsPath)
	if err != nil {
		return fmt.Errorf("unable to read %s: %w", requirementsPath, err)
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

	f, err := os.Create(requirementsPath)
	if err != nil {
		return err
	}
	defer safeClose(f.Close)

	var sortedKeys []string
	for key := range reqBase {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		fmt.Fprintf(f, "%s%s\n", key, reqBase[key])
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
				return nil, fmt.Errorf("invalid line format: %s", line)
			}

			result[submatch[1]] = submatch[2]
		}
	}

	return result, nil
}

func safeClose(fn func() error) {
	if err := fn(); err != nil {
		log.Println(err)
	}
}
