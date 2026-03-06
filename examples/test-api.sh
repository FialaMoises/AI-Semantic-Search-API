#!/bin/bash

# Script para testar a API

API_URL="http://localhost:8080"

echo "=== Testing AI Semantic Search API ==="
echo ""

# 1. Health Check
echo "1. Health Check..."
curl -s ${API_URL}/health | jq .
echo ""

# 2. Index Documents
echo "2. Indexing documents..."
curl -s -X POST ${API_URL}/api/v1/index-documents \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      {
        "id": "doc1",
        "text": "Go é uma linguagem de programação criada pelo Google, focada em simplicidade e performance",
        "metadata": {"category": "programming", "language": "go"}
      },
      {
        "id": "doc2",
        "text": "Python é muito usado para ciência de dados, machine learning e inteligência artificial",
        "metadata": {"category": "programming", "language": "python"}
      },
      {
        "id": "doc3",
        "text": "JavaScript é a linguagem principal da web, usada no frontend e backend",
        "metadata": {"category": "programming", "language": "javascript"}
      },
      {
        "id": "doc4",
        "text": "FAISS é uma biblioteca do Facebook para busca eficiente de vetores",
        "metadata": {"category": "library", "topic": "vector-search"}
      },
      {
        "id": "doc5",
        "text": "Docker permite empacotar aplicações em containers portáteis",
        "metadata": {"category": "devops", "topic": "containers"}
      }
    ]
  }' | jq .
echo ""

# 3. Search Test 1
echo "3. Search: 'inteligência artificial'"
curl -s -X POST ${API_URL}/api/v1/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": "inteligência artificial e machine learning",
    "top_k": 3
  }' | jq .
echo ""

# 4. Search Test 2
echo "4. Search: 'busca vetorial'"
curl -s -X POST ${API_URL}/api/v1/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": "busca vetorial e similaridade",
    "top_k": 3
  }' | jq .
echo ""

# 5. Search Test 3
echo "5. Search: 'linguagem rápida e performática'"
curl -s -X POST ${API_URL}/api/v1/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": "linguagem rápida e performática para backend",
    "top_k": 3
  }' | jq .
echo ""

echo "=== Tests completed ==="
