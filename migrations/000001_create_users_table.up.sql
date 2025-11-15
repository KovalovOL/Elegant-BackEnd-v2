CREATE TABLE users (
    id UUID UNIQUE NOT NULL,
    name VARCHAR(100),
    email VARCHAR(50) UNIQUE NOT NULL,
    git_url TEXT,
    likedin_url TEXT,
    bio TEXT
);