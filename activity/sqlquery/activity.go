package sqlquery

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/project-flogo/contrib/activity/sqlquery/util"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

const (
	ovResults = "results"
)

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{MaxIdleConns: 2}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	dbHelper, err := util.GetDbHelper(s.DbType)
	if err != nil {
		return nil, err
	}

	ctx.Logger().Debugf("DB: '%s'", s.DbType)

	// todo move this to a shared connection object
	db, err := getConnection(s)
	if err != nil {
		return nil, err
	}

	sqlStatement, err := util.NewSQLStatement(dbHelper, s.Query)
	if err != nil {
		return nil, err
	}

	if sqlStatement.Type() != util.StSelect {
		return nil, fmt.Errorf("only select statement is supported")
	}

	f, e := time.Parse(time.RFC3339, s.StartOffset) // "2012-11-01T22:08:41+00:00"
	if e != nil {
		ctx.Logger().Debug("time parsing error.")
		return nil, e
	}
	min, _ := strconv.Atoi(s.BatchSize)
	t := f.Add(time.Minute * time.Duration(min))

	act := &Activity{db: db, dbHelper: dbHelper, sqlStatement: sqlStatement, settings: s, fromdate: &f, todate: &t}

	if !s.DisablePrepared {
		ctx.Logger().Debugf("Using PreparedStatement: %s", sqlStatement.PreparedStatementSQL())
		act.stmt, err = db.Prepare(sqlStatement.PreparedStatementSQL())
		if err != nil {
			return nil, err
		}
	}

	return act, nil
}

// Activity is a Counter Activity implementation
type Activity struct {
	dbHelper     util.DbHelper
	db           *sql.DB
	sqlStatement *util.SQLStatement
	stmt         *sql.Stmt
	settings     *Settings
	fromdate     *time.Time
	todate       *time.Time
}

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Cleanup() error {
	if a.stmt != nil {
		err := a.stmt.Close()
		log.RootLogger().Warnf("error cleaning up SQL Query activity: %v", err)
	}

	log.RootLogger().Tracef("cleaning up SQL Query activity")

	return a.db.Close()
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	in := &Input{}
	err = ctx.GetInputObject(in)
	if err != nil {
		return false, err
	}

	in.Params["fromdate"] = a.fromdate.Format(time.RFC3339)
	in.Params["todate"] = a.todate.Format(time.RFC3339)

	results, err := a.doSelect(in.Params)
	if err != nil {
		return false, err
	}

	output := &Output{Results: results}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}
	ctx.Logger().Debugf("result: %v", len(results))

	if len(results) > 0 {
		lastRow := results[len(results)-1]
		last := lastRow["event_time"]
		f, _ := time.Parse(time.RFC3339, last.(string)) // "2012-11-01T22:08:41+00:00"
		min, _ := strconv.Atoi(a.settings.BatchSize)
		t := a.fromdate.Add(time.Minute * time.Duration(min))

		a.fromdate = &f
		a.todate = &t
	}

	return true, nil
}

func (a *Activity) doSelect(params map[string]interface{}) ([]map[string]interface{}, error) {
	var err error
	var rows *sql.Rows

	if a.stmt != nil {
		args := a.sqlStatement.GetPreparedStatementArgs(params)
		rows, err = a.stmt.Query(args...)
	} else {
		rows, err = a.db.Query(a.sqlStatement.ToStatementSQL(params))
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results, err := getLabeledResults(a.dbHelper, rows)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func getLabeledResults(dbHelper util.DbHelper, rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columnTypes))
		for i := range values {
			values[i] = dbHelper.GetScanType(columnTypes[i])
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		resMap := make(map[string]interface{}, len(columns))
		for i, column := range columns {
			resMap[column] = *(values[i].(*interface{}))
		}

		//todo do we need to do column mapping
		results = append(results, resMap)
	}

	return results, rows.Err()
}

//todo move to shared connection
func getConnection(s *Settings) (*sql.DB, error) {

	db, err := sql.Open(s.DriverName, s.DataSourceName)
	if err != nil {
		return nil, err
	}

	if s.MaxOpenConns > 0 {
		db.SetMaxOpenConns(s.MaxOpenConns)
	}

	if s.MaxIdleConns != 2 {
		db.SetMaxIdleConns(s.MaxIdleConns)
	}

	return db, err
}
