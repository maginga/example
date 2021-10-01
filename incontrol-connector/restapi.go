package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetBody(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	if resp.StatusCode == 200 {
		bodyBytes, err = ioutil.ReadAll(resp.Body)
	} else if err != nil {
		return nil, err
	} else {
		return nil, fmt.Errorf("The remote end did not return a HTTP 200 (OK) response.")
	} 

	return bodyBytes, nil
}
