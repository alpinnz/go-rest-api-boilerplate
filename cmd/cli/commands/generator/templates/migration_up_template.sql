-- Create PLACEHOLDER_entities table
CREATE TABLE IF NOT EXISTS PLACEHOLDER_entities (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    -- Add your columns here
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);

-- Create indexes
CREATE INDEX idx_PLACEHOLDER_entities_deleted_at ON PLACEHOLDER_entities(deleted_at);
CREATE INDEX idx_PLACEHOLDER_entities_created_at ON PLACEHOLDER_entities(created_at);

