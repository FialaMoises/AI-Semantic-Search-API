# PowerShell script para testar a API (Windows)

$API_URL = "http://localhost:8080"

Write-Host "=== Testing AI Semantic Search API ===" -ForegroundColor Green
Write-Host ""

# 1. Health Check
Write-Host "1. Health Check..." -ForegroundColor Yellow
Invoke-RestMethod -Uri "$API_URL/health" -Method Get | ConvertTo-Json
Write-Host ""

# 2. Index Documents
Write-Host "2. Indexing documents..." -ForegroundColor Yellow
$indexBody = @{
    documents = @(
        @{
            id = "doc1"
            text = "Go é uma linguagem de programação criada pelo Google, focada em simplicidade e performance"
            metadata = @{ category = "programming"; language = "go" }
        },
        @{
            id = "doc2"
            text = "Python é muito usado para ciência de dados, machine learning e inteligência artificial"
            metadata = @{ category = "programming"; language = "python" }
        },
        @{
            id = "doc3"
            text = "JavaScript é a linguagem principal da web, usada no frontend e backend"
            metadata = @{ category = "programming"; language = "javascript" }
        },
        @{
            id = "doc4"
            text = "FAISS é uma biblioteca do Facebook para busca eficiente de vetores"
            metadata = @{ category = "library"; topic = "vector-search" }
        },
        @{
            id = "doc5"
            text = "Docker permite empacotar aplicações em containers portáteis"
            metadata = @{ category = "devops"; topic = "containers" }
        }
    )
} | ConvertTo-Json -Depth 10

Invoke-RestMethod -Uri "$API_URL/api/v1/index-documents" -Method Post -Body $indexBody -ContentType "application/json" | ConvertTo-Json
Write-Host ""

# 3. Search Test 1
Write-Host "3. Search: 'inteligência artificial'" -ForegroundColor Yellow
$searchBody1 = @{
    query = "inteligência artificial e machine learning"
    top_k = 3
} | ConvertTo-Json

Invoke-RestMethod -Uri "$API_URL/api/v1/search" -Method Post -Body $searchBody1 -ContentType "application/json" | ConvertTo-Json -Depth 10
Write-Host ""

# 4. Search Test 2
Write-Host "4. Search: 'busca vetorial'" -ForegroundColor Yellow
$searchBody2 = @{
    query = "busca vetorial e similaridade"
    top_k = 3
} | ConvertTo-Json

Invoke-RestMethod -Uri "$API_URL/api/v1/search" -Method Post -Body $searchBody2 -ContentType "application/json" | ConvertTo-Json -Depth 10
Write-Host ""

# 5. Search Test 3
Write-Host "5. Search: 'linguagem rápida'" -ForegroundColor Yellow
$searchBody3 = @{
    query = "linguagem rápida e performática para backend"
    top_k = 3
} | ConvertTo-Json

Invoke-RestMethod -Uri "$API_URL/api/v1/search" -Method Post -Body $searchBody3 -ContentType "application/json" | ConvertTo-Json -Depth 10
Write-Host ""

Write-Host "=== Tests completed ===" -ForegroundColor Green
