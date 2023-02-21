package gh

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// GetLatestReleaseTagName find the latest release tag name.
func GetLatestReleaseTagName(owner, repositoryName string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	baseURL := fmt.Sprintf("https://github.com/%s/%s/releases", owner, repositoryName)

	resp, err := client.Get(baseURL + "/latest")
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		logResponseBody(resp)
		return "", fmt.Errorf("failed to get latest release tag name: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		logResponseBody(resp)
		return "", fmt.Errorf("failed to get latest release tag name on GitHub (%q), status: %s", baseURL+"/latest", resp.Status)
	}

	location, err := resp.Location()
	if err != nil {
		return "", fmt.Errorf("failed to get location header: %w", err)
	}

	return strings.TrimPrefix(location.String(), baseURL+"/tag/"), nil
}

func logResponseBody(resp *http.Response) {
	if resp.Body == nil {
		log.Println("The response body is empty")
		return
	}

	body, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		log.Println(errBody)
		return
	}

	log.Println("Body:", string(body))
}
