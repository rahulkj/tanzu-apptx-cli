package services

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type ServiceAccounts struct {
	url         string
	username    string
	password    string
	sa_username string
	sa_password string
	sa_alias    string
	operation   string
}

type serviceAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Alias    string `json:"alias"`
}

type serviceAccount struct {
	UUID     string `json:"uuid"`
	Alias    string `json:"alias"`
	Username string `json:"username"`
}

type response struct {
	Embedded struct {
		ServiceAccounts []struct {
			UUID     string `json:"uuid"`
			Alias    string `json:"alias"`
			Username string `json:"username"`
		} `json:"serviceAccounts"`
	} `json:"_embedded"`
}

func (serviceAccounts ServiceAccounts) Execute() {
	serviceAccounts = serviceAccounts.validate()

	request := Request{serviceAccounts.url, serviceAccounts.username, serviceAccounts.password}
	authResponse := Authenticate(request)

	switch serviceAccounts.operation {
	case "create":
		serviceAccounts.createServiceAccount(authResponse.Token)
	case "delete":
		serviceAccounts.deleteServiceAccount(authResponse.Token)
	default:
		fmt.Println("Operation not supported")
		os.Exit(1)
	}
}

func (serviceAccounts ServiceAccounts) validate() ServiceAccounts {
	saCmd := flag.NewFlagSet(SERVICE_ACCOUNT_CMD, flag.ExitOnError)
	operation := saCmd.String("operation", "", "create, delete")
	url := saCmd.String("url", "", "Iris URL, ex: appliance.example.com")
	username := saCmd.String("username", "", "Iris admin username")
	password := saCmd.String("password", "", "Iris admin password")
	sa_username := saCmd.String("service_username", "", "service account username")
	sa_password := saCmd.String("service_password", "", "service account password")
	sa_alias := saCmd.String("sa_alias", "", "service account alias")

	saCmd.Parse(os.Args[2:])

	if (*url == "" || *username == "" || *password == "") ||
		(*sa_username == "" || *sa_password == "" || *sa_alias == "") ||
		(strings.Contains(*url, "https://")) {
		fmt.Println("subcommand 'serviceAccount'")
		saCmd.PrintDefaults()
		os.Exit(1)
	}

	serviceAccounts = ServiceAccounts{*url, *username, *password, *sa_username, *sa_password, *sa_alias, *operation}
	return serviceAccounts
}

func (serviceAccounts ServiceAccounts) createServiceAccount(token string) {
	response := serviceAccounts.findServiceAccount(serviceAccounts.sa_alias, token)

	if len(response.Embedded.ServiceAccounts) > 0 {
		log.Println("Service Account already exists")
	} else {
		url := PROTOCOL + "://" + serviceAccounts.url + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "?action=register"
		request := serviceAccountRequest{serviceAccounts.username, serviceAccounts.password, serviceAccounts.sa_alias}
		body, _ := processRequest(token, url, "POST", request)

		serviceAccount := serviceAccount{}
		err := json.Unmarshal(body, &serviceAccount)
		if err != nil {
			log.Println("Failed to parse the response body.\n[ERROR] -", err)
			os.Exit(1)
		}

		if len(serviceAccount.UUID) > 0 {
			log.Println("Service Account created")
		}
	}
}

func (serviceAccounts ServiceAccounts) findServiceAccount(alias string, token string) (response response) {
	url := PROTOCOL + "://" + serviceAccounts.url + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "?page=0&size=10&alias=" + alias
	body, _ := processRequest(token, url, "GET", nil)

	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Failed to parse the response body.\n[ERROR] -", err)
		os.Exit(1)
	}
	return response
}

func (serviceAccounts ServiceAccounts) deleteServiceAccount(token string) {
	response := serviceAccounts.findServiceAccount(serviceAccounts.sa_alias, token)

	if len(response.Embedded.ServiceAccounts) > 0 {

		for _, serviceAccount := range response.Embedded.ServiceAccounts {
			if serviceAccount.Alias == serviceAccounts.sa_alias {
				url := PROTOCOL + "://" + serviceAccounts.url + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "/" + serviceAccount.UUID
				_, responseCode := processRequest(token, url, "DELETE", nil)

				if responseCode == 200 {
					log.Println("Deleted Service Account")
				} else {
					log.Println("Failed to delete Service Account")
				}
			} else {
				log.Println("Cannot delete Service Account as it does not exist")
			}
		}

	} else {
		log.Println("Cannot delete Service Account as it does not exist")
	}
}
