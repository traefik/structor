package manifest

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// FileName file name of the mkdocs manifest file
const FileName = "mkdocs.yml"

// GetManifestDocsDir returns the path to the directory pointed by "docs_dir" in the manifest file.
// If docs_dir is not set, then the directory containing the manifest is returned.
func GetManifestDocsDir(manifestFilePath string) (string, error) {
	bytes, err := ioutil.ReadFile(manifestFilePath)
	if err != nil {
		return "", errors.Wrap(err, "error when reading MkDocs Manifest.")
	}

	manif := make(map[string]interface{})

	err = yaml.Unmarshal(bytes, manif)
	if err != nil {
		return "", errors.Wrap(err, "error when during unmarshal of the MkDocs Manifest.")
	}

	manifestDir := filepath.Dir(manifestFilePath)

	if value, ok := manif["docs_dir"]; ok {

		return filepath.Join(manifestDir, value.(string)), nil
	}

	return manifestDir, nil
}
