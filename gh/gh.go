package gh

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// GetLatestReleaseTagName find the latest release tag name
func GetLatestReleaseTagName(owner, repositoryName string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	baseURL := fmt.Sprintf("https://github.com/%s/%s/releases", owner, repositoryName)

	resp, err := client.Get(baseURL + "/latest")
	if err != nil {
		displayResponseBody(resp)
		return "", errors.Wrap(err, "failed to get latest release tag name")
	}

	log.Println("Status:", resp.Status)

	if resp.StatusCode >= http.StatusBadRequest {
		displayResponseBody(resp)
		return "", errors.Errorf("failed to get latest release tag name on GitHub (%q), status: %s", baseURL+"/latest", resp.Status)
	}

	location, err := resp.Location()
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(location.String(), baseURL+"/tag/"), nil
}

func displayResponseBody(resp *http.Response) {
	if resp.Body == nil {
		log.Println("The response body is empty")
		return
	}

	defer safeClose(resp.Body.Close)

	body, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		log.Println(errBody)
		return
	}

	log.Println("Body:", string(body))
}

func safeClose(fn func() error) {
	err := fn()
	if err != nil {
		log.Println(err)
	}
}
