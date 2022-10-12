package manifest

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// FileName file name of the mkdocs manifest file.
	FileName = "mkdocs.yml"
	// TempPrefixEnvName temp prefix for environment variable.
	TempPrefixEnvName = "STRUCTOR_TEMP_"
)

// Read Reads the manifest.
func Read(manifestFilePath string) (map[string]interface{}, error) {
	content, err := os.ReadFile(manifestFilePath)
	if err != nil {
		return nil, fmt.Errorf("error when reading MkDocs Manifest: %w", err)
	}

	content = replaceEnvVariables(content)
	content = regexp.MustCompile(`: (!!python.+)`).ReplaceAll(content, []byte(`: '$1'`))

	manif := make(map[string]interface{})

	if err = yaml.Unmarshal(content, manif); err != nil {
		return nil, fmt.Errorf("error when during unmarshal of the MkDocs Manifest: %w", err)
	}

	return manif, nil
}

func replaceEnvVariables(bytes []byte) []byte {
	data := string(bytes)
	re := regexp.MustCompile(`!!python\/object\/apply:os\.getenv\s\[[\'|\"]?([A-Z0-9_-]+)[\'|\"]?\]`)
	result := re.FindAllStringSubmatch(data, -1)
	for _, value := range result {
		data = strings.ReplaceAll(data, value[0], TempPrefixEnvName+value[1])
	}

	return []byte(data)
}

func rewriteEnvVariables(bytes []byte) []byte {
	data := string(bytes)
	re := regexp.MustCompile(TempPrefixEnvName + `([A-Z0-9_-]+)`)
	result := re.FindAllStringSubmatch(data, -1)

	for _, value := range result {
		data = strings.ReplaceAll(data, value[0], fmt.Sprintf(`!!python/object/apply:os.getenv [%q]`, value[1]))
	}

	return []byte(data)
}

// Write Writes the manifest.
func Write(manifestFilePath string, manif map[string]interface{}) error {
	content, err := yaml.Marshal(manif)
	if err != nil {
		return fmt.Errorf("error when marshal MkDocs Manifest: %w", err)
	}

	content = rewriteEnvVariables(content)
	content = regexp.MustCompile(`: '(!!python.+)'`).ReplaceAll(content, []byte(`: $1`))

	return os.WriteFile(manifestFilePath, content, os.ModePerm)
}

// GetDocsDir returns the path to the directory pointed by "docs_dir" in the manifest file.
// If docs_dir is not set, then the directory containing the manifest is returned.
func GetDocsDir(manif map[string]interface{}, manifestFilePath string) string {
	return filepath.Join(filepath.Dir(manifestFilePath), getDocsDirAttribute(manif))
}

func getDocsDirAttribute(manif map[string]interface{}) string {
	if value, ok := manif["docs_dir"]; ok {
		return value.(string)
	}

	// https://www.mkdocs.org/user-guide/configuration/#docs_dir
	return "docs"
}

// AppendExtraJs Appends a file path to the "extra_javascript" in the manifest file.
func AppendExtraJs(manif map[string]interface{}, jsFile string) {
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
func AppendExtraCSS(manif map[string]interface{}, cssFile string) {
	if len(cssFile) > 0 {
		var extraCSS []interface{}
		if val, ok := manif["extra_css"]; ok {
			extraCSS = val.([]interface{})
		}
		extraCSS = append(extraCSS, cssFile)
		manif["extra_css"] = extraCSS
	}
}

// AddEditionURI Adds an edition URI to the "edit_uri" in the manifest file.
func AddEditionURI(manif map[string]interface{}, version, docsDirBase string, override bool) {
	v := version
	if v == "" {
		v = "master"
	}

	if val, ok := manif["edit_uri"]; ok && len(val.(string)) > 0 && !override {
		// noop
		return
	}

	docsDir := getDocsDirAttribute(manif)

	manif["edit_uri"] = path.Join("edit", v, docsDirBase, docsDir) + "/"
}
