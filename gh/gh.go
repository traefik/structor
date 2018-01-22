package gh

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/containous/structor/types"
	"github.com/pkg/errors"
)

// GetLatestReleaseTagName find the latest release tag name
func GetLatestReleaseTagName(repoID types.RepoID) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	baseURL := fmt.Sprintf("https://github.com/%s/%s/releases", repoID.Owner, repoID.RepositoryName)
	resp, err := client.Get(baseURL + "/latest")
	if err != nil {
		return "", errors.Wrap(err, "error when trying to get latest release tag name")
	}

	location, err := resp.Location()
	if err != nil {
		log.Println("Status:", resp.Status)

		body, errBody := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(errBody)
			return "", err
		}

		errClose := resp.Body.Close()
		if errClose != nil {
			log.Println(errClose)
			return "", err
		}

		log.Println("Body:", string(body))

		return "", err
	}

	tag := strings.TrimPrefix(location.String(), baseURL+"/tag/")

	return tag, nil
}
