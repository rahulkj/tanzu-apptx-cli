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

	var url *string
	var username *string
	var password *string
	var sa_username *string
	var sa_password *string
	var sa_alias *string

	if operation == REGISTER {
		url = registerCmd.String("url", "", "Iris URL, ex: appliance.example.com")
		username = registerCmd.String("username", "", "Iris admin username")
		password = registerCmd.String("password", "", "Iris admin password")
		sa_username = registerCmd.String("service-username", "", "service account username")
		sa_password = registerCmd.String("service-password", "", "service account password")
		sa_alias = registerCmd.String("sa-alias", "", "service account alias")

		registerCmd.Parse(os.Args[3:])

		if (len(*url) == 0 || len(*username) == 0 || len(*password) == 0) ||
			(len(*sa_username) == 0 || len(*sa_password) == 0 || len(*sa_alias) == 0) ||
			(strings.Contains(*url, "https://")) {
			fmt.Println("Usage: 'iris-cli serviceAccount register [flags]' \n")
			fmt.Println("Flags:")
			registerCmd.PrintDefaults()
			os.Exit(1)
		}
	} else if operation == UNREGISTER {
		url = unregisterCmd.String("url", "", "Iris URL, ex: appliance.example.com")
		username = unregisterCmd.String("username", "", "Iris admin username")
		password = unregisterCmd.String("password", "", "Iris admin password")
		sa_alias = unregisterCmd.String("sa-alias", "", "service account alias")

		unregisterCmd.Parse(os.Args[3:])

		if (len(*url) == 0 || len(*username) == 0 || len(*password) == 0) ||
			(len(*sa_alias) == 0) ||
			(strings.Contains(*url, "https://")) {
			fmt.Println("Usage: 'iris-cli serviceAccount unregister [flags]' \n")
			fmt.Println("Flags:")
			unregisterCmd.PrintDefaults()
			os.Exit(1)
		}
	} else {
		serviceAccounts.printUsage()
	}

	serviceAccounts = ServiceAccounts{*url, *username, *password, *sa_username, *sa_password, *sa_alias, operation}
	return serviceAccounts
}

func (serviceAccounts ServiceAccounts) printUsage() {
	fmt.Println("Usage: 'iris-cli serviceAccount [command]' \n")
	fmt.Println("Available Commands:")
	fmt.Printf("  %s \t\t\t%s \n", REGISTER, "Register service account")
	fmt.Printf("  %s \t\t\t%s \n", UNREGISTER, "Unregister service account")
	os.Exit(1)
}

func (serviceAccounts ServiceAccounts) createServiceAccount(token string, request Request) {
	response := serviceAccounts.findServiceAccount(serviceAccounts.sa_alias, token, request)

	if len(response.Embedded.ServiceAccounts) > 0 {
		log.Println("Service Account already exists")
	} else {
		url := PROTOCOL + "://" + request.URL + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "?action=register"
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
	response := serviceAccounts.findServiceAccount(serviceAccounts.sa_alias, token, request)

	if len(response.Embedded.ServiceAccounts) > 0 {

		for _, serviceAccount := range response.Embedded.ServiceAccounts {
			if serviceAccount.Alias == serviceAccounts.sa_alias {
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
