# 🤖 CortexRAG — High-Throughput Vector Search & RAG Gateway

A lightweight, concurrent Vector Search Engine and RAG (Retrieval-Augmented Generation) gateway written in Go. It ingests documents, indexes embeddings, performs vector similarity searches, and streams LLM answers via Server-Sent Events (SSE).

## Features
- **In-Memory Vector Search**: Concurrent cosine similarity computation using Goroutine worker pools.
- **Multi-tier LRU Cache**: Thread-safe cache to prevent redundant embedding API calls.
- **Document Ingestion**: Fast extraction and chunking of PDF and Markdown files.
- **Provider Agnostic**: Switchable backends (Gemini, Ollama, pgvector).
- **Streaming LLM Responses**: Real-time SSE token delivery.

## Architecture
See `docs/architecture.md` (coming soon).

## Quick Start
1. Copy `.env.example` to `.env` and configure your API keys.
2. Run `go run cmd/cortexrag/main.go`
3. Access the frontend (coming soon).

## License
MIT
