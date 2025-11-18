ALTER TABLE refresh_tokens
ALTER COLUMN ip TYPE inet USING ip::inet;
