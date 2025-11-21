SET session_replication_role = 'replica';

DROP TABLE IF EXISTS pr_reviewers CASCADE;
DROP TABLE IF EXISTS pull_requests CASCADE;
DROP TABLE IF EXISTS teams CASCADE;
DROP TABLE IF EXISTS users CASCADE;

SET session_replication_role = 'origin';

CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    user_name VARCHAR(255) NOT NULL,
    team_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE teams (
    team_name VARCHAR(255) PRIMARY KEY
);

CREATE TABLE pull_requests (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    author_id VARCHAR(255) NOT NULL REFERENCES users(id),
    reviewers_id VARCHAR(255)[] NOT NULL DEFAULT '{}',
    is_merged BOOLEAN NOT NULL DEFAULT false
    merged_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX idx_users_team_name ON users(team_name);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_pull_requests_author_id ON pull_requests(author_id);