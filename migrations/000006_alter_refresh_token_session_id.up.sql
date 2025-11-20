ALTER TABLE refresh_tokens
    ALTER COLUMN session_id DROP DEFAULT;

ALTER TABLE refresh_tokens
    ALTER COLUMN session_id TYPE UUID
    USING gen_random_uuid();

ALTER TABLE refresh_tokens
    ALTER COLUMN session_id SET DEFAULT gen_random_uuid();
