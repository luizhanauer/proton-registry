package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Estrutura do JSON final
type ProtonEntry struct {
	Version     string `json:"version"`
	DownloadURL string `json:"url"`
	Size        int64  `json:"size"`
	Date        string `json:"date"`
	Major       string `json:"major"`
}

// Estruturas da API do GitHub
type Release struct {
	TagName     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []Asset   `json:"assets"`
}
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

const (
	repoOwner      = "GloriousEggroll"
	repoName       = "proton-ge-custom"
	apiBase        = "https://api.github.com/repos/" + repoOwner + "/" + repoName
	fullIndexFile  = "api/full_index.json"
	smartIndexFile = "api/smart_index.json"
)

func main() {
	fmt.Println("🤖 Iniciando Proton Registry...")

	// Garante que a pasta api/ existe
	if err := os.MkdirAll("api", 0755); err != nil {
		fmt.Printf("❌ Erro ao criar pasta api: %v\n", err)
		return
	}

	client := &http.Client{Timeout: 60 * time.Second}

	if !needsUpdate(client) {
		fmt.Println("✅ O registro já está atualizado. Encerrando.")
		return
	}

	fmt.Println("🚀 Nova versão detectada! Atualizando índices...")
	fullList := scrapeAllReleases(client)

	saveJSON(fullIndexFile, fullList)
	smartList := generateSmartList(fullList)
	saveJSON(smartIndexFile, smartList)

	fmt.Println("🎉 Tudo pronto! Arquivos na pasta api/ atualizados.")
}

// --- Lógica de Verificação (Economia de Recursos) ---
func needsUpdate(client *http.Client) bool {
	// Pega a versão mais recente do GitHub (apenas 1 request leve)
	resp, err := client.Get(apiBase + "/releases/latest")
	if err != nil {
		fmt.Printf("⚠️ Erro ao checar latest: %v. Forçando update.\n", err)
		return true
	}
	defer resp.Body.Close()

	var latestRemote Release
	if err := json.NewDecoder(resp.Body).Decode(&latestRemote); err != nil {
		return true
	}

	// Tenta ler o arquivo local existente (trazido pelo git checkout)
	file, err := os.Open(fullIndexFile)
	if os.IsNotExist(err) {
		return true // Arquivo não existe, precisa criar
	}
	defer file.Close()

	var localIndex []ProtonEntry
	if err := json.NewDecoder(file).Decode(&localIndex); err != nil || len(localIndex) == 0 {
		return true
	}

	// COMPARAÇÃO: Se a tag remota for igual à primeira local, não mudou nada.
	if latestRemote.TagName == localIndex[0].Version {
		fmt.Printf("⏸️ Versão atual (%s) é idêntica à remota. Nenhuma ação necessária.\n", latestRemote.TagName)
		return false
	}

	fmt.Printf("🆕 Nova versão encontrada: %s (Local era: %s)\n", latestRemote.TagName, localIndex[0].Version)
	return true
}

// --- Lógica de Filtragem Inteligente (Otimização para o Cliente) ---
func generateSmartList(all []ProtonEntry) []ProtonEntry {
	var filtered []ProtonEntry
	seenMajors := make(map[string]bool)
	const keepRecent = 10

	limit := keepRecent
	if len(all) < limit {
		limit = len(all)
	}

	// 1. Bleeding Edge: Mantém as 10 últimas, não importa o quê
	for i := 0; i < limit; i++ {
		filtered = append(filtered, all[i])
		if all[i].Major != "" {
			seenMajors[all[i].Major] = true
		}
	}

	// 2. Legado Estável: Mantém apenas a ÚLTIMA de cada major version anterior
	for i := limit; i < len(all); i++ {
		major := all[i].Major
		// Se encontrarmos um Major (ex: Proton7) que AINDA não vimos...
		// significa que essa é a versão mais recente desse major (pois a lista está ordenada)
		if major != "" && major != "Outros" && !seenMajors[major] {
			filtered = append(filtered, all[i])
			seenMajors[major] = true
		}
	}

	fmt.Printf("🧠 Smart Index gerado: Reduzido de %d para %d entradas.\n", len(all), len(filtered))
	return filtered
}

// --- Scraper (O mesmo de antes, encapsulado) ---
func scrapeAllReleases(client *http.Client) []ProtonEntry {
	var index []ProtonEntry
	page := 1

	for {
		url := fmt.Sprintf("%s/releases?per_page=100&page=%d", apiBase, page)
		fmt.Printf("📄 Lendo página %d...\n", page)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Proton-Indexer-Bot")

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			if resp != nil {
				resp.Body.Close()
			}
			break
		}

		var releases []Release
		if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
			resp.Body.Close()
			break
		}
		resp.Body.Close()

		if len(releases) == 0 {
			break
		} // Fim das páginas

		for _, r := range releases {
			url, size := findTarGz(r.Assets)
			if url == "" {
				continue
			}

			index = append(index, ProtonEntry{
				Version:     r.TagName,
				DownloadURL: url,
				Size:        size,
				Date:        r.PublishedAt.Format("2006-01-02"),
				Major:       getMajorVersion(r.TagName),
			})
		}
		page++
		time.Sleep(200 * time.Millisecond)
	}
	return index
}

// --- Helpers ---
func saveJSON(filename string, data interface{}) {
	file, _ := os.Create(filename)
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(data)
}

func findTarGz(assets []Asset) (string, int64) {
	for _, a := range assets {
		if strings.HasSuffix(a.Name, ".tar.gz") && !strings.Contains(a.Name, "sha512") {
			return a.BrowserDownloadURL, a.Size
		}
	}
	return "", 0
}

func getMajorVersion(tag string) string {
	// Caso 1: Padrão Novo (GE-Proton10-29) -> Retorna "Proton10"
	if strings.HasPrefix(tag, "GE-Proton") {
		parts := strings.Split(tag, "-")
		if len(parts) >= 2 {
			return parts[1]
		}
	}

	// Caso 2: Padrão Antigo (7.3-GE-1, 6.21-GE-2) -> Retorna "Proton7", "Proton6", etc.
	// Se começar com um número, pegamos tudo antes do primeiro ponto.
	if len(tag) > 0 && tag[0] >= '0' && tag[0] <= '9' {
		parts := strings.Split(tag, ".")
		if len(parts) >= 1 {
			return "Proton" + parts[0]
		}
	}

	// Caso de fallback (segurança)
	return "Outros"
}
