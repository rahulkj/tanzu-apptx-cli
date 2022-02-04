package services

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
			fmt.Println("Failed to parse the request payload.\n[ERROR] -", err)
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
		fmt.Println("HTTP request failed.\n[ERROR] -", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to parse HTTP response.\n[ERROR] -", err)
		os.Exit(1)
	}

	return body, resp.StatusCode
}

func getCertificateThumbprint(endpoint string, port int, checksum string) (thumprint string) {

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", endpoint, port), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		panic("failed to connect: " + err.Error())
	}

	var fingerprint string

	cert := conn.ConnectionState().PeerCertificates[0]

	if checksum == "md5" {
		fingerprint = insertNth(strings.ToUpper(fmt.Sprintf("%x", md5.Sum(cert.Raw))), 2)
	} else if checksum == "sha1" {
		fingerprint = insertNth(strings.ToUpper(fmt.Sprintf("%x", sha1.Sum(cert.Raw))), 2)
	} else if checksum == "sha256" {
		fingerprint = insertNth(strings.ToUpper(fmt.Sprintf("%x", sha256.Sum256(cert.Raw))), 2)
	} else if checksum == "sha512" {
		fingerprint = insertNth(strings.ToUpper(fmt.Sprintf("%x", sha512.Sum512(cert.Raw))), 2)
	}

	conn.Close()

	return fingerprint
}

func insertNth(s string, n int) string {
	var buffer bytes.Buffer
	var n1 = n - 1
	var l1 = len(s) - 1
	for i, runei := range s {
		buffer.WriteRune(runei)
		if i%n == n1 && i != l1 {
			buffer.WriteRune(':')
		}
	}
	return buffer.String()
}
