CREATE TABLE auth_providers (
    provider TEXT NOT NULL,
    provider_user_id TEXT,
    user_id UUID NOT NULL,
    created_at TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);