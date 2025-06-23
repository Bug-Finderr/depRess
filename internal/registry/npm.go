package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	BASE_URL = "https://registry.npmjs.org/"
	Timeout  = 10 * time.Second
)

type PackageInfo struct {
	Name     string                     `json:"name"`
	Versions map[string]PackageVersions `json:"versions"`
}

type PackageVersions struct {
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
}

type Client struct {
	httpClient *http.Client
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: Timeout,
		},
	}
}

func (c *Client) GetPkgInfo(pkg string) (*PackageInfo, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s%s", BASE_URL, pkg))

	if err != nil {
		return nil, fmt.Errorf("failed to get package info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("npm registry returned status %d for package %s", resp.StatusCode, pkg)
	}

	var pkgInfo PackageInfo
	if err := json.NewDecoder(resp.Body).Decode(&pkgInfo); err != nil {
		return nil, fmt.Errorf("failed to decode package info: %w", err)
	}

	return &pkgInfo, nil
}
