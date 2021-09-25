# Go-timetracking

This is a pet project to quickly track activities and see how much time you spend doing things.

# Requirements

* Golang
* SQLite3

# Setting uo SQLITE 3

* Make sure you have `sqlite3` installed.
* Init the application with `tt init`.
* From the root folder of this project, run `sqlite3 ~/.gott/gott.db` in your terminal to connect to a sqlite3 shell to the database.
* Run `.read migrations/001_create_activities.up.sql` inside the `sqlite3` prompt.
* Run `.read migrations/002_create_activity_logs.up.sql` inside the sqlite3 prompt.

You now have all the tables created. If you need to drop the activities table:

* `.read migrations/001_create_activities.down.sql`

If you need to drop the activity logs table:

* `.read migrations/002_create_activity_logs.down.sql`

# Commands

To see a list of all the supported commands and how to use them, please run `tt help`. You can also
run `tt COMMAND --help`
