CREATE TABLE IF NOT EXISTS USERS
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    bio      TEXT DEFAULT '',
    CONSTRAINT username_not_empty CHECK ( username <> '' )
); 