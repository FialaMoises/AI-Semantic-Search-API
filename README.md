# 🔍 AI Semantic Search API

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Python Version](https://img.shields.io/badge/Python-3.11+-3776AB?style=flat&logo=python&logoColor=white)](https://python.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker&logoColor=white)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Uma API de busca semântica inteligente usando embeddings vetoriais e busca por similaridade.

## O que ela faz

- **Recebe texto do usuário** via API REST
- **Converte em embedding** usando modelos de IA
- **Busca documentos semelhantes** no banco vetorial
- **Retorna respostas relevantes** ranqueadas por similaridade

## Stack Tecnológica

- **Go** - Backend API performático
- **Gin** - Framework web
- **Vector Search** - Busca por similaridade (cosine similarity)
- **sentence-transformers** - Embeddings locais (gratuito)
- **OpenAI API** - Embeddings premium (opcional)
- **Docker** - Containerização

## Arquitetura

```
┌─────────────┐
│   Cliente   │
└──────┬──────┘
       │ HTTP
       ▼
┌─────────────────────┐
│   API Go (Gin)      │
│  - /search          │
│  - /index-documents │
└──────┬──────────────┘
       │
       ├──────────────┐
       │              │
       ▼              ▼
┌─────────────┐  ┌──────────────┐
│  Embeddings │  │    Vector    │
│   Provider  │  │    Store     │
└─────────────┘  └──────────────┘
  │         │
  ▼         ▼
[Local]  [OpenAI]
```

## Estrutura do Projeto

```
ai-semantic-search-api/
├── api/                  # Entry point da aplicação
├── internal/
│   ├── api/              # Handlers e rotas HTTP
│   ├── config/           # Configurações
│   ├── embeddings/       # Provedores de embeddings
│   └── search/           # Vector store (cosine similarity)
├── embedding-service/    # Microserviço Python
│   ├── main.py
│   ├── requirements.txt
│   └── Dockerfile
├── data/                 # Índices vetoriais (persistência)
├── examples/             # Scripts de teste
├── docker-compose.yml
├── QUICKSTART.md
└── README.md
```

## Instalação e Setup

### Opção 1: Docker (Recomendado)

**Com embeddings locais (gratuito):**

```bash
# 1. Clone o repositório
git clone <seu-repo>
cd ai-semantic-search

# 2. Configure o ambiente
make setup
# ou manualmente:
cp .env.example .env
mkdir -p data

# 3. Inicie os serviços
docker-compose up -d

# 4. Verifique os logs
docker-compose logs -f
```

**Com OpenAI (pago):**

```bash
# 1. Configure sua API key
echo "OPENAI_API_KEY=sk-your-key-here" > .env

# 2. Inicie com OpenAI
docker-compose -f docker-compose.yml -f docker-compose.openai.yml up -d
```

### Opção 2: Local Development

**Pré-requisitos:**
- Go 1.21+
- Python 3.11+

```bash
# 1. Instale dependências Go
go mod download

# 2. Configure Python embedding service
cd embedding-service
pip install -r requirements.txt
python main.py &  # Roda em background na porta 8001

# 3. Configure .env
cd ..
cp .env.example .env

# 4. Rode a API
go run ./api
```

## Uso da API

### 1. Health Check

```bash
curl http://localhost:8080/health
```

**Resposta:**
```json
{
  "status": "healthy",
  "document_count": 0,
  "embedding_dimension": 384
}
```

### 2. Indexar Documentos

```bash
curl -X POST http://localhost:8080/api/v1/index-documents \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      {
        "id": "doc1",
        "text": "Go é uma linguagem de programação criada pelo Google",
        "metadata": {"category": "tech"}
      },
      {
        "id": "doc2",
        "text": "Python é muito usado para ciência de dados e IA",
        "metadata": {"category": "tech"}
      },
      {
        "id": "doc3",
        "text": "JavaScript é a linguagem da web",
        "metadata": {"category": "tech"}
      }
    ]
  }'
```

**Resposta:**
```json
{
  "message": "documents indexed successfully",
  "indexed_count": 3,
  "total_documents": 3
}
```

### 3. Buscar Documentos

```bash
curl -X POST http://localhost:8080/api/v1/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": "linguagem para inteligência artificial",
    "top_k": 3
  }'
```

**Resposta:**
```json
{
  "query": "linguagem para inteligência artificial",
  "results": [
    {
      "document": {
        "id": "doc2",
        "text": "Python é muito usado para ciência de dados e IA",
        "metadata": {"category": "tech"}
      },
      "score": 0.234
    },
    {
      "document": {
        "id": "doc1",
        "text": "Go é uma linguagem de programação criada pelo Google",
        "metadata": {"category": "tech"}
      },
      "score": 0.456
    }
  ],
  "count": 2
}
```

## Features

✅ **Embeddings flexíveis**
- Suporte local (gratuito) com sentence-transformers
- Suporte OpenAI (pago, melhor qualidade)
- Fácil alternância via env vars

✅ **Busca vetorial eficiente**
- Cosine similarity para busca semântica
- Persistência automática do índice

✅ **API REST completa**
- Endpoint de busca semântica
- Endpoint de indexação
- Health checks

✅ **Containerização**
- Docker Compose para deploy fácil
- Ambiente isolado e reproduzível

✅ **Arquitetura limpa**
- Interfaces para fácil extensão
- Separação de responsabilidades
- Código testável

## Configuração

### Variáveis de Ambiente

Edite o arquivo `.env`:

```env
# Provedor: "local" ou "openai"
EMBEDDING_PROVIDER=local

# Serviço local de embeddings
EMBEDDING_SERVICE_URL=http://localhost:8001

# OpenAI (apenas se EMBEDDING_PROVIDER=openai)
OPENAI_API_KEY=sk-xxxxx
OPENAI_MODEL=text-embedding-3-small

# API
API_PORT=8080

# Vector Store
FAISS_INDEX_PATH=./data/vector.index
FAISS_DIMENSION=384  # 384 para local, 1536 para OpenAI
```

## Comparação: Local vs OpenAI

| Aspecto | Local (sentence-transformers) | OpenAI |
|---------|-------------------------------|--------|
| **Custo** | Gratuito | ~$0.0001/1K tokens |
| **Qualidade** | Boa (85-90%) | Excelente (95%+) |
| **Latência** | Baixa (local) | Média (API externa) |
| **Privacidade** | Total | Dados enviados para OpenAI |
| **Dimensão** | 384 | 1536 |
| **Requer GPU** | Não (CPU ok) | N/A |

**Recomendação:** Comece com local para desenvolvimento e teste. Migre para OpenAI em produção se precisar de melhor qualidade.

## Comandos Úteis

```bash
# Build
make build

# Run local
make run

# Docker
make docker-up          # Local embeddings
make docker-up-openai   # OpenAI embeddings
make docker-down        # Stop services
make docker-logs        # Ver logs

# Limpar dados
make clean
```

## Troubleshooting

### Erro: "connection refused" ao embedding service

```bash
# Verifique se o serviço está rodando
docker-compose ps
curl http://localhost:8001/health

# Reinicie o serviço
docker-compose restart embedding-service
```

### Erro: "failed to create index"

```bash
# Verifique permissões da pasta data
mkdir -p data
chmod 755 data

# Limpe índice corrupto
rm -rf data/*
```

### OpenAI rate limit

Configure retry logic ou use `text-embedding-3-small` (mais barato).

## Próximos Passos

- [ ] Adicionar autenticação (API keys)
- [ ] Implementar rate limiting
- [ ] Cache de embeddings
- [ ] Suporte a filtros por metadata
- [ ] Paginação nos resultados
- [ ] Métricas e observabilidade
- [ ] Testes unitários e de integração
- [ ] CI/CD pipeline

## 📄 Licença

MIT License - sinta-se livre para usar em seus projetos!

## Contribuindo

Pull requests são bem-vindos! Para mudanças grandes, abra uma issue primeiro.

---

**Desenvolvido com Go, Python e IA** 🚀
