DROP TABLE IF EXISTS works CASCADE;
DROP TABLE IF EXISTS files CASCADE;

CREATE TABLE files (
    id UUID PRIMARY KEY,
    file_name TEXT NOT NULL,
    content_type TEXT NOT NULL
);

CREATE TABLE works (
    id UUID PRIMARY KEY,
    user_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    type_work TEXT NOT NULL,
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE
);
