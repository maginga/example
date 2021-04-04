package rdb

import (
	"database/sql"
	"strings"
)

type DbType int8
type BindType int8

const (
	BtUnknown BindType = iota
	BtAt
	BtColon
	BtDollar
	BtQuestion

	DbUnknown DbType = iota
	// DbMySql
	// DbOracle
	// DbPostgres
	// DbSQLite
	DbSqlServer
)

// // ToTypeEnum get the data type that corresponds to the specified name
// func ToDbType(typeStr string) (DbType, error) {

// 	switch strings.ToLower(typeStr) {
// 	case "mysql":
// 		return DbMySql, nil
// 	case "oracle":
// 		return DbOracle, nil
// 	case "postgres":
// 		return DbPostgres, nil
// 	case "sqlite":
// 		return DbSQLite, nil
// 	case "sqlserver":
// 		return DbSqlServer, nil
// 	default:
// 		return DbUnknown, fmt.Errorf("unknown type: %s", typeStr)
// 	}
// }

type DbHelper interface {
	DbType() DbType
	BindType() BindType
	ToSQLStatementVal(val interface{}) string
	GetScanType(columnType *sql.ColumnType) interface{}
}

func GetDbHelper() (DbHelper, error) {

	// dbType, err := ToDbType(typeStr)
	// if err != nil {
	// 	return nil, err
	// }

	// switch dbType {
	// case DbMySql:
	// 	return &mySqlDBHelper{}, nil
	// case DbOracle:
	// 	return &oracleDBHelper{}, nil
	// case DbPostgres:
	// 	return &postgresDBHelper{}, nil
	// case DbSQLite:
	// 	return &sqliteDBHelper{}, nil
	// case DbSqlServer:
	// 	return &sqlserverDBHelper{}, nil
	// }

	return &sqlserverDBHelper{}, nil
	//return nil, fmt.Errorf("unsupported db: %s", typeStr)
}

// SQLServer

type sqlserverDBHelper struct {
}

func (*sqlserverDBHelper) DbType() DbType {
	return DbSqlServer
}

func (*sqlserverDBHelper) BindType() BindType {
	return BtAt
}

func (*sqlserverDBHelper) ToSQLStatementVal(val interface{}) string {

	switch t := val.(type) {
	case bool:
		if t {
			return "TRUE"
		} else {
			return "FALSE"
		}
	}

	return toSQLStatementVal(val)
}

func (*sqlserverDBHelper) GetScanType(columnType *sql.ColumnType) interface{} {

	if strings.HasPrefix(columnType.DatabaseTypeName(), "VARCHAR") {
		return new(string)
	}

	return new(interface{})
}

func toSQLStatementVal(val interface{}) string {
	switch t := val.(type) {
	case bool:
		if t {
			return "true"
		} else {
			return "false"
		}
	case int, int32, int64, float32, float64:
		s, _ := ToString(val)
		return s
	default:
		s, _ := ToString(val)
		return "'" + s + "'"
	}
}
