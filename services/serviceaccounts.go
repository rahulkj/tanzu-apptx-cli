package services

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var url *string
var username *string
var password *string
var isDefault *bool
var sa_username *string
var sa_password *string
var sa_alias *string
var sa_type *string
var operation *string

type ServiceAccounts struct {
}

type ServiceAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Alias    string `json:"alias"`
}

type ServiceAccountResponse struct {
}

type ServiceAccount struct {
	UUID     string `json:"uuid"`
	Alias    string `json:"alias"`
	Username string `json:"username"`
}

type Response struct {
	Embedded struct {
		ServiceAccounts []struct {
			UUID     string `json:"uuid"`
			Alias    string `json:"alias"`
			Username string `json:"username"`
		} `json:"serviceAccounts"`
	} `json:"_embedded"`
}

func (serviceAccounts ServiceAccounts) Execute() {
	validate()

	request := Request{*url, *username, *password}
	authResponse := Authenticate(request)

	switch *operation {
	case "create":
		createServiceAccount(authResponse.Token)
	case "delete":
		deleteServiceAccount(authResponse.Token)
	default:
		fmt.Println("Operation not supported")
		os.Exit(1)
	}
}

func validate() {
	inputCmd := flag.NewFlagSet("serviceAccount", flag.ExitOnError)
	operation = inputCmd.String("operation", "", "create, delete")
	url = inputCmd.String("url", "", "Iris URL, ex: appliance.example.com")
	username = inputCmd.String("username", "", "Iris admin username")
	password = inputCmd.String("password", "", "Iris admin password")
	sa_username = inputCmd.String("service_username", "", "service account username")
	sa_password = inputCmd.String("service_password", "", "service account password")
	sa_alias = inputCmd.String("sa_alias", "", "service account alias")
	isDefault = inputCmd.Bool("default", false, "service account default (false)")
	sa_type = inputCmd.String("service_account_type", "", "service account type, ex: VCs, VRNIs, LINUX_VMs")

	inputCmd.Parse(os.Args[2:])

	if (*url == "" || *username == "" || *password == "") ||
		(*sa_username == "" || *sa_password == "" || *sa_alias == "") ||
		(*isDefault && *sa_type == "") ||
		(strings.Contains(*url, "https://")) {
		fmt.Println("subcommand 'serviceAccount'")
		inputCmd.PrintDefaults()
		os.Exit(1)
	}
}

func createServiceAccount(token string) {
	response := findServiceAccount(*sa_alias, token)

	if len(response.Embedded.ServiceAccounts) > 0 {
		log.Println("Service Account already exists")
	} else {
		url := PROTOCOL + "://" + *url + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "?action=register"
		request := ServiceAccountRequest{*username, *password, *sa_alias}
		body, _ := processRequest(token, url, "POST", request)

		serviceAccount := ServiceAccount{}
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

func findServiceAccount(alias string, token string) (response Response) {
	url := PROTOCOL + "://" + *url + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "?page=0&size=10&alias=" + alias
	body, _ := processRequest(token, url, "GET", nil)

	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Failed to parse the response body.\n[ERROR] -", err)
		os.Exit(1)
	}

	log.Println(response)
	return response
}

func deleteServiceAccount(token string) {
	response := findServiceAccount(*sa_alias, token)

	if len(response.Embedded.ServiceAccounts) > 0 {

		for _, serviceAccount := range response.Embedded.ServiceAccounts {
			if serviceAccount.Alias == *sa_alias {
				url := PROTOCOL + "://" + *url + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "/" + serviceAccount.UUID
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
