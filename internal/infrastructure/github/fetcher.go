package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/luizhanauer/proton-registry/internal/domain"
)

type githubRelease struct {
	TagName     string        `json:"tag_name"`
	PublishedAt time.Time     `json:"published_at"`
	Assets      []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

type Fetcher struct {
	baseURL string
	client  *http.Client
}

// NewFetcher agora recebe a baseURL como dependência
func NewFetcher(baseURL string) *Fetcher {
	return &Fetcher{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func (f *Fetcher) GetLatestTagName() (string, error) {
	resp, err := f.client.Get(f.baseURL + "/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var latest githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&latest); err != nil {
		return "", err
	}
	return latest.TagName, nil
}

func (f *Fetcher) FetchAll() (domain.ReleaseCollection, error) {
	var releases []domain.Release
	page := 1

	for {
		pageReleases, hasMore := f.fetchPage(page)
		if !hasMore {
			break
		}
		releases = append(releases, pageReleases...)
		page++
		time.Sleep(200 * time.Millisecond) // Respeita rate limit
	}

	return domain.ReleaseCollection{Releases: releases}, nil
}

func (f *Fetcher) fetchPage(page int) ([]domain.Release, bool) {
	url := fmt.Sprintf("%s/releases?per_page=100&page=%d", f.baseURL, page)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Proton-Registry-Bot")

	resp, err := f.client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		f.closeResp(resp)
		return nil, false
	}
	defer resp.Body.Close()

	var remotes []githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&remotes); err != nil {
		return nil, false
	}

	if len(remotes) == 0 {
		return nil, false
	}

	return f.mapToDomain(remotes), true
}

func (f *Fetcher) mapToDomain(remotes []githubRelease) []domain.Release {
	var parsed []domain.Release
	for _, r := range remotes {
		if release, valid := f.extractValidRelease(r); valid {
			parsed = append(parsed, release)
		}
	}
	return parsed
}

func (f *Fetcher) extractValidRelease(r githubRelease) (domain.Release, bool) {
	url, size := f.findTarGz(r.Assets)
	if url == "" {
		return domain.Release{}, false
	}

	entry := domain.Release{
		Version:     r.TagName,
		DownloadURL: url,
		Size:        size,
		Date:        r.PublishedAt.Format("2006-01-02"),
		Major:       f.getMajorVersion(r.TagName),
	}
	return entry, true
}

func (f *Fetcher) findTarGz(assets []githubAsset) (string, int64) {
	for _, a := range assets {
		if f.isValidTarGz(a.Name) {
			return a.BrowserDownloadURL, a.Size
		}
	}
	return "", 0
}

func (f *Fetcher) isValidTarGz(name string) bool {
	return strings.HasSuffix(name, ".tar.gz") && !strings.Contains(name, "sha512")
}

func (f *Fetcher) getMajorVersion(tag string) string {
	if strings.HasPrefix(tag, "GE-Proton") {
		return f.parseNewStandard(tag)
	}
	if len(tag) > 0 && tag[0] >= '0' && tag[0] <= '9' {
		return f.parseOldStandard(tag)
	}
	return "Outros"
}

func (f *Fetcher) parseNewStandard(tag string) string {
	parts := strings.Split(tag, "-")
	if len(parts) >= 2 {
		return parts[1]
	}
	return "Outros"
}

func (f *Fetcher) parseOldStandard(tag string) string {
	parts := strings.Split(tag, ".")
	if len(parts) >= 1 {
		return "Proton" + parts[0]
	}
	return "Outros"
}

func (f *Fetcher) closeResp(resp *http.Response) {
	if resp != nil {
		resp.Body.Close()
	}
}
