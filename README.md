# Go-timetracking

This is a pet project to quickly track activities and see how much time you spend doing things.

# Requirements

* Golang
* SQLite3

# Running SQLITE 3 migrations

* Make sure you have `sqlite3` installed.
* From the root folder of this project, run `sqlite3 gott.db` in your terminal to create a new database called `gott`.
* Run `.read migrations/001_create_activities.up.sql` inside the `sqlite3` prompt.
* Run `.read migrations/002_create_activity_logs.up.sql` inside the sqlite3 prompt.

You now have all the tables created. If you need to drop the activities table:

* `.read migrations/001_create_activities.down.sql`

If you need to drop the activity logs table:

* `.read migrations/002_create_activity_logs.down.sql`
