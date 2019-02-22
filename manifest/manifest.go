package manifest

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// FileName file name of the mkdocs manifest file.
const FileName = "mkdocs.yml"

// Read Reads the manifest.
func Read(manifestFilePath string) (map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile(manifestFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "error when reading MkDocs Manifest.")
	}

	manif := make(map[string]interface{})

	err = yaml.Unmarshal(bytes, manif)
	if err != nil {
		return nil, errors.Wrap(err, "error when during unmarshal of the MkDocs Manifest.")
	}

	return manif, nil
}

// Write Writes the manifest.
func Write(manifestFilePath string, manif map[string]interface{}) error {
	out, err := yaml.Marshal(manif)
	if err != nil {
		return errors.Wrap(err, "error when marshal MkDocs Manifest.")
	}

	return ioutil.WriteFile(manifestFilePath, out, os.ModePerm)
}

// GetDocsDir returns the path to the directory pointed by "docs_dir" in the manifest file.
// If docs_dir is not set, then the directory containing the manifest is returned.
func GetDocsDir(manifestFilePath string, manif map[string]interface{}) (string, error) {
	manifestDir := filepath.Dir(manifestFilePath)

	if value, ok := manif["docs_dir"]; ok {
		return filepath.Join(manifestDir, value.(string)), nil
	}
	// https://www.mkdocs.org/user-guide/configuration/#docs_dir
	return filepath.Join(manifestDir, "docs"), nil
}

// AppendExtraJs Appends a file path to the "extra_javascript" in the manifest file.
func AppendExtraJs(jsFile string, manif map[string]interface{}) {
	if len(jsFile) > 0 {
		var extraJs []interface{}
		if val, ok := manif["extra_javascript"]; ok {
			extraJs = val.([]interface{})
		}
		extraJs = append(extraJs, jsFile)
		manif["extra_javascript"] = extraJs
	}
}

// AppendExtraCSS Appends a file path to the "extra_css" in the manifest file.
func AppendExtraCSS(cssFile string, manif map[string]interface{}) {
	if len(cssFile) > 0 {
		var extraCSS []interface{}
		if val, ok := manif["extra_css"]; ok {
			extraCSS = val.([]interface{})
		}
		extraCSS = append(extraCSS, cssFile)
		manif["extra_css"] = extraCSS
	}
}
