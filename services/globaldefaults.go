package services

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type GlobalDefaults struct {
	url       string
	username  string
	password  string
	isDefault bool
	sa_alias  string
	sa_type   string
	operation string
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

func (globalDefaults GlobalDefaults) Execute() {
	globalDefaults = globalDefaults.validate()

	request := Request{globalDefaults.url, globalDefaults.username, globalDefaults.password}
	authResponse := Authenticate(request)

	switch globalDefaults.operation {
	case "create":
		globalDefaults.createServiceAccount(authResponse.Token)
	case "delete":
		globalDefaults.deleteServiceAccount(authResponse.Token)
	default:
		fmt.Println("Operation not supported")
		os.Exit(1)
	}
}

func (globalDefaults GlobalDefaults) validate() GlobalDefaults {
	inputCmd := flag.NewFlagSet(GLOBAL_DEFAULT_CMD, flag.ExitOnError)
	url := inputCmd.String("url", "", "Iris URL, ex: appliance.example.com")
	username := inputCmd.String("username", "", "Iris admin username")
	password := inputCmd.String("password", "", "Iris admin password")
	sa_alias := inputCmd.String("sa_alias", "", "service account alias")
	isDefault := inputCmd.Bool("default", false, "service account default (false)")
	sa_type := inputCmd.String("service_account_type", "", "service account type, ex: VCs, VRNIs, LINUX_VMs")
	operation := inputCmd.String("operation", "", "create, delete")

	inputCmd.Parse(os.Args[2:])

	if (*url == "" || *username == "" || *password == "") ||
		(*operation == "" || *sa_alias == "") ||
		(*isDefault && *sa_type == "") ||
		(strings.Contains(*url, "https://")) {
		fmt.Println("subcommand 'globalDefault'")
		inputCmd.PrintDefaults()
		os.Exit(1)
	}

	globalDefaults = GlobalDefaults{*url, *username, *password, *isDefault, *sa_alias, *sa_type, *operation}
	return globalDefaults
}

func (globalDefaults GlobalDefaults) createServiceAccount(token string) {
	response := globalDefaults.findServiceAccount(globalDefaults.sa_alias, token)

	if len(response.Embedded.ServiceAccounts) > 0 {
		log.Println("Service Account already exists")
	} else {
		url := PROTOCOL + "://" + globalDefaults.url + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "?action=register"
		request := ServiceAccountRequest{globalDefaults.username, globalDefaults.password, globalDefaults.sa_alias}
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

func (globalDefaults GlobalDefaults) findServiceAccount(alias string, token string) (response Response) {
	url := PROTOCOL + "://" + globalDefaults.url + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "?page=0&size=10&alias=" + alias
	body, _ := processRequest(token, url, "GET", nil)

	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Failed to parse the response body.\n[ERROR] -", err)
		os.Exit(1)
	}

	log.Println(response)
	return response
}

func (globalDefaults GlobalDefaults) deleteServiceAccount(token string) {
	response := globalDefaults.findServiceAccount(globalDefaults.sa_alias, token)

	if len(response.Embedded.ServiceAccounts) > 0 {

		for _, serviceAccount := range response.Embedded.ServiceAccounts {
			if serviceAccount.Alias == globalDefaults.sa_alias {
				url := PROTOCOL + "://" + globalDefaults.url + "/" + PREFIX + "/" + SERVICE_ACCOUNTS + "/" + serviceAccount.UUID
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
