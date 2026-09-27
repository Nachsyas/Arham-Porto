-- Sprint 1: GitHub Repository Ingestion Foundation
-- Stores synchronized GitHub repository metadata for Nachsyas
CREATE TABLE IF NOT EXISTS github_projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    github_id BIGINT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    full_name TEXT NOT NULL,
    description TEXT,
    html_url TEXT,
    homepage TEXT,
    language TEXT,
    stars INTEGER DEFAULT 0,
    forks INTEGER DEFAULT 0,
    topics JSONB DEFAULT '[]',
    readme TEXT,
    preview_image TEXT,
    synced_at TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Fast lookup indexes
CREATE INDEX IF NOT EXISTS idx_github_projects_github_id ON github_projects(github_id);
CREATE INDEX IF NOT EXISTS idx_github_projects_updated_at ON github_projects(updated_at);
