CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    event_time TIMESTAMP NOT NULL,
    event_type VARCHAR(20) NOT NULL CHECK (event_type IN ('view', 'like', 'comment')),
    post_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    content TEXT
);

CREATE INDEX idx_events_post ON events(post_id, event_type);
CREATE INDEX idx_events_time ON events(event_time);
