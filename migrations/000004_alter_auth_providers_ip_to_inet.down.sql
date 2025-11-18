ALTER TABLE refresh_tokens
ALTER COLUMN ip TYPE text USING ip::text;
