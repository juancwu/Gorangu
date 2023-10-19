-- Write your DOWN migration SQL here.
-- SQLite does not support drop column
CREATE TABLE todos_tmp (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    checked INTEGER NOT NULL CHECK(checked IN (0,1))
);

INSERT INTO todos_tmp (id, name, checked)
SELECT id, name, checked FROM todos;

DROP TABLE todos;

ALTER TABLE todos_tmp RENAME TO todos;
