package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func Authenticate(request Request) (authResponse AuthResponse) {

	authRequest := AuthRequest{request.Username, request.Password}

	client := getHTTPSClient()

	url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + SESSION
	reqBody, err := json.Marshal(authRequest)

	if err != nil {
		log.Println("Failed to parse the request payload.\n[ERROR] -", err)
		os.Exit(1)
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("HTTP request failed.\n[ERROR] -", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to parse HTTP response.\n[ERROR] -", err)
		os.Exit(1)
	}

	err = json.Unmarshal(body, &authResponse)

	if err != nil {
		log.Println("Failed to parse the response body.\n[ERROR] -", err)
		os.Exit(1)
	}

	return authResponse
}
