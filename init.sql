CREATE TABLE admin (
    admin_id INTEGER PRIMARY KEY,
    id TEXT,
    pw TEXT
);
CREATE TABLE story (
    story_id INTEGER PRIMARY KEY,
    id TEXT,
    title TEXT,
    date_marked TEXT
    primary_media_id TEXT,
    content TEXT
);
CREATE TABLE media (
    media_id INTEGER PRIMARY KEY,
    id TEXT,
    kind TEXT,
    content BLOB
)