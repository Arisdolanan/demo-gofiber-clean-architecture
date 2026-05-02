-- Migration: Create settings table
-- Version: 016
-- Description: Drop settings table

DROP TABLE IF EXISTS settings;

-- Migration: Create backups table (rollback)
-- Version: 018

DROP TABLE IF EXISTS backups;
