package services

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func getHTTPSClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	return client
}

func processRequest(token string, url string, method string, payload interface{}) (body []byte, responseCode int) {

	var req *http.Request

	if payload != nil {
		reqBody, err := json.Marshal(payload)

		if err != nil {
			log.Println("Failed to parse the request payload.\n[ERROR] -", err)
			os.Exit(1)
		}
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := getHTTPSClient().Do(req)
	if err != nil {
		log.Println("HTTP request failed.\n[ERROR] -", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to parse HTTP response.\n[ERROR] -", err)
		os.Exit(1)
	}

	return body, resp.StatusCode
}
