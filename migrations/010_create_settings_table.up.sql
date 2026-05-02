-- Migration: Create settings table
-- Version: 016
-- Description: Create settings table for flexible configuration

CREATE TABLE IF NOT EXISTS settings (
    id BIGSERIAL PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    setting_key VARCHAR(255) NOT NULL,
    setting_value TEXT,
    group_name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    updated_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    deleted_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(school_id, setting_key)
);

-- Create index for group_name
CREATE INDEX IF NOT EXISTS idx_settings_group_name ON settings(group_name) WHERE deleted_at IS NULL;

-- Create index for setting_key
CREATE INDEX IF NOT EXISTS idx_settings_key ON settings(setting_key) WHERE deleted_at IS NULL;

-- Create index for school_id
CREATE INDEX IF NOT EXISTS idx_settings_school_id ON settings(school_id);

-- Create index for soft delete
CREATE INDEX IF NOT EXISTS idx_settings_deleted_at ON settings(deleted_at);

-- Add comment to table
COMMENT ON TABLE settings IS 'Settings table for flexible system configurations';
COMMENT ON COLUMN settings.id IS 'Primary key, auto-incrementing';
COMMENT ON COLUMN settings.setting_key IS 'Unique key identifying the configuration';
COMMENT ON COLUMN settings.setting_value IS 'Value of the configuration';
COMMENT ON COLUMN settings.group_name IS 'Category grouping (e.g., umum, backup, integrasi)';
COMMENT ON COLUMN settings.description IS 'Description of what the configuration does';


-- Migration: Create backups table
-- Version: 018

CREATE TABLE IF NOT EXISTS backups (
    id BIGSERIAL PRIMARY KEY,
    school_id BIGINT NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    storage_path TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,
    status VARCHAR(50) DEFAULT 'success',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_backups_school_id ON backups(school_id);



-- 101. Create Integration Definitions Table
-- This table stores metadata for integrations (name, provider, category)

CREATE TABLE IF NOT EXISTS integration_definitions (
    id BIGSERIAL PRIMARY KEY,
    school_id BIGINT REFERENCES schools(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    provider VARCHAR(100),
    category VARCHAR(50) NOT NULL, -- messaging, payment, meeting, calendar, other
    description TEXT,
    is_system BOOLEAN DEFAULT false, -- system integrations cannot be deleted
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(school_id, code)
);

-- Index for faster queries
CREATE INDEX IF NOT EXISTS idx_integration_definitions_school_id ON integration_definitions(school_id);
CREATE INDEX IF NOT EXISTS idx_integration_definitions_category ON integration_definitions(category);

-- Trigger for updated_at
CREATE TRIGGER update_integration_definitions_updated_at 
    BEFORE UPDATE ON integration_definitions 
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- Add comments
COMMENT ON TABLE integration_definitions IS 'Stores integration metadata definitions';
COMMENT ON COLUMN integration_definitions.code IS 'Unique integration code (e.g. whatsapp, midtrans)';
COMMENT ON COLUMN integration_definitions.provider IS 'Provider name (e.g. Wablas, Midtrans)';
COMMENT ON COLUMN integration_definitions.category IS 'Integration category: messaging, payment, meeting, calendar, other';
COMMENT ON COLUMN integration_definitions.is_system IS 'System integrations cannot be deleted by users';