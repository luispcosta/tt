CREATE TABLE activities (
  id integer PRIMARY KEY AUTOINCREMENT,
  name string NOT NULL,
  alias string,
  description text
);

CREATE UNIQUE INDEX name_alias_index
ON activities(name, alias);

CREATE UNIQUE INDEX name_index
ON activities(name);

CREATE INDEX alias_index
ON activities(alias);