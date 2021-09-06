package client

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

// DruidClient connect to druid
type DruidClient struct {
	Addr   string
	client *httpClient
}

type Message struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	Jti          string `json:"jti"`
}

// New returns a client
func NewDruidClient(addr string) (*DruidClient, error) {
	return &DruidClient{
		Addr:   addr,
		client: newHttpClient(),
	}, nil
}

func (d *DruidClient) url(path string) string {
	if strings.HasPrefix(d.Addr, "http") {
		return fmt.Sprintf("%s%s", d.Addr, path)
	}
	return fmt.Sprintf("http://%s%s", d.Addr, path)
}

func (c *DruidClient) getAuth() string {
	resp, err := http.PostForm(c.url("/oauth/token"),
		url.Values{"grant_type": {"password"},
			"client_id":     {"polaris_trusted"},
			"client_secret": {"secret"},
			"scope":         {"read write"},
			"username":      {"admin"},
			"password":      {"admin"},
		})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	var msg Message
	err = json.Unmarshal(respBody, &msg)
	if err != nil {
		panic(err)
	}
	return msg.AccessToken
}

func (c *DruidClient) InitOrganization() (string, error) {
	auth := "Bearer " + c.getAuth()
	println(auth)

	buff := bytes.NewBufferString("")
	req, err := http.NewRequest("POST", c.url("/api/organizations/init"), buff)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", auth)

	_, err = c.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	str := string(b)
	return str, err
}

func (c *DruidClient) ConfigureUser() (string, error) {
	gvURL := fmt.Sprintf("%v", viper.Get("grandview.url"))

	auth := "Bearer " + c.getAuth()
	println(auth)

	json := `{
		"clientName": "metatron Grandview",
		"logoFilePath": "/static/grandview.png",
		"redirectUri": "http://` + gvURL + `/api/oauth/login",
		"autoApprove": "true"
		}`

	buff := bytes.NewBufferString(json)
	req, err := http.NewRequest("POST", c.url("/api/oauth/polaris_trusted"), buff)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", auth)

	_, err = c.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	str := string(b)
	return str, err
}

func (c *DruidClient) Create(filePath string) (string, error) {
	auth := "Bearer " + c.getAuth()
	println(auth)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	buff := bytes.NewBufferString(string(content))
	req, err := http.NewRequest("POST", c.url("/api/datasources"), buff)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", auth)

	_, err = c.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	str := string(b)
	return str, err
}

func (c *DruidClient) CreateAlarm() (string, error) {
	bootstrapServer := fmt.Sprintf("%v", viper.Get("stream.bootstrapServers"))

	auth := "Bearer " + c.getAuth()
	println(auth)

	json := `{
		"name": "apm_alarm",
		"dsType": "MASTER",
		"connType": "ENGINE",
		"srcType": "REALTIME",
		"granularity": "SECOND",
		"segGranularity": "DAY",
		"ingestion": {
		  "type": "realtime",
		  "topic": "apm-alarm",
		  "consumerType": "KAFKA",
		  "consumerProperties": {
			"bootstrap.servers": "` + bootstrapServer + `"
		  },
		  "taskOptions": {
			"useEarliestOffset": true
		  },
		  "format": {
			"type": "json",
			"flattenRules": [
			  {
				"name": "event_time",
				"expr": "$.timestamp"
			  },
			  {
				"name": "alarmName",
				"expr": "$.name"
			  },
			  {
				"name": "rule",
				"expr": "$.alarmRule"
			  },
			  {
				"name": "paramId",
				"expr": "$.parameterNames"
			  }
			]
		  },
		  "rollup": false
		},
		"fields": [
		  {
			"name": "event_time",
			"type": "TIMESTAMP",
			"role": "TIMESTAMP",
			"seq": 0,
			"format": {
			  "type": "time_unix",
			  "unit": "millisecond"
			}
		  },
		  {
			"name": "uuid",
			"type": "STRING",
			"role": "DIMENSION",
			"seq": 1
		  },
		  {
			"name": "type",
			"type": "STRING",
			"role": "DIMENSION",
			"seq": 2
		  },
		  {
			"name": "assetId",
			"type": "STRING",
			"role": "DIMENSION",
			"seq": 3
		  },
		  {
			"name": "alarmName",
			"type": "STRING",
			"role": "DIMENSION",
			"seq": 4
		  },
		  {
			"name": "rule",
			"type": "STRING",
			"role": "DIMENSION",
			"seq": 5
		  },
		  {
			"name": "scores",
			"type": "ARRAY",
			"logicalType": "ARRAY",
			"role": "MEASURE",
			"seq": 6
		  },
		  {
			"name": "timestamps",
			"type": "ARRAY",
			"logicalType": "ARRAY",
			"role": "DIMENSION",
			"seq": 7
		  },
		  {
			"name": "assetStates",
			"type": "ARRAY",
			"logicalType": "ARRAY",
			"role": "DIMENSION",
			"seq": 8
		  },
		  {
			"name": "paramId",
			"type": "ARRAY",
			"logicalType": "ARRAY",
			"role": "DIMENSION",
			"seq": 9
		  },
		  {
			"name": "modelType",
			"type": "STRING",
			"role": "DIMENSION",
			"seq": 10
		  },
		  {
			"name": "modelName",
			"type": "STRING",
			"role": "DIMENSION",
			"seq": 11
		  }
		]
	  }`

	log.Printf("alarm: %s", json)
	buff := bytes.NewBufferString(json)
	req, err := http.NewRequest("POST", c.url("/api/datasources"), buff)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", auth)

	_, err = c.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	str := string(b)
	return str, err
}

func (c *DruidClient) CreateScore(nestID string) (string, error) {
	bootstrapServer := fmt.Sprintf("%v", viper.Get("stream.bootstrapServers"))

	auth := "Bearer " + c.getAuth()
	println(auth)

	json := `{
		"name": "apm_score_` + strings.ReplaceAll(nestID, "-", "_") + `",
		"dsType": "MASTER",
		"connType": "ENGINE",
		"srcType": "REALTIME",
		"granularity": "SECOND",
		"segGranularity": "DAY",
		"ingestion": {
		  "type": "realtime",
		  "topic": "apm-score-` + strings.ReplaceAll(nestID, "_", "-") + `",
		  "consumerType": "KAFKA",
		  "consumerProperties": {
			"bootstrap.servers": "` + bootstrapServer + `"
		  },
		  "taskOptions": {
			"useEarliestOffset": true
		  },
		  "format": {
			"type": "json",
			"flattenRules": [
			  {
				"name": "event_time",
				"expr": "$.timestamp"
			  },
			  {
				"name": "score",
				"expr": "$.assetScore"
			  },
			  {
				"name": "paramId",
				"expr": "$.parameters"
			  }
			]
		  },
		  "rollup": false
		},
		"fields": [
		  {
			"name": "event_time",
			"type": "TIMESTAMP",
			"role": "TIMESTAMP",
			"seq": 0,
			"format": {
			  "type": "time_unix",
			  "unit": "millisecond",
			  "timeZone": "Asia/Seoul"
			}
		  },
		  {
			"name": "assetId",
			"type": "STRING",
			"role": "DIMENSION",
			"seq": 1
		  },
		  {
			"name": "modelType",
			"type": "STRING",
			"role": "DIMENSION",
			"seq": 2
		  },
		  {
			"name": "paramId",
			"type": "ARRAY",
			"logicalType": "ARRAY",
			"role": "DIMENSION",
			"seq": 3
		  },
		  {
			"name": "healthIndexes",
			"type": "ARRAY",
			"logicalType": "ARRAY",
			"role": "MEASURE",
			"seq": 4
		  },
		  {
			"name": "score",
			"type": "DOUBLE",
			"logicalType": "DOUBLE",
			"role": "MEASURE",
			"seq": 5
		  },
		  {
			"name": "threshold",
			"type": "DOUBLE",
			"logicalType": "DOUBLE",
			"role": "MEASURE",
			"seq": 6
		  },
		  {
			"name": "modelName",
			"type": "STRING",
			"role": "DIMENSION",
			"seq": 7
		  },
		  {
			"name": "assetState",
			"type": "STRING",
			"role": "DIMENSION",
			"seq": 8
		  },
		  {
			"name": "context",
			"type": "STRING",
			"role": "DIMENSION",
			"seq": 9
		  }
		]
	  }`

	log.Printf("score: %s", json)
	buff := bytes.NewBufferString(json)
	req, err := http.NewRequest("POST", c.url("/api/datasources"), buff)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", auth)

	_, err = c.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	str := string(b)
	return str, err
}

func (c *DruidClient) CreateTrace(nestID string) (string, error) {
	bootstrapServer := fmt.Sprintf("%v", viper.Get("stream.bootstrapServers"))

	auth := "Bearer " + c.getAuth()
	println(auth)

	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var dataType string
	var physicalName string
	rows, err := db.Query(`SELECT p.data_type, p.physical_name
	FROM nest n, nest_egg g, sensor s, sensor_param_group_join j, parameter p
	WHERE n.id=g.nest_id
	and n.id=?
	and s.asset_id=g.asset_id
	and s.id = j.sensor_id
	and j.param_group_id = p.param_group_id
	GROUP BY p.data_type, p.physical_name`, nestID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	seq := 5
	var columns []string
	for rows.Next() {
		err := rows.Scan(&dataType, &physicalName)
		if err != nil {
			panic(err)
		}

		str := `
		{
			"name": "` + physicalName + `",
			"type": "` + dataType + `",
			"role": "MEASURE",
			"seq": ` + fmt.Sprintf("%v", seq) + `
		  }`

		columns = append(columns, str)
		seq++
	}

	var json string

	if len(columns) > 0 {
		columnJSON := strings.Join(columns, ",")
		json = `{
			"name": "apm_trace_` + strings.ReplaceAll(nestID, "-", "_") + `",
			"dsType": "MASTER",
			"connType": "ENGINE",
			"srcType": "REALTIME",
			"granularity": "SECOND",
			"segGranularity": "DAY",
			"ingestion": {
			  "type": "realtime",
			  "topic": "apm-trace-` + strings.ReplaceAll(nestID, "_", "-") + `",
			  "consumerType": "KAFKA",
			  "consumerProperties": {
				"bootstrap.servers": "` + bootstrapServer + `"
			  },
			  "taskOptions": {
				"useEarliestOffset": true
			  },
			  "format": {
				"type": "json"
			  },
			  "rollup": false
			},
			"fields": [
			  {
				"name": "event_time",
				"type": "TIMESTAMP",
				"role": "TIMESTAMP",
				"seq": 0,
				"format": {
				  "type": "time_format",
				  "format": "yyyy-MM-dd'T'HH:mm:ss.SSSZ",
				  "timeZone": "UTC",
				  "locale": "en"
				}
			  },
			  {
				"name": "sensorType",
				"type": "STRING",
				"role": "DIMENSION",
				"seq": 1
			  },
			  {
				"name": "sensorId",
				"type": "STRING",
				"role": "DIMENSION",
				"seq": 2
			  },
			  {
				"name": "assetId",
				"type": "STRING",
				"role": "DIMENSION",
				"seq": 3
			  },
			  {
				"name": "context",
				"type": "STRING",
				"role": "DIMENSION",
				"seq": 4
			  },` + columnJSON + `
			]
		  }`
	} else {
		json = `{
			"name": "apm_trace_` + strings.ReplaceAll(nestID, "-", "_") + `",
			"dsType": "MASTER",
			"connType": "ENGINE",
			"srcType": "REALTIME",
			"granularity": "SECOND",
			"segGranularity": "DAY",
			"ingestion": {
			  "type": "realtime",
			  "topic": "apm-trace-` + strings.ReplaceAll(nestID, "_", "-") + `",
			  "consumerType": "KAFKA",
			  "consumerProperties": {
				"bootstrap.servers": "` + bootstrapServer + `"
			  },
			  "taskOptions": {
				"useEarliestOffset": true
			  },
			  "format": {
				"type": "json"
			  },
			  "rollup": false
			},
			"fields": [
			  {
				"name": "event_time",
				"type": "TIMESTAMP",
				"role": "TIMESTAMP",
				"seq": 0,
				"format": {
				  "type": "time_format",
				  "format": "yyyy-MM-dd'T'HH:mm:ss.SSSZ",
				  "timeZone": "UTC",
				  "locale": "en"
				}
			  },
			  {
				"name": "sensorType",
				"type": "STRING",
				"role": "DIMENSION",
				"seq": 1
			  },
			  {
				"name": "sensorId",
				"type": "STRING",
				"role": "DIMENSION",
				"seq": 2
			  },
			  {
				"name": "assetId",
				"type": "STRING",
				"role": "DIMENSION",
				"seq": 3
			  },
			  {
				"name": "context",
				"type": "STRING",
				"role": "DIMENSION",
				"seq": 4
			  }
			]
		  }`
	}

	log.Printf("trace: %s", json)
	buff := bytes.NewBufferString(json)
	req, err := http.NewRequest("POST", c.url("/api/datasources"), buff)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", auth)

	_, err = c.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	str := string(b)
	return str, err
}
