package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type SensorData struct {
	Data []DataValue
}

type DataValue struct {
	Ai    []float32
	Aiabh []bool
	Aiabl []bool
	Aialh []bool
	Aiall []bool
	Di    []bool
	Diab  []bool
	Dial  []bool
	Time  string
}

func Call(client *http.Client, url string) SensorData {
	//client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err.Error())
	}

	var sensorDataObj SensorData
	json.Unmarshal(bodyBytes, &sensorDataObj)
	//log.Printf("API Response as struct %+v\n", sensorDataObj)
	return sensorDataObj
}
