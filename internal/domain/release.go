package domain

// Release representa a entidade central do nosso domínio.
type Release struct {
	Version     string `json:"version"`
	DownloadURL string `json:"url"`
	Size        int64  `json:"size"`
	Date        string `json:"date"`
	Major       string `json:"major"`
}

// ReleaseCollection é uma First-Class Collection para aplicar regras sobre listas.
type ReleaseCollection struct {
	Releases []Release
}

func (c *ReleaseCollection) IsEmpty() bool {
	return len(c.Releases) == 0
}

func (c *ReleaseCollection) First() Release {
	if c.IsEmpty() {
		return Release{}
	}
	return c.Releases[0]
}
