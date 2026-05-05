package updatecheck

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const (
	FallbackVersion  = "1.1.0"
	latestVersionURL = "https://proxy.golang.org/github.com/igorrochap/commitgen/@latest"
)

type latestVersionResponse struct {
	Version string `json:"Version"`
	TagName string `json:"tag_name"`
}

type Result struct {
	Current string
	Latest  string
	Newer   bool
}

func Check(ctx context.Context) (Result, error) {
	current := CurrentVersion()
	latest, err := LatestVersion(ctx, http.DefaultClient, latestVersionURL)
	if err != nil {
		return Result{Current: current}, err
	}

	return Result{
		Current: current,
		Latest:  latest,
		Newer:   CompareVersions(latest, current) > 0,
	}, nil
}

func CheckWithTimeout(timeout time.Duration) (Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return Check(ctx)
}

func CurrentVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return FallbackVersion
}

func LatestVersion(ctx context.Context, client *http.Client, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "commitgen")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("latest version returned status %s", resp.Status)
	}

	var latest latestVersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&latest); err != nil {
		return "", err
	}
	version := latest.Version
	if version == "" {
		version = latest.TagName
	}
	if strings.TrimSpace(version) == "" {
		return "", fmt.Errorf("latest version response missing version")
	}
	return version, nil
}

func CompareVersions(left, right string) int {
	leftParts := versionParts(left)
	rightParts := versionParts(right)

	for i := 0; i < 3; i++ {
		if leftParts[i] > rightParts[i] {
			return 1
		}
		if leftParts[i] < rightParts[i] {
			return -1
		}
	}
	return 0
}

func versionParts(version string) [3]int {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")
	version = strings.Split(version, "-")[0]

	var parts [3]int
	for i, part := range strings.Split(version, ".") {
		if i >= len(parts) {
			break
		}
		n, err := strconv.Atoi(part)
		if err != nil {
			continue
		}
		parts[i] = n
	}
	return parts
}
