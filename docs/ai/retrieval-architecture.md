# Retrieval Architecture & Vector Search Specification (Phase 5 Certified)

## 1. Overview
The internal retrieval service (`internal/retrieval/service.go`) provides semantic evidence discovery using exact nearest-neighbor search over pgvector embeddings.

---

## 2. Embedding Identity & Vector Space Segregation
Per Gate #2:
- Embeddings produced by different models must **never** be compared as though they belong to one vector space.
- Every indexed chunk records:
  - `embedding_provider`
  - `embedding_model`
  - `embedding_dimensions`
- Retrieval queries filter by active model identity before vector distance ranking:
  ```sql
  WHERE c.embedding IS NOT NULL
    AND c.embedding_provider = $provider
    AND c.embedding_model = $model
    AND c.embedding_dimensions = $dimensions
    AND ($repoFilter = '' OR s.repository = $repoFilter)
    AND ($typeFilter = '' OR s.source_type = $typeFilter)
  ORDER BY c.embedding <=> $queryVector ASC
  LIMIT $limit
  ```
- Changing models requires re-embedding the corpus. Vectors from different providers or models are never mixed.

---

## 3. Explicit Embedding Purpose Contract
Per Gate #3:
- Document embeddings and query embeddings have distinct semantic roles.
- The provider abstraction mandates explicit purpose:
  ```go
  type Purpose string
  const (
      PurposeDocument Purpose = "document"
      PurposeQuery    Purpose = "query"
  )
  ```
- Providers map this purpose into native modes:
  - Gemini: `RETRIEVAL_DOCUMENT` vs `RETRIEVAL_QUERY`.
  - Fake / local: purpose-aware seed hashing.

---

## 4. Vector Dimension Contract
Per Gate #4:
- Configured deployment dimension: `768`.
- Supported default model: `gemini-embedding-2` with explicit `outputDimensionality: 768`.
- Every returned vector is strictly validated:
  - `len(v) == 768`
  - Zero NaN values
  - Zero Infinity values
- No truncation, padding, or silent coercion is permitted.

---

## 5. Exact Search First (No HNSW / ANN)
Per Gate #5, #37:
- Initial portfolio corpus contains dozens to hundreds of chunks.
- pgvector exact nearest-neighbor cosine search (`<=>`) is deterministic, provides perfect recall, and eliminates premature ANN configuration overhead.
- HNSW or IVFFlat indexes are explicitly deferred until corpus size or real query latency warrants approximate indexing.

---

## 6. Score & Distance Semantics
Per Gate #41:
- `Distance`: Cosine distance returned directly by pgvector (`c.embedding <=> $queryVector`), ranging from `0.0` (identical) to `2.0` (opposite).
- `Similarity`: Normalized cosine similarity defined as `1.0 - Distance`.
- Top-K limits: Default `5`, maximum `20`.

---

## 7. Immutable Citation Generation
Per Gate #38:
- Immutable GitHub URLs are constructed strictly from validated metadata:
  ```
  https://github.com/<owner>/<repo>/blob/<40-char-commit-sha>/<clean-path>
  ```
- Arbitrary URLs or absolute local paths are never embedded in citations.

---

## 8. Phase 6 Boundary & Public Surface Freeze
Per Gate #42, #43, #64:
- **No Public HTTP Routes**: No `/api/v1/search`, `/api/v1/retrieve`, or `/api/v1/ai` endpoints.
- **No Generative AI**: Zero chat completions, summarizations, or reviewer brief generation.
- Retrieval remains an internal package and CLI tool.
