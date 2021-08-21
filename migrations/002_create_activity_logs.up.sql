CREATE TABLE activity_logs (
  id integer PRIMARY KEY AUTOINCREMENT,
  day date,
  started_at timestamp,
  stopped_at timestamp,  
  activity_id integer,
  FOREIGN KEY(activity_id) REFERENCES activities(id)
);

CREATE INDEX date_index
ON activity_logs(day);
