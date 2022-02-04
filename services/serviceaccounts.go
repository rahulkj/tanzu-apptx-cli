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
	url        string
	username   string
	password   string
	saUsername string
	saPassword string
	saAlias    string
	operation  string
}

func (serviceAccounts ServiceAccounts) Execute() {
	serviceAccounts = serviceAccounts.validate()

	request := Request{serviceAccounts.url, serviceAccounts.username, serviceAccounts.password}
	authResponse := Authenticate(request)

	switch serviceAccounts.operation {
	case REGISTER:
		serviceAccounts.createServiceAccount(authResponse.Token, request)
	case UNREGISTER:
		serviceAccounts.deleteServiceAccount(authResponse.Token, request)
	default:
		fmt.Println("Operation not supported")
		serviceAccounts.printUsage()
		os.Exit(1)
	}
}

func (serviceAccounts ServiceAccounts) validate() ServiceAccounts {
	registerCmd := flag.NewFlagSet(REGISTER, flag.ExitOnError)
	unregisterCmd := flag.NewFlagSet(UNREGISTER, flag.ExitOnError)

	if len(os.Args) < 3 {
		serviceAccounts.printUsage()
	}

	operation := os.Args[2]

	var url string
	var username string
	var password string
	var saUsername string
	var saPassword string
	var saAlias string

	if operation == REGISTER {
		registerCmd.StringVar(&url, "fqdn", "", "Application Transformer FQDN / IP, ex: appliance.example.com")
		registerCmd.StringVar(&username, "username", "", "Application Transformer admin username")
		registerCmd.StringVar(&password, "password", "", "Application Transformer admin password")
		registerCmd.StringVar(&saUsername, "service-username", "", "service account username")
		registerCmd.StringVar(&saPassword, "service-password", "", "service account password")
		registerCmd.StringVar(&saAlias, "sa-alias", "", "service account alias")

		registerCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(len(saUsername) == 0 || len(saPassword) == 0 || len(saAlias) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Printf("Usage: '%s %s %s [flags]' \n", CLI_NAME, SERVICE_ACCOUNT_CMD, REGISTER)
			fmt.Println("Available Flags:")
			registerCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == UNREGISTER {
		unregisterCmd.StringVar(&url, "fqdn", "", "Application Transformer FQDN / IP, ex: appliance.example.com")
		unregisterCmd.StringVar(&username, "username", "", "Application Transformer admin username")
		unregisterCmd.StringVar(&password, "password", "", "Application Transformer admin password")
		unregisterCmd.StringVar(&saAlias, "sa-alias", "", "service account alias")

		unregisterCmd.Parse(os.Args[3:])

		if (len(url) == 0 || len(username) == 0 || len(password) == 0) ||
			(len(saAlias) == 0) ||
			(strings.Contains(url, "https://")) {
			fmt.Printf("Usage: '%s %s %s [flags]' \n", CLI_NAME, SERVICE_ACCOUNT_CMD, UNREGISTER)
			fmt.Println("Available Flags:")
			unregisterCmd.PrintDefaults()
			os.Exit(1)
		}
	} else {
		serviceAccounts.printUsage()
	}

	serviceAccounts = ServiceAccounts{url, username, password, saUsername, saPassword, saAlias, operation}
	return serviceAccounts
}

func (serviceAccounts ServiceAccounts) printUsage() {
	fmt.Printf("Usage: '%s %s [command]' \n", CLI_NAME, SERVICE_ACCOUNT_CMD)
	fmt.Println("Available Commands:")
	fmt.Printf("  %s \t\t\t%s \n", REGISTER, "Register service account")
	fmt.Printf("  %s \t\t\t%s \n", UNREGISTER, "Unregister service account")
	os.Exit(1)
}

func (serviceAccounts ServiceAccounts) createServiceAccount(token string, request Request) {
	response := serviceAccounts.findServiceAccount(serviceAccounts.saAlias, token, request)

	if len(response.Embedded.ServiceAccounts) > 0 {
		log.Println("Service Account already exists")
	} else {
		url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "?action=register"
		request := serviceAccountRequest{serviceAccounts.saUsername, serviceAccounts.saPassword, serviceAccounts.saAlias}
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

func (serviceAccounts ServiceAccounts) findServiceAccount(alias string, token string, request Request) (response response) {
	url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "?page=0&size=10&alias=" + alias
	body, _ := processRequest(token, url, "GET", nil)

	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Failed to parse the response body.\n[ERROR] -", err)
		os.Exit(1)
	}
	return response
}

func (serviceAccounts ServiceAccounts) deleteServiceAccount(token string, request Request) {
	response := serviceAccounts.findServiceAccount(serviceAccounts.saAlias, token, request)

	if len(response.Embedded.ServiceAccounts) > 0 {

		for _, serviceAccount := range response.Embedded.ServiceAccounts {
			if serviceAccount.Alias == serviceAccounts.saAlias {
				url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "/" + serviceAccount.UUID
				_, responseCode := processRequest(token, url, "DELETE", nil)

				if responseCode == 200 {
					log.Println("Deleted Service Account")
				} else {
					log.Println("Failed to delete Service Account. Response Code:", responseCode)
				}
			} else {
				log.Println("Cannot delete Service Account as it does not exist")
			}
		}

	} else {
		log.Println("Cannot delete Service Account as it does not exist")
	}
}
