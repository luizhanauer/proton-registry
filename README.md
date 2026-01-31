# 🦀 Proton Registry
Automated indexer for GloriousEggroll's Proton-GE releases. Updates daily via GitHub Actions.

## Official Repository de Download

https://github.com/GloriousEggroll/proton-ge-custom/releases

O **Proton Registry** é um indexador automatizado para as versões do [Proton-GE](https://github.com/GloriousEggroll/proton-ge-custom). Ele foi criado para resolver problemas de *Rate Limit* na API do GitHub e fornecer um "backend" leve para aplicações de gerenciamento de Proton no Linux.

## 🚀 Como Funciona

1.  **Automação:** Uma GitHub Action roda diariamente à meia-noite (UTC).
2.  **Scraping Inteligente:** O indexador verifica se há novas releases. Se houver, ele varre todo o histórico do repositório original.
3.  **Filtragem (Smart Index):** O motor de indexação agrupa as versões por "Major Version" (Proton 10, 9, 8, 7, 6, 5, 4) e extrai apenas o que é essencial.
4.  **Entrega:** Os dados são salvos em arquivos JSON estáticos, servidos via GitHub Raw.

## 📦 API (Endpoints JSON)

Seu cliente ou script pode consumir os seguintes arquivos:

| Arquivo | Descrição | Uso Recomendado |
| :--- | :--- | :--- |
| `smart_index.json` | Top 10 recentes + última estável de cada versão anterior. | **Produção / Apps de usuário.** |
| `full_index.json` | Histórico completo de todas as versões já lançadas. | Arquivamento / Pesquisa. |

**Exemplo de URL para o `smart_index.json`:**
`https://raw.githubusercontent.com/luizhanauer/proton-registry/main/api/smart_index.json`

## 🛠️ Tecnologias

* **Linguagem:** Go (Golang)
* **CI/CD:** GitHub Actions
* **Storage:** Flat JSON files (GitHub Pages/Raw)

## 🔧 Compilação Local

Caso deseje rodar o indexador manualmente:

```bash
go build -ldflags="-s -w" -o proton-registry main.go
./proton-registry