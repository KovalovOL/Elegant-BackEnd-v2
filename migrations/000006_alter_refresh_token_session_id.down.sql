ALTER TABLE refresh_tokens
    ALTER COLUMN session_id DROP DEFAULT;

CREATE SEQUENCE refresh_tokens_session_id_seq;

ALTER TABLE refresh_tokens
    ALTER COLUMN session_id TYPE INTEGER
    USING nextval('refresh_tokens_session_id_seq');

ALTER TABLE refresh_tokens
    ALTER COLUMN session_id SET DEFAULT nextval('refresh_tokens_session_id_seq');
