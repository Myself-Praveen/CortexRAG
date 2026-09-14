# 🤖 CortexRAG
**High-Throughput Vector Search & RAG Gateway**

A lightweight, concurrent Vector Search Engine and RAG (Retrieval-Augmented Generation) gateway written in Go that ingests documents, indexes embeddings, performs vector similarity searches, and streams LLM answers.

## Features
- **In-Memory Vector Search Engine**: Implements brute-force and HNSW vector search with cosine similarity computed concurrently across CPU cores using worker pools.
- **Smart Embedding Cache**: Multi-tier LRU cache in Go to prevent duplicate expensive embedding calls to OpenAI / Gemini / Ollama for frequent queries.
- **SSE Streaming**: Delivers streaming token generation to the frontend with sub-50ms Time-To-First-Token (TTFT).
- **Document Ingestion Pipeline**: Extracts text from PDFs and Markdown, chunks it iteratively, and embeds it via an LLM API.
- **Pluggable Architecture**: Dependency Injection for Vector Store (Memory vs HNSW) and LLM Provider (Gemini vs Ollama).
- **React UI**: A polished, dark-mode frontend built with Vite, React, and Tailwind CSS.

## Getting Started

### Backend
1. Copy `.env.example` to `.env` and configure your API keys (e.g., `GEMINI_API_KEY`).
2. Run `go run ./cmd/cortex --port 8080`

### Frontend
1. Navigate to the `ui` directory.
2. Run `npm install`
3. Run `npm run dev`

## Architecture
- `internal/extract`: PDF and text parsing
- `internal/embedding`: LLM Embedding interfaces and Cache
- `internal/vectorstore`: In-memory and HNSW vector stores
- `internal/rag`: Orchestration pipeline (chunking -> embedding -> search -> generation)
- `internal/llm`: Generation providers (Gemini, Ollama)
- `internal/server`: Chi-based HTTP router and SSE handlers
- `ui/`: React frontend

## Built With
- **Go 1.22+**: Core backend logic and concurrency
- **google.golang.org/genai**: Official Gemini Go SDK
- **Chi**: Lightweight routing
- **React + Vite + Tailwind**: Frontend
