package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/luispcosta/go-tt/config"
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

const DatabaseName = "gott.db"

// NewSqliteRepository creates a new SQLite repository struct
func NewSqliteRepository() (*SqliteRepository, error) {

	return &SqliteRepository{
		Clock: utils.NewLiveClock(),
	}, nil
}

// Initialize initializes the connection to the database
func (repo *SqliteRepository) Initialize(config config.Config) error {
	dbFilePath := fmt.Sprintf("%s%s", config.UserDataLocation, DatabaseName)
	db, err := sql.Open("sqlite3", dbFilePath)

	if err != nil {
		return nil
	}

	repo.db = db
	repo.dbFile = dbFilePath
	return nil
}

// Shutdown shutsdown the database
func (repo *SqliteRepository) Shutdown() error {
	return repo.db.Close()
}

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

// LogsForPeriod returns a list of activity logs for a given period
func (repo *SqliteRepository) LogsForPeriod(period core.Period) (map[string][]core.ActivityDurationDayAggregation, error) {
	stmt := `
		SELECT activities.id,
			   activities.name,
			   activities.alias,
			   activities.description,
			   agg.day,
			   agg.duration_in_seconds
		FROM (
			SELECT
				activity_id,
				day,
				SUM(CAST((JulianDay(stopped_at) - JulianDay(started_at)) * 24 * 60 * 60 AS integer)) AS duration_in_seconds
			FROM
				activity_logs
			WHERE
				day BETWEEN '%s' AND '%s'
			GROUP BY activity_id, day
		) AS agg, activities
		WHERE activities.id = agg.activity_id;
	`

	query := fmt.Sprintf(stmt, period.StartDateDay(), period.EndDateDay())
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make(map[string][]core.ActivityDurationDayAggregation)

	for rows.Next() {
		var activityId int
		var activityName string
		var activityAlias string
		var activityDescription string
		var day time.Time
		var durationInSeconds int

		err = rows.Scan(&activityId, &activityName, &activityAlias, &activityDescription, &day, &durationInSeconds)
		if err != nil {
			return nil, err
		}

		date := day.Format(utils.DateFormat)

		if err != nil {
			return nil, err
		}

		activity := core.Activity{Id: activityId, Name: activityName, Alias: activityAlias, Description: activityDescription}
		result[date] = append(result[date], core.ActivityDurationDayAggregation{Activity: activity, Date: date, Duration: durationInSeconds})
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Start starts tracking the time for an activity
func (repo *SqliteRepository) Start(activity core.Activity) error {
	activityStartedAndNotStopped, err := repo.CurrentlyTrackedActivity()

	if err != nil {
		return err
	}

	if activityStartedAndNotStopped != nil {
		return fmt.Errorf("you are already tracking the activity '%s', please stop that one before starting a new one", activityStartedAndNotStopped.Name)
	}

	startTime := utils.TimeToStandardDateTimeFormat(repo.Clock.Now())
	sql := fmt.Sprintf("INSERT INTO activity_logs (day, started_at, activity_id) VALUES (DATE(), '%v', '%v')", startTime, activity.Id)
	_, err = repo.db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

// WipeLogsPeriodAndActivity deletes logs for a given activity and for a given period
func (repo *SqliteRepository) WipeLogsPeriodAndActivity(period core.Period, activity *core.Activity) error {
	sql := `
		DELETE
		FROM activity_logs
		WHERE activity_id = %v AND day BETWEEN '%s' AND '%s';
	`

	query := fmt.Sprintf(sql, activity.Id, period.StartDateDay(), period.EndDateDay())
	res, err := repo.db.Exec(query)

	if err != nil {
		return err
	}

	_, errWipe := res.RowsAffected()

	if errWipe != nil {
		return errWipe
	}

	return nil
}

// WipeLogsPeriodAndActivity deletes logs for a given period
func (repo *SqliteRepository) WipeLogsPeriod(period core.Period) error {
	sql := `
		DELETE
		FROM activity_logs
		WHERE day BETWEEN '%s' AND '%s';
	`

	query := fmt.Sprintf(sql, period.StartDateDay(), period.EndDateDay())
	res, err := repo.db.Exec(query)

	if err != nil {
		return err
	}

	_, errWipe := res.RowsAffected()

	if errWipe != nil {
		return errWipe
	}

	return nil
}

// CurrentlyTrackedActivity returns the activity beeing currently tracked, if any
func (repo *SqliteRepository) CurrentlyTrackedActivity() (*core.Activity, error) {
	activityLogs, err := repo.activityLogStartedAt(repo.Clock.Now())

	if err != nil {
		return nil, err
	}

	var activityStartedAndNotStopped *core.Activity

	for i := range activityLogs {
		if activityLogs[i].StartedAt != nil && activityLogs[i].StoppedAt == nil {
			activityStartedAndNotStopped = &activityLogs[i].Activity
			break
		}
	}

	if activityStartedAndNotStopped != nil {
		return activityStartedAndNotStopped, nil
	}

	return nil, nil
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

	var logIncompleteToday *core.ActivityLog
	for i := range activityLogs {
		if activityLogs[i].Activity.Id == activity.Id && (activityLogs[i].StartedAt != nil && activityLogs[i].StoppedAt == nil) {
			logIncompleteToday = &activityLogs[i]
			break
		}
	}

	if logIncompleteToday == nil {
		return fmt.Errorf("you are not tracking this activity. Please start tracking it with `tt start %s`", activity.Name)
	}

	stopTime := utils.TimeToStandardDateTimeFormat(repo.Clock.Now())
	updateQuery := fmt.Sprintf("UPDATE activity_logs SET stopped_at = '%v' WHERE id = %v", stopTime, logIncompleteToday.Id)
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
	date := instant.Format(utils.DateFormat)

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
