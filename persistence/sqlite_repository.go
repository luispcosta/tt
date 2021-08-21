package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/utils"
	_ "github.com/mattn/go-sqlite3"
)

// SqliteRepository represents the connection to a SQLite database.
type SqliteRepository struct {
	db     *sql.DB
	dbFile string
	Clock  utils.Clock
}

// NewSqliteRepository creates a new SQLite repository struct
func NewSqliteRepository() (*SqliteRepository, error) {
	return &SqliteRepository{
		dbFile: "./gott.db",
		Clock:  utils.NewLiveClock(),
	}, nil
}

// Initialize initializes the connection to the database
func (repo *SqliteRepository) Initialize() error {
	db, err := sql.Open("sqlite3", repo.dbFile)

	if err != nil {
		return nil
	}

	repo.db = db
	return nil
}

// Shutdown shutsdown the database
func (repo *SqliteRepository) Shutdown() error {
	return repo.db.Close()
}

// TODO: Adding two activities with the same name.
// Add adds a new activity to the database
func (repo *SqliteRepository) Add(activity core.Activity) error {
	sql := fmt.Sprintf("INSERT INTO activities (name, alias, description) VALUES ('%s', '%s', '%s')", activity.Name, activity.Alias, activity.Description)

	_, err := repo.db.Exec(sql)

	if err != nil {
		return err
	}

	return nil
}

// Delete deletes an activity from the database
func (repo *SqliteRepository) Delete(activityNameOrAlias string) error {
	sql := fmt.Sprintf("DELETE FROM activities WHERE name = '%s' OR alias = '%s'", activityNameOrAlias, activityNameOrAlias)

	res, err := repo.db.Exec(sql)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("activity not found")
	}

	return nil
}

// List returns a list with all the activities in the database
func (repo *SqliteRepository) List() ([]core.Activity, error) {
	rows, err := repo.db.Query("SELECT name, alias, description FROM activities")

	if err != nil {
		return []core.Activity{}, err
	}

	defer rows.Close()

	var activities []core.Activity
	for rows.Next() {
		var activityId int
		var activityName string
		var activityAlias string
		var activityDesc string
		err = rows.Scan(&activityName, &activityAlias, &activityDesc)
		if err != nil {
			return []core.Activity{}, err
		}

		activities = append(activities, core.Activity{Id: activityId, Name: activityName, Alias: activityAlias, Description: activityDesc})
	}

	err = rows.Err()

	if err != nil {
		return []core.Activity{}, err
	}

	return activities, nil
}

// Find returns an activity
func (repo *SqliteRepository) Find(activityNameOrAlias string) (*core.Activity, error) {
	query := fmt.Sprintf("SELECT id, name, alias, description FROM activities WHERE name = '%s' OR alias = '%s'", activityNameOrAlias, activityNameOrAlias)
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var activity *core.Activity
	for rows.Next() {
		var activityId int
		var activityName string
		var activityAlias string
		var activityDesc string
		err = rows.Scan(&activityId, &activityName, &activityAlias, &activityDesc)
		if err != nil {
			return nil, err
		}

		activity = &core.Activity{Id: activityId, Name: activityName, Alias: activityAlias, Description: activityDesc}
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	if activity == nil {
		return nil, errors.New("activity not found")
	}

	return activity, nil
}

// Update updates an activity
func (repo *SqliteRepository) Update(activityNameOrAlias string, updateOp core.UpdateActivity) error {
	activity, err := repo.Find(activityNameOrAlias)
	if err != nil {
		return err
	}

	updateOp.Visit(activity)

	updateQuery := fmt.Sprintf("UPDATE activities SET name = '%s', alias = '%s', description = '%s' WHERE id = %v", activity.Name, activity.Alias, activity.Description, activity.Id)
	res, err := repo.db.Exec(updateQuery)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		panic("More than one activity updated! This shouldn't happen, please check your data")
	}

	return nil
}

// LogsForDay returns a list of activity logs for a given day
func (repo *SqliteRepository) LogsForDay(day time.Time) ([]core.ActivityLog, error) {
	d := utils.TimeToStandardFormat(day)

	query := fmt.Sprintf("SELECT id, day, started_at, stopped_at FROM activity_logs WHERE day = '%s'", d)
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var logs []core.ActivityLog

	for rows.Next() {
		var logId int
		var logDay string
		var startedAt time.Time
		var stoppedAt time.Time
		err = rows.Scan(&logId, &logDay, &startedAt, &stoppedAt)
		if err != nil {
			return nil, err
		}

		logs = append(logs, core.ActivityLog{Id: logId, Date: logDay, StartedAt: &startedAt, StoppedAt: &stoppedAt})
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return logs, nil
}

// Start starts tracking the time for an activity
func (repo *SqliteRepository) Start(activity core.Activity) error {
	activityLogs, err := repo.activityLogStartedAt(repo.Clock.Now())

	if err != nil {
		return err
	}

	var activityStartedAndNotStopped *core.Activity

	for i := range activityLogs {
		if activityLogs[i].StartedAt != nil && activityLogs[i].StoppedAt == nil {
			activityStartedAndNotStopped = &activityLogs[i].Activity
			break
		}
	}

	if activityStartedAndNotStopped != nil {
		return fmt.Errorf("you are already tracking the activity '%s', please stop that one before starting a new one", activityStartedAndNotStopped.Name)
	}

	sql := fmt.Sprintf("INSERT INTO activity_logs (day, started_at, activity_id) VALUES (DATE(), %v, %v)", repo.Clock.Now().Unix(), activity.Id)
	_, err = repo.db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

// Stop stops tracking the time for an activity
func (repo *SqliteRepository) Stop(activity core.Activity) error {
	activityLogs, err := repo.activityLogStartedAt(repo.Clock.Now())

	if err != nil {
		return err
	}

	if len(activityLogs) == 0 {
		return errors.New("you are not tracking any activity today, please start tracking one with the 'start' command")
	}

	var activityStartedAndNotStopped *core.Activity

	for i := range activityLogs {
		if activityLogs[i].StartedAt != nil && activityLogs[i].StoppedAt == nil {
			activityStartedAndNotStopped = &activityLogs[i].Activity
			break
		}
	}

	if activityStartedAndNotStopped != nil && activityStartedAndNotStopped.Id != activity.Id {
		return fmt.Errorf("you are already tracking the activity '%s', please stop that one before starting a new one", activityStartedAndNotStopped.Name)
	}

	var logForActivityToday *core.ActivityLog
	for i := range activityLogs {
		if activityLogs[i].Activity.Id == activity.Id {
			logForActivityToday = &activityLogs[i]
			break
		}
	}

	if logForActivityToday == nil {
		return fmt.Errorf("you are not tracking this activity. Please start tracking it with `tt start %s`", activity.Name)
	}

	if logForActivityToday.StoppedAt != nil {
		return errors.New("you have already stopped this activity")
	}

	updateQuery := fmt.Sprintf("UPDATE activity_logs SET stopped_at = %v WHERE id = %v", repo.Clock.Now().Unix(), logForActivityToday.Id)
	res, err := repo.db.Exec(updateQuery)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		panic("More than one activity log updated! This shouldn't happen, please check your data")
	}

	return nil
}

func (repo *SqliteRepository) activityLogStartedAt(instant time.Time) ([]core.ActivityLog, error) {
	date := instant.Format("2006-01-02")

	query := `
		SELECT activity_logs.id, 
			   activity_logs.day, 
			   activity_logs.started_at, 
			   activity_logs.stopped_at, 
			   activities.id,
			   activities.name,
			   activities.alias,
			   activities.description
		FROM activity_logs, activities
		WHERE day = DATE('%s') AND activities.id = activity_logs.activity_id
	`

	rows, err := repo.db.Query(fmt.Sprintf(query, date))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var activityLogs []core.ActivityLog
	for rows.Next() {
		var logId int
		var logDay string
		var logStartedAt *time.Time
		var logStoppedAt *time.Time
		var activityId int
		var activityName string
		var activityAlias string
		var activityDesc string
		err = rows.Scan(
			&logId,
			&logDay,
			&logStartedAt,
			&logStoppedAt,
			&activityId,
			&activityName,
			&activityAlias,
			&activityDesc,
		)
		if err != nil {
			return nil, err
		}

		activityLog := core.ActivityLog{
			Id:        logId,
			Date:      logDay,
			StartedAt: logStartedAt,
			StoppedAt: logStoppedAt,
			Activity: core.Activity{
				Id:          activityId,
				Name:        activityName,
				Alias:       activityAlias,
				Description: activityDesc,
			},
		}

		activityLogs = append(activityLogs, activityLog)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return activityLogs, nil
}
