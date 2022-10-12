package menu

import (
	"fmt"
	"os"
	"path/filepath"
)

const menuCSSFileName = "structor-menu.css"

func writeCSSFile(manifestDocsDir string, menuContent Content) (string, error) {
	if len(menuContent.CSS) == 0 {
		return "", nil
	}

	cssDir := filepath.Join(manifestDocsDir, "theme", "css")
	if _, errStat := os.Stat(cssDir); os.IsNotExist(errStat) {
		errDir := os.MkdirAll(cssDir, os.ModePerm)
		if errDir != nil {
			return "", fmt.Errorf("error when create CSS folder: %w", errDir)
		}
	}

	err := os.WriteFile(filepath.Join(cssDir, menuCSSFileName), menuContent.CSS, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error when trying ro write CSS file: %w", err)
	}

	return filepath.Join("theme", "css", menuCSSFileName), nil
}
