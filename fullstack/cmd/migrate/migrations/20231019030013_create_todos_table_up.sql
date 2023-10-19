-- Write your UP migration SQL here.
CREATE TABLE IF NOT EXISTS todos (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    checked INTEGER NOT NULL CHECK(checked IN (0, 1))
);
