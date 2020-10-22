package sqlquery

import "github.com/project-flogo/core/data/coerce"

type Settings struct {
	DbType          string `md:"dbType,allowed(mysql,oracle,postgres,sqlite,sqlserver), required"`
	DriverName      string `md:"driverName,required"`
	DataSourceName  string `md:"dataSourceName,required"`
	Query           string `md:"query,required"`
	MaxOpenConns    int    `md:"maxOpenConnections"`
	MaxIdleConns    int    `md:"maxIdleConnections"`
	DisablePrepared bool   `md:"disablePrepared"`
	TimeZone        string `md:"timeZone"`
	StartOffset     string `md:"startOffset,required"`
	BatchSize       string `md:"batchSize,required"`
}

type Input struct {
	Params map[string]interface{} `md:"params"`
}

// FromMap converts the values from a map into the struct Input
func (i *Input) FromMap(values map[string]interface{}) error {
	params, err := coerce.ToObject(values["params"])
	if err != nil {
		return err
	}
	i.Params = params
	return nil
}

// ToMap converts the struct Input into a map
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"params": i.Params,
	}
}

type Output struct {
	// ColumnNames []interface{} `md:"columnNames"`
	Result interface{} `md:"result"`
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"result": o.Result,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Result, err = coerce.ToObject(values["result"])
	if err != nil {
		return err
	}

	return nil
}
