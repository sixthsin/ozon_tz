CREATE TABLE comments (
    id VARCHAR(36) PRIMARY KEY,
    post_id VARCHAR(36) NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    parent_id VARCHAR(36),
    author_id VARCHAR(36) NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comments_post_id_created_at ON comments(post_id, created_at);