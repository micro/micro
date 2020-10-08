// +build integration kind

package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func statusRunning(service, branch string, statusOutput []byte) bool {
	reg, _ := regexp.Compile(service + "\\s+" + branch + "\\s+\\S+\\s+running")
	return reg.Match(statusOutput)
}

func curl(serv Server, namespace, path string) (string, map[string]interface{}, error) {
	client := &http.Client{}
	url := fmt.Sprintf("http://127.0.0.1:%v/%v", serv.APIPort(), path)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Micro-Namespace", namespace)
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	m := map[string]interface{}{}
	return string(body), m, json.Unmarshal(body, &m)
}
