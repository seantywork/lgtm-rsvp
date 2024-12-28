CREATE TABLE admin (
    admin_id INTEGER PRIMARY KEY,
    id TEXT,
    session_id TEXT,
    pw TEXT
);
CREATE TABLE story (
    story_id INTEGER PRIMARY KEY,
    id TEXT,
    title TEXT,
    date_marked TEXT
    primary_media_name TEXT,
    content TEXT
);

