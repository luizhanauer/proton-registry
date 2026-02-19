# 🎮 Proton Registry

![Go Version](https://img.shields.io/badge/Go-1.25.6-00ADD8?style=flat&logo=go)
![GitHub Actions](https://img.shields.io/badge/Build-Passing-brightgreen?style=flat&logo=githubactions)
![License](https://img.shields.io/badge/License-MIT-blue.svg)

**Proton Registry** é uma API estática de alto desempenho projetada para indexar, filtrar e servir as versões do [Proton-GE Custom](https://github.com/GloriousEggroll/proton-ge-custom).

Este projeto atua como um *middleware* leve, consumindo a API original do GitHub, aplicando regras de negócio para otimização de payload (Smart Filter) e publicando os resultados estáticos no GitHub Pages. Isso garante downloads rápidos, zero custo de servidor e evita bloqueios por *rate limit*.

---

## ❓ O que é o Proton e o Proton-GE?

**Proton** é uma camada de compatibilidade desenvolvida pela Valve (baseada no Wine e outras tecnologias) que permite que jogos nativos de Windows sejam executados no Linux de forma transparente e com alta performance — é a tecnologia central que faz o **Steam Deck** funcionar.

O **Proton-GE Custom** (mantido por *GloriousEggroll*) é um *fork* comunitário focado em entregar as últimas novidades. Ele inclui patches experimentais, correções avançadas de codecs de áudio/vídeo e suporte *Day One* para lançamentos recentes que ainda não foram integrados à versão oficial do Proton da Valve.

---

## 🌐 Endpoints da API

A API é estática e atualizada diariamente. Você pode consumi-la via requisições `GET` simples.

### 1. Smart Index (Recomendado)
Mantém as 10 versões mais recentes (*bleeding edge*) e apenas a última versão estável de cada *major release* anterior (ex: Proton9, Proton8, etc). Ideal para clientes que precisam de economia de banda e processamento.
```http
GET https://luizhanauer.github.io/proton-registry/api/smart_index.json
```

### 2. Full Index

Contém o histórico completo de todas as versões extraídas do repositório original.

```http
GET https://luizhanauer.github.io/proton-registry/api/full_index.json
```

### Exemplo de Resposta (JSON)

```json
[
  {
    "version": "GE-Proton10-32",
    "url": "[https://github.com/GloriousEggroll/proton-ge-custom/releases/download/GE-Proton10-32/GE-Proton10-32.tar.gz](https://github.com/GloriousEggroll/proton-ge-custom/releases/download/GE-Proton10-32/GE-Proton10-32.tar.gz)",
    "size": 515874267,
    "date": "2026-02-16",
    "major": "Proton10"
  }
]

```

---

## 🏗️ Arquitetura e Engenharia

O núcleo do projeto foi desenvolvido em **Go** utilizando princípios avançados de engenharia de software:

* **Clean Architecture & DDD:** O código está dividido em camadas de Domínio (`domain`), Casos de Uso (`usecase`) e Infraestrutura (`infrastructure`), isolando as regras de negócio de detalhes de implementação (como APIs externas ou manipulação de arquivos).
* **Object Calisthenics:** Foco em baixo nível de indentação, *early returns*, ausência do uso indiscriminado de `else` e encapsulamento em *First-Class Collections*.
* **Inversão de Dependência (SOLID):** Uso de interfaces para comunicação entre camadas, permitindo 100% de cobertura de testes unitários com *mocks* locais.
* **CI/CD via GitHub Actions:** Um bot automatizado compila o código, roda os testes, gera os arquivos JSON na pasta `public/` e faz o deploy seguro utilizando *Artifacts* diretamente para o GitHub Pages.

---

## 💻 Desenvolvimento Local

Se você deseja rodar ou modificar o projeto localmente:

### Pré-requisitos

* Go 1.25+

### Executando

1. Clone o repositório:
```bash
git clone https://github.com/luizhanauer/proton-registry.git
```

2. Execute o orquestrador:
```bash
cd proton-registry
```

```bash
go run ./cmd/registry

```

*Os arquivos `full_index.json` e `smart_index.json` serão gerados dentro da pasta `public/api/`.*

### Rodando os Testes

A suíte de testes cobre as regras do domínio, os fluxos do orquestrador e as simulações de I/O e HTTP.

```bash
go test -v ./... --cover

```

---

## 🙏 Créditos e Agradecimentos

* Todo o mérito do desenvolvimento e manutenção dos binários do **Proton-GE Custom** pertence ao projeto original [GloriousEggroll/proton-ge-custom](https://github.com/GloriousEggroll/proton-ge-custom). Este repositório é apenas um facilitador de indexação.

## ☕ Apoie o Projeto

Se esta ferramenta facilitou a sua vida ou ajudou a construir seus próprios projetos, considere pagar um café!

---

*Desenvolvido com 💜 por [Luiz Hanauer*](https://github.com/luizhanauer)

