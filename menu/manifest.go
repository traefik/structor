package menu

import (
	"github.com/traefik/structor/manifest"
)

func editManifest(manif map[string]interface{}, versionJsFile, versionCSSFile string) {
	// Append menu JS file
	manifest.AppendExtraJs(manif, versionJsFile)

	// Append menu CSS file
	manifest.AppendExtraCSS(manif, versionCSSFile)

	// reset site URL
	manif["site_url"] = ""
}
