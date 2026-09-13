-- Phase 0 Initial Database Migration: Enable pgvector extension
-- Canonical portfolio data lives in data/ for Phases 0-3.
-- Knowledge tables will be added in Phase 5.
CREATE EXTENSION IF NOT EXISTS vector;
