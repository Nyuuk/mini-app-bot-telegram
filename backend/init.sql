-- Initialize database for Mini App Bot Telegram
-- This script will be executed when the PostgreSQL container starts

-- Create database if not exists (already handled by POSTGRES_DB env var)
-- CREATE DATABASE IF NOT EXISTS postgres;

-- You can add any initial database setup here
-- For example, creating additional users, setting up extensions, etc.

-- Create extensions if needed
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Initial setup completed
SELECT 'Database initialized successfully' as status;
