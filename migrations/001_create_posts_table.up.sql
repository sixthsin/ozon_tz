CREATE TABLE posts (
    id VARCHAR(36) PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author_id VARCHAR(36) NOT NULL,
    allow_comments BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_posts_created_at ON posts(created_at);