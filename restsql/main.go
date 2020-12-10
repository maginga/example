package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

func main() {
	filename, _ := filepath.Abs("./config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	for {

		requestBody, err := json.Marshal(map[string]string{
			"query": config.Sql,
		})

		resp, err := http.Post(config.Url, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var result []map[string]interface{}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)
			println(str)
		}

		if err := json.Unmarshal(respBody, &result); err != nil {
			panic(err)
		}

		currentTime := time.Now()
		eventtime := fmt.Sprintf("%v", result[0]["__time"])
		oldTime, err := time.Parse(time.RFC3339, eventtime)
		if err != nil {
			fmt.Println(err)
		}
		diff := currentTime.Sub(oldTime)
		log.Println("elapsed time: ", diff.Seconds(), " (sec)")
		// 결과 출력
		// bytes, _ := ioutil.ReadAll(resp.Body)
		// str := string(bytes) //바이트를 문자열로
		// fmt.Println(str)

		time.Sleep(time.Second * 1)
	}
}
