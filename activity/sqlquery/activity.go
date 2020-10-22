package sqlquery

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/project-flogo/contrib/activity/sqlquery/util"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
)

var layout string = "2006-01-02 15:04:05.000"

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

	ctx.Logger().Infof("DB: '%s'", s.DbType)

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

	utc, _ := time.Parse(layout, s.StartOffset)
	loc, _ := time.LoadLocation(s.TimeZone) //"Asia/Seoul"
	f := utc.In(loc)

	min, _ := strconv.Atoi(s.BatchSize)
	t := f.Add(time.Minute * time.Duration(min))

	ctx.Logger().Infof("start: %v", f)
	ctx.Logger().Infof("end: %v", t)

	act := &Activity{db: db, dbHelper: dbHelper, sqlStatement: sqlStatement, settings: s, fromdate: &f, todate: &t}

	if !s.DisablePrepared {
		ctx.Logger().Infof("Using PreparedStatement: %s", sqlStatement.PreparedStatementSQL())
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

	in.Params["fromdate"] = a.fromdate.Format(layout)
	in.Params["todate"] = a.todate.Format(layout)

	ctx.Logger().Infof("query from: %s, to: %s", in.Params["fromdate"], in.Params["todate"])

	r, err := a.doSelect(in.Params)
	if err != nil {
		return false, err
	}

	output := &Output{Result: r}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}

	if len(r) > 0 {
		lastRow := r[len(r)-1].(map[string]interface{})
		utc := lastRow["event_time"].(time.Time)
		loc, _ := time.LoadLocation(a.settings.TimeZone) //"Asia/Seoul"
		f := utc.In(loc)

		min, _ := strconv.Atoi(a.settings.BatchSize)
		t := a.fromdate.Add(time.Minute * time.Duration(min))

		a.fromdate = &f
		a.todate = &t
	}

	ctx.Logger().Infof("Result: %v rows.", len(r))
	return true, nil
}

func (a *Activity) doSelect(params map[string]interface{}) ([]interface{}, error) {
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

func getLabeledResults(dbHelper util.DbHelper, rows *sql.Rows) ([]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	log.RootLogger().Infof("columns: %v", columns)

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var results []interface{}

	for rows.Next() {
		values := make([]interface{}, len(columnTypes))
		for i := range values {
			values[i] = dbHelper.GetScanType(columnTypes[i])
			log.RootLogger().Infof("col type: %v", *columnTypes[i])
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		resMap := make(map[string]interface{}, len(columns))
		for i, column := range columns {
			resMap[column] = *(values[i].(*interface{}))
			//resMap[column] = values[i]
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
